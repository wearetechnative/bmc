# BMC (Bill McCloud) Technative AWS/Terraform DevOps tools

A single Go binary that simplifies working with AWS — profile selection, EC2/ECS operations, and console access.

## Installation

### Homebrew (macOS / Linux)
```bash
brew install wearetechnative/tap/bmc
```

### Nix — nix-env
```bash
nix-env -iA bmc -f https://github.com/wearetechnative/bmc/archive/main.tar.gz
```

### Nix — nix profile
```bash
nix profile add github:wearetechnative/bmc
```

### NixOS — configuration.nix
```nix
{
  inputs.bmc.url = "github:wearetechnative/bmc";
  # ...
  environment.systemPackages = [ inputs.bmc.packages.${system}.bmc ];
}
```

### Binary download (GitHub Releases)

Download from [GitHub Releases](https://github.com/wearetechnative/bmc/releases) for your platform:
- `linux/amd64`, `linux/arm64`
- `darwin/amd64` (Intel Mac), `darwin/arm64` (Apple Silicon)

## Setup

### 1. Shell integration (required for profsel)

```bash
bmc install-shell-integration
```

This installs a shell wrapper that allows `bmc profsel` to set `AWS_PROFILE` in your current shell.

- **zsh / bash**: appends to `~/.zshrc` or `~/.bashrc`
- **Fish**: writes `~/.config/fish/functions/bmc.fish` (auto-loaded, no restart needed)
- **NixOS + Fish**: prints a `programs.fish.functions` home-manager snippet instead of writing files

#### NixOS / home-manager

If `~/.zshrc` is managed by home-manager, the command will print manual snippets instead of writing to the file. Add the wrapper yourself using the method that matches your setup:

**home-manager (`home.nix`)**
```nix
programs.zsh.initContent = ''
  bmc() {
    if [[ "$1" == "profsel" ]]; then
      eval "$(command bmc profsel "$@")"
    else
      command bmc "$@"
    fi
  }
'';
```

**Fish shell (`~/.config/fish/functions/bmc.fish`)**
```fish
function bmc
  if test "$argv[1]" = "profsel"
    eval (command bmc profsel $argv)
  else
    command bmc $argv
  end
end
```

### 2. Configuration

Create `~/.config/bmc/config.json`:

```json
{
  "mfa": {
    "enabled": true,
    "totp_script": "/usr/bin/rbw get my-aws-mfa-entry --field totp",
    "copy_command": "wl-copy",
    "paste_command": "wl-paste | wtype -"
  },
  "ec2": {
    "auto_start_stopped": "prompt",
    "columns": ["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]
  },
  "console": {
    "firefox_containers": true
  }
}
```

The `columns` field controls which columns appear in EC2 instance tables (`ec2ls`, `ec2connect`, `ec2stopstart`, `ec2scheduler`, `ec2find`) and in what order. Available column names:

| Column | Description |
|--------|-------------|
| `InstanceId` | EC2 instance ID (e.g. `i-0abc123`) |
| `Name` | Value of the `Name` tag |
| `PrivateIP` | Private IPv4 address |
| `PublicIP` | Public IPv4 address (empty if none) |
| `State` | Instance state (`running`, `stopped`, etc.) |
| `Hibernate` | Whether hibernation is enabled (`yes`/`no`) |
| `Scheduler` | Whether InstanceScheduler tag is set (`yes`/`no`) |
| `Profile` | AWS profile name (always shown in `ec2find`) |

Unknown column names are silently rendered as `n/a`.


### 3. Check prerequisites

```bash
bmc doctor
```

## MFA

bmc handles MFA automatically — there is no separate `bmc mfa` command. When you run `bmc profsel` or `bmc console`, bmc checks whether your session credentials are still valid and refreshes them if needed.

**Requirements:**

1. Set `mfa.enabled = true` in `~/.config/bmc/config.json`
2. Add a `[profile-long-term]` section to `~/.aws/credentials` for the source profile that has MFA:

```ini
[technative-long-term]
aws_access_key_id     = AKIA...
aws_secret_access_key = ...
aws_mfa_device        = arn:aws:iam::123456789012:mfa/your-username
```

When the session expires, bmc prompts for a 6-digit TOTP code. If `totp_script` is configured, bmc runs that command via `sh -c` to fetch the code automatically and optionally copies it to the clipboard and pastes it into the focused window.

`totp_script` is executed via `sh -c`, so quoted arguments with spaces work correctly and interactive TUI selection tools (e.g. `gum filter`) render on the terminal and accept keyboard input.

### Clipboard integration

`copy_command` receives the TOTP code via **stdin** and copies it to the clipboard. `paste_command` runs 300ms later and simulates a paste keystroke in whatever window is currently focused — useful for automatically filling an MFA field in a browser.

**Wayland (wl-clipboard)**

```json
{
  "mfa": {
    "copy_command": "wl-copy",
    "paste_command": "wl-paste"
  }
}
```

- `wl-copy` reads stdin and writes it to the Wayland clipboard
- `wl-paste` outputs the clipboard content — combine with a type tool if you want to simulate keystrokes instead: `"paste_command": "wl-paste | wtype -"`

**X11 (xclip)**

```json
{
  "mfa": {
    "copy_command": "xclip -selection clipboard",
    "paste_command": "xdotool key ctrl+v"
  }
}
```

- `xclip -selection clipboard` reads stdin into the clipboard
- `xdotool key ctrl+v` simulates Ctrl+V in the focused window

Both fields are optional. If only `copy_command` is set, the code is copied but not auto-pasted. If `copy_command` is not set, no clipboard interaction occurs.

## Commands

### Profile selection
```bash
bmc profsel              # Interactive profile selection
bmc profsel -p myprofile # Pre-select a profile
bmc profsel -l           # List all profiles
bmc profsel --json       # JSON output for scripting
```

### AWS Console
```bash
bmc console              # Open console for selected/current profile
bmc console -p myprofile # Open console for specific profile
bmc console -p           # Force interactive profile selection (ignores AWS_PROFILE)
bmc console -s ec2       # Open console at specific service
```

Recently used profiles appear at the top of the interactive selector (last 10, labelled "recent"). History is stored in `~/.local/share/bmc/console-history.json`.

### EC2
```bash
bmc ec2ls                # List EC2 instances
bmc ec2connect           # Connect via SSH or SSM
bmc ec2connect -i i-xxx  # Connect to specific instance
bmc ec2connect -u ubuntu # SSH as specific user
bmc ec2stopstart         # Stop or start an instance
bmc ec2find <search>     # Find instances across profiles
bmc ec2scheduler         # Toggle InstanceScheduler tag
```

### ECS
```bash
bmc ecsconnect           # Interactive shell into ECS container
```

### Shell completion
```bash
# Bash
source <(bmc completion bash)

# Zsh
source <(bmc completion zsh)
```

### Other
```bash
bmc version              # Show version
bmc doctor               # System health check
bmc install-shell-integration  # Install profsel wrapper
```

## Configuration reference

`~/.config/bmc/config.json`:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `mfa.enabled` | bool | `false` | Enable MFA session management |
| `mfa.totp_script` | string | `""` | Command to generate TOTP codes |
| `mfa.copy_command` | string | `""` | Command to copy TOTP to clipboard (receives code via stdin) |
| `mfa.paste_command` | string | `""` | Command to simulate paste keystroke after copy (e.g. `xdotool key ctrl+v`); runs 300ms after copy |
| `ec2.auto_start_stopped` | string | `"prompt"` | `always` / `never` / `prompt` |
| `console.firefox_containers` | bool | `false` | Open console in Firefox container tab via [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) extension |
| `console.chrome_profiles` | bool | `false` | **Experimental.** Open console in a bmc-managed isolated Chrome profile per AWS profile |
| `console.chrome_binary` | string | `"google-chrome"` | Chrome binary to use when `chrome_profiles = true` (e.g. `"chromium"`, `"brave-browser"`) |

### Experimental: Chrome profile isolation

When `chrome_profiles = true`, `bmc console` opens the AWS console in a dedicated Chrome instance isolated per AWS profile. Profile data is stored at `~/.config/bmc/chrome/profiles/<profile-name>/`.

On first use for a profile, bmc copies your extensions and preferences from the default Chrome profile (without copying cookies or login data), so your usual extensions are available immediately.

```json
{
  "console": {
    "chrome_profiles": true,
    "chrome_binary": "google-chrome"
  }
}
```

For Brave or Chromium, set `chrome_binary` accordingly:

```json
{
  "console": {
    "chrome_profiles": true,
    "chrome_binary": "brave-browser"
  }
}
```

> **Note:** Profile directories under `~/.config/bmc/chrome/profiles/` can be deleted at any time to reset a profile. They are not managed by bmc after creation.

## Migration from bash version

1. Install the new binary (same name: `bmc`)
2. Run `bmc install-shell-integration` (replaces `source bmc profsel` pattern)
3. Create `~/.config/bmc/config.json` (replaces `~/.config/bmc/config.env`)
4. Run `bmc doctor` to verify setup

**Breaking changes:**
- `source bmc profsel` → `eval "$(bmc profsel)"` (handled automatically by shell wrapper)
- `~/.config/bmc/config.env` → `~/.config/bmc/config.json`
- Shell completions now via cobra: `bmc completion bash/zsh` (replaces `bmc gencompletions`)

## Optional dependencies

Some commands require additional tools (checked lazily):

| Command | Requires |
|---------|----------|
| `ec2connect` SSH | `ssh` |
| `ec2connect` SSM | `aws` CLI v2 + `session-manager-plugin` |
| `ecsconnect` | `aws` CLI v2 + `session-manager-plugin` |

Run `bmc doctor` to check all dependencies with install instructions.
