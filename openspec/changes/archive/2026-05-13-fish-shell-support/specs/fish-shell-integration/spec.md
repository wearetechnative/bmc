## ADDED Requirements

### Requirement: Install Fish wrapper to functions directory
When `$SHELL` ends with `fish` and the system is not NixOS, `bmc install-shell-integration` SHALL write a Fish wrapper function to `~/.config/fish/functions/bmc.fish`, creating the directory if it does not exist.

#### Scenario: Fish shell detected, clean install
- **WHEN** `$SHELL` ends with `fish`
- **AND** the system is not NixOS
- **AND** `~/.config/fish/functions/bmc.fish` does not exist
- **THEN** the system SHALL create `~/.config/fish/functions/bmc.fish` with the Fish wrapper function
- **AND** the system SHALL report success to stdout
- **AND** no shell restart or `source` is required (Fish auto-loads functions)

#### Scenario: Fish wrapper already installed
- **WHEN** `$SHELL` ends with `fish`
- **AND** the system is not NixOS
- **AND** `~/.config/fish/functions/bmc.fish` already exists
- **THEN** the system SHALL report "Shell integration already installed" and exit without modifying the file

#### Scenario: Fish functions directory does not exist
- **WHEN** `$SHELL` ends with `fish`
- **AND** `~/.config/fish/functions/` does not exist
- **THEN** the system SHALL create the directory before writing `bmc.fish`

### Requirement: NixOS Fish users receive manual instructions
When `$SHELL` ends with `fish` and NixOS is detected, `bmc install-shell-integration` SHALL skip auto-installation and print manual instructions including a `programs.fish.functions` home-manager snippet.

#### Scenario: NixOS detected with Fish shell
- **WHEN** `$SHELL` ends with `fish`
- **AND** `/etc/nixos/` directory exists
- **THEN** the system SHALL NOT write any files
- **AND** the system SHALL print an explanation that NixOS requires declarative configuration
- **AND** the system SHALL print a `programs.fish.functions` home-manager snippet for the bmc wrapper

### Requirement: NixOS detection via /etc/nixos/
The system SHALL detect NixOS by checking for the existence of the `/etc/nixos/` directory.

#### Scenario: /etc/nixos/ exists
- **WHEN** `/etc/nixos/` directory is present on the filesystem
- **THEN** the system SHALL treat the OS as NixOS for the purposes of shell integration

#### Scenario: /etc/nixos/ does not exist
- **WHEN** `/etc/nixos/` directory is absent
- **THEN** the system SHALL NOT treat the OS as NixOS and SHALL proceed with normal auto-install flow

### Requirement: Fish wrapper function syntax
The Fish wrapper written to `~/.config/fish/functions/bmc.fish` SHALL use valid Fish syntax and correctly wrap the `profsel` subcommand.

#### Scenario: profsel invoked via Fish wrapper
- **WHEN** the user runs `bmc profsel` in Fish shell with the wrapper installed
- **THEN** `eval (command bmc profsel $argv)` is executed
- **AND** `AWS_PROFILE` is set in the current shell session

#### Scenario: other subcommands pass through
- **WHEN** the user runs any `bmc` subcommand other than `profsel`
- **THEN** `command bmc $argv` is executed directly without eval wrapping
