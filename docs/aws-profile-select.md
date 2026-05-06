# AWS Profile Selection

## Overview

`bmc profsel` is an interactive TUI for selecting an AWS profile from your `~/.aws/config`. After selection it sets `AWS_PROFILE` in your current shell session.

The TUI renders with full colors and styling even when run inside `eval "$(...)"` subshells, by redirecting I/O through `/dev/tty`.

---

## Prerequisites

### AWS CLI configuration

Your `~/.aws/config` must contain named profiles. Profiles are grouped by a common prefix (e.g. `ACT-`, `AMS-`). Example:

```ini
[default]
region = eu-west-1

[profile ACT-mycompany]
region = eu-west-1
role_arn = arn:aws:iam::123456789012:role/DevOpsAdministrator
source_profile = mycompany-userauth

[profile ACT-mycompany-staging]
region = eu-west-1
role_arn = arn:aws:iam::210987654321:role/DevOpsAdministrator
source_profile = mycompany-userauth
```

For MFA-based profiles using a source profile and long-term credentials, see the [MFA documentation](../README.md).

---

## Installation

`bmc` is distributed as a single binary. See the [releases page](https://github.com/wearetechnative/bmc/releases) or install via Homebrew:

```bash
brew install wearetechnative/tap/bmc
```

### Shell integration

To make `bmc profsel` set `AWS_PROFILE` in your current shell, install the shell wrapper:

```bash
bmc install-shell-integration
```

This appends a `bmc()` function to your `~/.zshrc` or `~/.bashrc`:

```zsh
bmc() {
  if [[ "$1" == "profsel" ]]; then
    eval "$(command bmc profsel "$@")"
  else
    command bmc "$@"
  fi
}
```

If your rc file is managed by home-manager or is otherwise not writable, the command prints manual snippets for zsh, bash, and Fish.

---

## Usage

### Interactive profile selection

With the shell wrapper installed, simply run:

```bash
bmc profsel
```

Or use the suggested alias:

```bash
alias aws-switch='bmc profsel'
aws-switch
```

A two-step TUI appears:

1. Select an account **group**
2. Select a **profile** within that group

After confirming, `AWS_PROFILE` is set in your current shell.

### Non-interactive (pre-select a profile)

```bash
bmc profsel --profile ACT-mycompany
```

### List all profiles

```bash
bmc profsel --list
```

### JSON output

```bash
bmc profsel --json
```

Returns `{"source_profile": "...", "profile_name": "...", "profile_arn": "..."}`.

---

## Color rendering in subshell contexts

When `bmc profsel` is invoked via `eval "$(command bmc profsel)"`, stdout is a pipe rather than a terminal. `bmc` detects this and opens `/dev/tty` directly for TUI I/O, so colors, borders, and highlighting render correctly. If `/dev/tty` is unavailable (e.g. in a CI environment), the TUI falls back to a plain numbered list on stderr.

---

## Keyboard shortcuts

| Key       | Action                        |
|-----------|-------------------------------|
| `↑` / `↓` | Navigate items                |
| `/`       | Filter/search                 |
| `Enter`   | Confirm selection             |
| `Esc`     | Go back (group → profile nav) |
| `Ctrl+C`  | Cancel without selecting      |
