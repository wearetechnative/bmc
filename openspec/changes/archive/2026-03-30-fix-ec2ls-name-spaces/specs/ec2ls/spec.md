# ec2ls Spec

## MODIFIED Requirements

### Requirement: List All EC2 Instances
The `bmc ec2ls` command SHALL list all EC2 instances in the current AWS profile with formatted instance details displayed in an interactive table, correctly parsing field values that contain spaces.

#### Scenario: Display instance information columns
- **WHEN** displaying the instance table
- **THEN** the command SHALL show at minimum: instance ID, private IP address, public IP address, state, hibernation status, instance name (from Name tag), and scheduler configuration status
- **AND** the command SHALL correctly parse and display instance information including Name tags that contain spaces
- **AND** the command SHALL use tab character as the field separator for awk processing
- **AND** the command SHALL preserve tabs during variable expansion by quoting variables
- **AND** the command SHALL correctly extract all fields regardless of spaces within field values
