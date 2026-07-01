import re

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'r') as f:
    content = f.read()

content = re.sub(
    r'dcTheme := ""\n\s*if session\.PhishLure != nil \{\n\s*dcTheme = session\.PhishLure\.DeviceCodeTheme\n\s*\}',
    r'',
    content
)

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'w') as f:
    f.write(content)
print("Done fixing dcTheme!")
