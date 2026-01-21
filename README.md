# BMC (Bill McCloud) Technative AWS/Terraform DevOps tools

These scripts simplify working with aws-cli and the AWS Console.

- [AWS Profile Switcher](./docs/aws-profile-select.md) - set Environment Vars to a profile from .aws/config
- [AWS aws_config2browserext](./docs/aws_config2browserext) - Convert config for AWS browser Externsion (https://addons.mozilla.org/en-US/firefox/addon/aws-extend-switch-roles3/)
- ...

## Configuration

BMC can be configured via `~/.config/bmc/config.env`. Available options:

### EC2 Instance Auto-Start
- `BMC_AUTO_START_STOPPED_INSTANCES` - Controls behavior when selecting stopped instances in `bmc ec2connect`
  - `"prompt"` (default) - Ask user before starting stopped instances
  - `"always"` - Automatically start stopped instances without prompting
  - `"never"` - Never start stopped instances, show error and exit

Example:
```bash
BMC_AUTO_START_STOPPED_INSTANCES="always"
```

### MFA / TOTP Configuration
- `totpScript` - Array containing command and arguments to generate TOTP codes for MFA authentication
- `clipboardCopyCommand` - Array containing command and arguments to copy text to clipboard
- `clipboardPasteCommand` - Array containing command and arguments to paste text from clipboard

Examples:
```bash
# Using rbw-menu.sh for TOTP generation
totpScript=("/path/to/rbw-menu.sh" "-t" "code" "-q" "new")

# Using pass for TOTP generation
totpScript=("pass" "otp" "aws/mfa")

# Simple TOTP script without arguments
totpScript=("/usr/local/bin/get-totp.sh")

# Clipboard commands (Linux with xclip)
clipboardCopyCommand=("xclip" "-selection" "clipboard")
clipboardPasteCommand=("xclip" "-selection" "clipboard" "-o")

# Clipboard commands (macOS)
clipboardCopyCommand=("pbcopy")
clipboardPasteCommand=("pbpaste")

# Clipboard commands (custom wrapper)
clipboardCopyCommand=("/usr/local/bin/clipcopy")
clipboardPasteCommand=("/usr/local/bin/clippaste")
```

**Note**: All commands should be configured as bash arrays to properly handle arguments and paths with spaces.

## TODO

- [ ] ci testing
- [ ] central command?
- [ ] naming conventions
- [ ] documentation (github pages)
- [ ] share code?
- [ ] coding style
- [ ] 2 versions of aws_config2browserext(2)

