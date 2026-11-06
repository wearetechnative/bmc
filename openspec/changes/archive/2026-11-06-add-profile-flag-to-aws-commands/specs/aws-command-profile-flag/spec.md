## ADDED Requirements

### Requirement: AWS operation commands accept a --profile flag
The commands `ec2connect`, `ec2ls`, `ec2`, `ec2scheduler`, `ec2stopstart`, and `ecsconnect` SHALL each accept a `-p`/`--profile <name>` flag that specifies the AWS profile to use, bypassing both the `AWS_PROFILE` environment variable check and the interactive profile picker.

#### Scenario: User specifies a valid profile via -p flag
- **WHEN** the user runs an AWS command with `-p <profile-name>`
- **THEN** the command SHALL use the specified profile without prompting for interactive selection
- **AND** SHALL still run the MFA validity check for that profile

#### Scenario: User specifies a valid profile via --profile flag
- **WHEN** the user runs an AWS command with `--profile <profile-name>`
- **THEN** the command SHALL behave identically to using `-p <profile-name>`

#### Scenario: User specifies an unknown profile name
- **WHEN** the user runs an AWS command with `-p <nonexistent-profile>`
- **THEN** the command SHALL return an error indicating the profile was not found

#### Scenario: -p flag takes priority over AWS_PROFILE env var
- **GIVEN** `AWS_PROFILE` is set in the environment
- **WHEN** the user also passes `-p <different-profile>`
- **THEN** the command SHALL use the profile from the `-p` flag, not the environment variable

#### Scenario: console and profsel -p flags are unaffected
- **WHEN** the user runs `bmc console -p` or `bmc profsel -p`
- **THEN** those commands SHALL behave exactly as before this change
