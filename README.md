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

### AWS Profile Selection
- `bmc profsel` - Select and set AWS profile environment variables

The `profsel` command helps you interactively select AWS profiles from your `~/.aws/config` file and set the `AWS_PROFILE` environment variable.

**Basic usage:**
```bash
# Interactive profile selection
. bmc profsel

# Pre-select a specific profile
. bmc profsel -p my-profile

# List available profiles
bmc profsel -l
```

**JSON output for scripting:**

The `--json` flag enables machine-readable output for integration with scripts and automation tools:

```bash
# Interactive selection with JSON output (progress visible)
PROFILE_DATA=$(bmc profsel --json 3>&1 >/dev/null)

# Non-interactive with specific profile
PROFILE_DATA=$(bmc profsel -p my-profile --json 3>&1 >/dev/null)

# Extract fields from JSON
PROFILE_NAME=$(echo "$PROFILE_DATA" | jq -r '.profile_name')
SOURCE_PROFILE=$(echo "$PROFILE_DATA" | jq -r '.source_profile')
PROFILE_ARN=$(echo "$PROFILE_DATA" | jq -r '.profile_arn')
```

**JSON output format:**
```json
{
  "source_profile": "my-org",
  "profile_name": "my-dev-profile",
  "profile_arn": "arn:aws:iam::123456789012:role/DevRole"
}
```

**Error cases:**
```json
{"error": "profile not found"}
{"error": "no profile selected"}
```

**File descriptor 3 support:**

When using `--json`, output is directed to file descriptor 3 (if available), allowing progress messages to remain visible during interactive selection. This provides the best experience for both interactive use and scripting:

- JSON output → fd 3 (captured in variable)
- Progress messages → stdout (visible during execution)
- Backward compatible: falls back to stdout if fd 3 is not redirected

### Other Commands
- `bmc console` - Open AWS Console in browser with profile selection
- `bmc ecsconnect` - Connect to ECS container
- `bmc gencompletions` - Generate shell completion scripts for bash or zsh

## Shell Completion

BMC supports tab-completion for bash and zsh shells. This enables auto-completion of commands and improves usability.

### Bash Completion

Add one of the following to your `~/.bashrc`:

**Option 1: Direct sourcing (recommended)**
```bash
source <(bmc gencompletions bash)
```

**Option 2: Save to file**
```bash
bmc gencompletions bash > ~/.bmc-completion.bash
echo 'source ~/.bmc-completion.bash' >> ~/.bashrc
```

**Option 3: System-wide installation (requires root)**
```bash
bmc gencompletions bash | sudo tee /etc/bash_completion.d/bmc
```

Then restart your shell or run: `source ~/.bashrc`

### Zsh Completion

Add one of the following to your `~/.zshrc`:

**Option 1: Direct sourcing (recommended)**
```bash
source <(bmc gencompletions zsh)
```

**Option 2: Save to completion directory**
```bash
mkdir -p ~/.zsh/completions
bmc gencompletions zsh > ~/.zsh/completions/_bmc

# Add to ~/.zshrc if not already present:
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit
compinit
```

Then restart your shell or run: `source ~/.zshrc`

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

