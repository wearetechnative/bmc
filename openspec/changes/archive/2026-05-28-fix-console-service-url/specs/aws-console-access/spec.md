## MODIFIED Requirements

### Requirement: Service Selection with -s Flag
The `bmc console` command SHALL support a `-s <path>` flag to open a specific AWS service or console sub-page. The value is treated as a console path and may include a `/` to target a sub-page (e.g., `systems-manager/parameters`).

#### Scenario: Open specific service by short name
- **WHEN** user runs `bmc console -s rds`
- **THEN** the command SHALL open the AWS console at `https://<region>.console.aws.amazon.com/rds/home` using the region from the selected profile

#### Scenario: Open console sub-page with path
- **WHEN** user runs `bmc console -s systems-manager/parameters`
- **THEN** the command SHALL open the AWS console at `https://<region>.console.aws.amazon.com/systems-manager/parameters` using the region from the selected profile

#### Scenario: Combine service selection with profile
- **WHEN** user runs `bmc console -s ec2` with or without profile flags
- **THEN** the command SHALL open the AWS console to the EC2 service page using the selected/specified profile and the region resolved from that profile

#### Scenario: No service flag
- **WHEN** user runs `bmc console` without `-s`
- **THEN** the command SHALL open the AWS console homepage at `https://<region>.console.aws.amazon.com/`
