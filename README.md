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

This installs a shell wrapper in `~/.zshrc` or `~/.bashrc` that allows `bmc profsel` to set `AWS_PROFILE` in your current shell.

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

**Fish shell (`~/.config/fish/config.fish`)**
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

Create `~/.config/bmc/config.toml`:

```toml
[mfa]
enabled = true
totp_script = "/usr/bin/rbw get my-aws-mfa-entry --field totp"
clipboard_command = "xclip -selection clipboard"

[ec2]
auto_start_stopped = "prompt"   # always | never | prompt
columns = ["InstanceId", "Name", "PrivateIP", "PublicIP", "State", "Hibernate", "Scheduler"]

[console]
firefox_containers = true   # open in Firefox container tab via Granted extension
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
bmc console -s ec2       # Open console at specific service
```

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

`~/.config/bmc/config.toml`:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `mfa.enabled` | bool | `false` | Enable MFA session management |
| `mfa.totp_script` | string | `""` | Command to generate TOTP codes |
| `mfa.clipboard_command` | string | `""` | Command to copy TOTP to clipboard |
| `ec2.auto_start_stopped` | string | `"prompt"` | `always` / `never` / `prompt` |
| `console.firefox_containers` | bool | `false` | Open console in Firefox container tab via [Granted](https://addons.mozilla.org/en-US/firefox/addon/granted/) extension |

## Migration from bash version

1. Install the new binary (same name: `bmc`)
2. Run `bmc install-shell-integration` (replaces `source bmc profsel` pattern)
3. Create `~/.config/bmc/config.toml` (replaces `~/.config/bmc/config.env`)
4. Run `bmc doctor` to verify setup

**Breaking changes:**
- `source bmc profsel` → `eval "$(bmc profsel)"` (handled automatically by shell wrapper)
- `~/.config/bmc/config.env` → `~/.config/bmc/config.toml`
- Shell completions now via cobra: `bmc completion bash/zsh` (replaces `bmc gencompletions`)

## Optional dependencies

Some commands require additional tools (checked lazily):

| Command | Requires |
|---------|----------|
| `ec2connect` SSH | `ssh` |
| `ec2connect` SSM | `aws` CLI v2 + `session-manager-plugin` |
| `ecsconnect` | `aws` CLI v2 + `session-manager-plugin` |

Run `bmc doctor` to check all dependencies with install instructions.
