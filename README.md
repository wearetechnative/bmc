# BMC (Bill McCloud) Technative AWS/Terraform DevOps tools

These scripts simplify working with aws-cli and the AWS Console.

- [AWS Profile Switcher](./docs/aws-profile-select.md) - set Environment Vars to a profile from .aws/config
- [AWS aws_config2browserext](./docs/aws_config2browserext) - Convert config for AWS browser Externsion (https://addons.mozilla.org/en-US/firefox/addon/aws-extend-switch-roles3/)
- ...

## Commands

Run `bmc usage` to see all available commands. Key commands include:

### EC2 Management
- `bmc ec2ls` - List running EC2 instances
- `bmc ec2connect` - Connect to an EC2 instance via SSH or SSM
- `bmc ec2stopstart` - Stop or start an EC2 instance
- `bmc ec2find` - Find EC2 instances
- `bmc ec2scheduler` - Toggle InstanceScheduler tag to temporarily enable/disable instance scheduling

The `ec2scheduler` command helps manage instances that use the `InstanceScheduler` tag for automatic start/stop scheduling. Use this to:
- List all EC2 instances and their scheduler status (enabled/disabled/none)
- Toggle between `InstanceScheduler` and `InstanceScheduler_DISABLED` tags for instances with existing scheduler tags
- Get guidance and instructions to add scheduler tags to instances that don't have them
- Open AWS Console directly to the selected instance details page using `assumego` for immediate tag addition
- Temporarily disable scheduling for maintenance or long-running tasks
- Re-enable scheduling when ready

### Other Commands
- `bmc console` - Open AWS Console in browser with profile selection
- `bmc profsel` - Select and set AWS profile environment variables
- `bmc ecsconnect` - Connect to ECS container

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

