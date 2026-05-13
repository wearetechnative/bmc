---
title: "Shell Integration"
weight: 10
---

Shell integration is required for `bmc profsel` to set `AWS_PROFILE` in your current shell session.

## Install

```bash
bmc install-shell-integration
```

This installs a shell wrapper function:

- **zsh / bash**: appends to `~/.zshrc` or `~/.bashrc`
- **Fish**: writes `~/.config/fish/functions/bmc.fish` (auto-loaded, no restart needed)
- **NixOS + Fish**: prints a `programs.fish.functions` snippet instead of writing files

After installation, restart your shell or source the config file.

## How it works

The wrapper intercepts `bmc profsel` and evaluates the output so `AWS_PROFILE` is set in the parent shell:

**zsh / bash**
```bash
bmc() {
  if [[ "$1" == "profsel" ]]; then
    eval "$(command bmc profsel "$@")"
  else
    command bmc "$@"
  fi
}
```

**Fish**
```fish
function bmc
  if test "$argv[1]" = "profsel"
    eval (command bmc profsel $argv)
  else
    command bmc $argv
  end
end
```

## NixOS / home-manager

If your shell config is managed by home-manager, `bmc install-shell-integration` prints manual snippets instead of writing files.

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

**Fish (`home.nix`)**
```nix
programs.fish.functions.bmc = ''
  if test "$argv[1]" = "profsel"
    eval (command bmc profsel $argv)
  else
    command bmc $argv
  end
'';
```
