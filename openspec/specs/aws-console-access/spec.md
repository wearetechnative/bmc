# aws-console-access Specification

## Purpose
This specification defines the behavior of the `bmc console` command, which provides convenient access to the AWS Console in Firefox with automatic profile management. The command supports environment variable integration, flexible profile selection, service shortcuts, and profile listing functionality.
## Requirements
### Requirement: Respect AWS_PROFILE Environment Variable
The `bmc console` command SHALL check for the `AWS_PROFILE` environment variable and use it when set, avoiding redundant profile selection prompts.

#### Scenario: AWS_PROFILE is set
- **WHEN** user runs `bmc console` and `AWS_PROFILE` environment variable is set
- **THEN** the command SHALL use the profile from `AWS_PROFILE` without prompting for selection

#### Scenario: AWS_PROFILE is not set
- **WHEN** user runs `bmc console` and `AWS_PROFILE` environment variable is not set
- **THEN** the command SHALL prompt for profile selection using the interactive profile selector

### Requirement: Force Profile Selection with -p Flag
The `bmc console` command SHALL support a `-p` flag without arguments to force profile selection even when `AWS_PROFILE` is set.

#### Scenario: Force selection when AWS_PROFILE is set
- **WHEN** user runs `bmc console -p` and `AWS_PROFILE` environment variable is set
- **THEN** the command SHALL ignore `AWS_PROFILE` and prompt for profile selection

#### Scenario: Specify profile directly with -p argument
- **WHEN** user runs `bmc console -p <profile-name>`
- **THEN** the command SHALL use the specified profile name directly without prompting

### Requirement: Service Selection with -s Flag
The `bmc console` command SHALL support a `-s <service>` flag to open a specific AWS service in the console.

#### Scenario: Open specific service
- **WHEN** user runs `bmc console -s ec2`
- **THEN** the command SHALL open the AWS console directly to the EC2 service page

#### Scenario: Combine service selection with profile
- **WHEN** user runs `bmc console -s s3` with or without profile flags
- **THEN** the command SHALL open the AWS console to the S3 service page using the selected/specified profile

### Requirement: List Profiles with -l Flag
The `bmc console` command SHALL support a `-l` flag to list available AWS profiles without opening the console.

#### Scenario: List available profiles
- **WHEN** user runs `bmc console -l`
- **THEN** the command SHALL print the list of available AWS profiles and exit without opening the console

