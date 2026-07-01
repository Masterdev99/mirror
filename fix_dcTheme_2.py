import re

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'r') as f:
    content = f.read()

content = re.sub(
    r'\s*dcTheme := l\.DeviceCodeTheme\n',
    r'\n',
    content
)

content = re.sub(
    r'\s*dcTheme := ""\n\s*if l \!= nil \{\n\s*dcTheme = l\.DeviceCodeTheme\n\s*\}\n',
    r'\n',
    content
)

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'w') as f:
    f.write(content)
print("Done fixing remaining dcThemes!")
