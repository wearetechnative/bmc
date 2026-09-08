## MODIFIED Requirements

### Requirement: Force Profile Selection with -p Flag
The `bmc console` command SHALL use `-P`/`--pick` (a value-less boolean flag) to force interactive profile selection even when `AWS_PROFILE` is set. The `-p`/`--profile` flag SHALL always take a profile name as its value (both `-p NAME` and `-p=NAME` forms), and SHALL NOT be usable without a value. A bare `-p` with no value SHALL be an error (`flag needs an argument`), not a request for interactive selection.

#### Scenario: Force selection when AWS_PROFILE is set
- **WHEN** user runs `bmc console -P` (or `bmc console --pick`) and `AWS_PROFILE` environment variable is set
- **THEN** the command SHALL ignore `AWS_PROFILE` and prompt for profile selection

#### Scenario: Specify profile directly with -p argument
- **WHEN** user runs `bmc console -p <profile-name>` (space form) or `bmc console -p=<profile-name>` (equals form)
- **THEN** the command SHALL use the specified profile name directly without prompting

#### Scenario: Bare -p is no longer a force-selection trigger
- **WHEN** user runs `bmc console -p` with no value
- **THEN** the command SHALL report that the flag needs an argument
- **AND** SHALL NOT prompt for interactive selection
