# Implementation Tasks

## 1. Core Implementation

- [x] 1.1 Create `ec2scheduler.sh` script with main function
- [x] 1.2 Implement tag listing functionality to query instances with `InstanceScheduler` or `InstanceScheduler_DISABLED` tags
- [x] 1.3 Implement interactive instance selection using `gum` table
- [x] 1.4 Implement tag toggle logic to rename between `InstanceScheduler` ↔ `InstanceScheduler_DISABLED` while preserving tag value
- [x] 1.5 Add clear user feedback messages for toggle operations
- [x] 1.6 Register `ec2scheduler` command in main `bmc` dispatcher script

## 2. Integration & Testing

- [x] 2.1 Test tag listing with instances having `InstanceScheduler` tag
- [x] 2.2 Test tag listing with instances having `InstanceScheduler_DISABLED` tag
- [x] 2.3 Test tag listing with instances having both or neither tag
- [x] 2.4 Test toggle from `InstanceScheduler` to `InstanceScheduler_DISABLED`
- [x] 2.5 Test toggle from `InstanceScheduler_DISABLED` to `InstanceScheduler`
- [x] 2.6 Verify tag value preservation during toggle
- [x] 2.7 Test error handling for instances without the expected tags
- [x] 2.8 Verify AWS profile selection integration
- [x] 2.9 Test command help output and documentation

## 3. Documentation

- [x] 3.1 Update README.md with `ec2scheduler` command documentation
- [x] 3.2 Add entry to CHANGELOG.md under "NEXT VERSION"
- [x] 3.3 Ensure help text in command registration is clear and descriptive
