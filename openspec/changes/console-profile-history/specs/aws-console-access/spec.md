## MODIFIED Requirements

### Requirement: Interactive profile selection shows recent profiles
When `AWS_PROFILE` is not set and no `-p` flag is given, the `bmc console` command SHALL present the interactive profile selector with recently used profiles shown at the top of the list.

#### Scenario: AWS_PROFILE is not set, history exists
- **WHEN** user runs `bmc console` and `AWS_PROFILE` environment variable is not set
- **AND** the console history file contains one or more entries
- **THEN** the command SHALL prompt for profile selection using the interactive profile selector
- **AND** recent profiles SHALL appear at the top of the list with a "recent" label

#### Scenario: AWS_PROFILE is not set, no history
- **WHEN** user runs `bmc console` and `AWS_PROFILE` environment variable is not set
- **AND** the console history file is empty or does not exist
- **THEN** the command SHALL prompt for profile selection using the interactive profile selector without a recent section
