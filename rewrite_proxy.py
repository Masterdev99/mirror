with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'r') as f:
    content = f.read()

import re

# 1. Remove the 5 regex vars:
content = re.sub(r'\n\t// Document access themed device code pages \(5 themes\).*?\n\tdcSharePointRe\s*=\s*regexp\.MustCompile\(`\^/access/sharepoint/\(\[a-zA-Z0-9_-\]\+\)\$`\)', '', content, flags=re.DOTALL)

# 2. Remove the themedRoutes block from 1021 to 1089:
content = re.sub(r'\n\t\t\t// --- Begin themed document access device code endpoints ---.*?\t\t\t// --- End themed document access device code endpoints ---', '', content, flags=re.DOTALL)

with open('/Users/macpro14/Downloads/evil-token-main/core/http_proxy.go', 'w') as f:
    f.write(content)
print("Done proxy!")
