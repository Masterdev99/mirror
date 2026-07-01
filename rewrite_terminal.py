import re

with open('/Users/macpro14/Downloads/evil-token-main/core/terminal.go', 'r') as f:
    content = f.read()

# 1. Replace the themes list in handleQuickstart (around line 380)
# and in handlePhishlets (around line 2100)
# We need to remove the 4 themes.
content = re.sub(
    r'\{\"onedrive\", \"OneDrive\"\},\n\s*\{\"authenticator\", \"Authenticator\"\},\n\s*\{\"adobe\", \"Adobe PDF\"\},\n\s*\{\"docusign\", \"DocuSign\"\},\n',
    '',
    content,
    flags=re.MULTILINE
)

# 2. Replace the theme labels map (around line 2285)
content = re.sub(
    r'\"onedrive\":\s*\"OneDrive\",\n\s*\"authenticator\":\s*\"Authenticator\",\n\s*\"adobe\":\s*\"Adobe PDF\",\n\s*\"docusign\":\s*\"DocuSign\",\n',
    '',
    content,
    flags=re.MULTILINE
)

# 3. Replace validThemes array (around line 2560)
content = re.sub(
    r'validThemes := \[\]string\{\"default\", \"onedrive\", \"authenticator\", \"adobe\", \"docusign\", \"sharepoint\"\}',
    r'validThemes := []string{"default", "sharepoint"}',
    content
)

# 4. Replace error message (around line 2568)
content = re.sub(
    r'\"edit: invalid theme \'%s\' \(valid: default, onedrive, authenticator, adobe, docusign, sharepoint\)\"',
    r'"edit: invalid theme \'%s\' (valid: default, sharepoint)"',
    content
)

# 5. Replace readline autocomplete valid themes (around line 3240)
content = re.sub(
    r'readline.PcItem\(\"onedrive\"\), readline.PcItem\(\"authenticator\"\), readline.PcItem\(\"adobe\"\), readline.PcItem\(\"docusign\"\), ',
    '',
    content
)

# 6. Replace lures help string (around line 3268)
content = re.sub(
    r'set themed landing page \(default, onedrive, authenticator, adobe, docusign, sharepoint\)',
    r'set themed landing page (default, sharepoint)',
    content
)

# 7. Replace switch block in lure display (around line 3760)
content = re.sub(
    r'case \"onedrive\":\n\s*dcThemeStr = hiblue.Sprint\(dcTheme\)\n\s*case \"authenticator\":\n\s*dcThemeStr = hcyan.Sprint\(dcTheme\)\n\s*case \"adobe\":\n\s*dcThemeStr = color.New\(color.FgHiRed\).Sprint\(dcTheme\)\n\s*case \"docusign\":\n\s*dcThemeStr = yellow.Sprint\(dcTheme\)\n\s*',
    '',
    content,
    flags=re.MULTILINE
)


with open('/Users/macpro14/Downloads/evil-token-main/core/terminal.go', 'w') as f:
    f.write(content)
print("Done terminal!")
