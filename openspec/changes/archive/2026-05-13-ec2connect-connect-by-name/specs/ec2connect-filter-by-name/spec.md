## ADDED Requirements

### Requirement: Filter instances by partial name via positional argument
The `ec2connect` command SHALL accept an optional positional argument as a search fragment. When provided, it SHALL filter the full list of instances to those whose `InstanceID`, `Name`, `PrivateIP`, or `PublicIP` contain the fragment (case-insensitive substring match).

#### Scenario: Single match proceeds without picker
- **WHEN** the user runs `bmc ec2connect <fragment>`
- **AND** exactly one instance matches the fragment
- **THEN** the system SHALL use that instance directly
- **AND** the system SHALL NOT show the interactive instance picker

#### Scenario: Multiple matches show filtered picker
- **WHEN** the user runs `bmc ec2connect <fragment>`
- **AND** two or more instances match the fragment
- **THEN** the system SHALL show the interactive instance picker
- **AND** the picker SHALL contain only the matching instances

#### Scenario: No matches returns an error
- **WHEN** the user runs `bmc ec2connect <fragment>`
- **AND** no instances match the fragment
- **THEN** the system SHALL return an error indicating no instances were found for that fragment
- **AND** the system SHALL NOT show the interactive instance picker

#### Scenario: Search is case-insensitive
- **WHEN** the user provides a fragment in any casing (e.g. `NIXHOST`, `nixhost`, `NixHost`)
- **THEN** the system SHALL match instances whose Name, InstanceID, PrivateIP, or PublicIP contain the fragment regardless of case

#### Scenario: Search matches InstanceID fragment
- **WHEN** the user provides a fragment that is a substring of an instance ID (e.g. `i-0abc`)
- **THEN** the system SHALL include that instance in the matches

#### Scenario: Search matches IP fragment
- **WHEN** the user provides a fragment that is a substring of a private or public IP (e.g. `10.0.1`)
- **THEN** the system SHALL include instances whose PrivateIP or PublicIP contains the fragment

### Requirement: Positional argument does not override -i flag
When both the `-i` instance ID flag and a positional search argument are provided to `ec2connect`, the `-i` flag SHALL take precedence and the positional argument SHALL be ignored.

#### Scenario: -i flag wins over positional argument
- **WHEN** the user runs `bmc ec2connect -i <instance-id> <fragment>`
- **THEN** the system SHALL connect to the instance specified by `-i`
- **AND** the system SHALL print a warning to stderr that the positional argument is being ignored

#### Scenario: No positional argument with -i flag behaves as before
- **WHEN** the user runs `bmc ec2connect -i <instance-id>` without a positional argument
- **THEN** the system SHALL behave exactly as before this change
