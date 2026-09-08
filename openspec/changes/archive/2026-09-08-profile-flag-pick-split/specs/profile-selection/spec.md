## ADDED Requirements

### Requirement: Profile Selection Command-Line Flags
Every command that selects an AWS profile SHALL expose two distinct, single-purpose flags:

- `-p` / `--profile <NAME>` — SHALL accept a profile name as its value. Both the space form (`-p NAME`) and the equals form (`-p=NAME`) SHALL be accepted. The flag SHALL NOT consume or interfere with positional arguments such as an instance search term.
- `-P` / `--pick` — a boolean flag that SHALL force the interactive profile picker, ignoring any `AWS_PROFILE` environment variable.

When neither flag is provided, the command SHALL use `AWS_PROFILE` if it is set, and otherwise fall back to the interactive picker.

Bare `-p` with no value SHALL NOT be a valid way to request interactive selection; forcing interactive selection is done exclusively via `-P`/`--pick`.

This contract SHALL apply uniformly to: `ec2`, `ec2connect`, `ec2ls`, `ec2scheduler`, `ec2stopstart`, `ecsconnect`, `console`, and `profsel`.

#### Scenario: Profile name and positional search argument together
- **WHEN** the user runs `bmc ec2connect -p TN-Production compute2`
- **THEN** the system SHALL resolve the profile to `TN-Production`
- **AND** SHALL treat `compute2` as the instance search term
- **AND** SHALL NOT report an "accepts at most 1 arg(s)" error

#### Scenario: Equals form for profile value
- **WHEN** the user runs `bmc ec2connect -p=TN-Production compute2`
- **THEN** the system SHALL resolve the profile to `TN-Production`
- **AND** SHALL treat `compute2` as the instance search term

#### Scenario: Force interactive picker with a positional argument
- **WHEN** the user runs `bmc ec2connect --pick compute2`
- **THEN** the system SHALL open the interactive profile picker, ignoring `AWS_PROFILE`
- **AND** after a profile is chosen SHALL treat `compute2` as the instance search term

#### Scenario: Force interactive picker via short flag
- **WHEN** the user runs a profile-taking command with `-P`
- **THEN** the system SHALL open the interactive profile picker, ignoring `AWS_PROFILE`

#### Scenario: No profile flag falls back to environment then picker
- **GIVEN** `AWS_PROFILE` is set to a valid profile
- **WHEN** the user runs a profile-taking command with neither `-p` nor `-P`
- **THEN** the system SHALL use the `AWS_PROFILE` value
- **AND WHEN** `AWS_PROFILE` is not set
- **THEN** the system SHALL open the interactive picker

#### Scenario: Bare -p no longer triggers interactive selection
- **WHEN** the user runs a profile-taking command with a bare `-p` and no value
- **THEN** the system SHALL report that the flag needs an argument
- **AND** SHALL NOT open the interactive picker
