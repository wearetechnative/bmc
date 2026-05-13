---
title: "NixOS / home-manager"
weight: 20
---

## Installing BMC

### NixOS (configuration.nix)

```nix
{
  inputs.bmc.url = "github:wearetechnative/bmc";

  outputs = { nixpkgs, bmc, ... }: {
    nixosConfigurations.myhost = nixpkgs.lib.nixosSystem {
      modules = [
        {
          environment.systemPackages = [
            bmc.packages.${system}.bmc
          ];
        }
      ];
    };
  };
}
```

### home-manager

```nix
{
  inputs.bmc.url = "github:wearetechnative/bmc";

  home.packages = [ inputs.bmc.packages.${system}.bmc ];
}
```

## Shell integration with home-manager

`bmc install-shell-integration` prints snippets instead of writing files when it detects a managed shell config. Add them manually:

### zsh

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

### bash

```nix
programs.bash.initExtra = ''
  bmc() {
    if [[ "$1" == "profsel" ]]; then
      eval "$(command bmc profsel "$@")"
    else
      command bmc "$@"
    fi
  }
'';
```

### Fish

```nix
programs.fish.functions.bmc = ''
  if test "$argv[1]" = "profsel"
    eval (command bmc profsel $argv)
  else
    command bmc $argv
  end
'';
```

## Shell completions

```nix
programs.zsh.initContent = ''
  source <(bmc completion zsh)
'';

# or for bash:
programs.bash.initExtra = ''
  source <(bmc completion bash)
'';
```
