## ADDED Requirements

### Requirement: Optional positional argument for instance pre-selection
The `ec2connect` command SHALL accept an optional positional argument. When no positional argument and no `-i` flag is given, the command SHALL display the full interactive instance picker as before.

#### Scenario: No arguments shows full picker
- **WHEN** the user runs `bmc ec2connect` with no positional argument and no `-i` flag
- **THEN** the system SHALL show the full interactive instance picker with all instances
- **AND** behaviour SHALL be identical to pre-change functionality
