## 1. New command file

- [x] 1.1 Create `cmd/ec2.go` with `ec2Cmd` Cobra command (`Use: "ec2 [search]"`, `Args: cobra.MaximumNArgs(1)`)
- [x] 1.2 Register `ec2Cmd` in `cmd/root.go` via `rootCmd.AddCommand(ec2Cmd)`

## 2. Instance selection with optional search

- [x] 2.1 In `runEC2` load config and list instances (same as `ec2connect`)
- [x] 2.2 Implement search filter: case-insensitive substring match on InstanceID+Name+PrivateIP+PublicIP
- [x] 2.3 Handle zero-match case: return error with search term
- [x] 2.4 Handle single-match case: skip picker, use that instance directly
- [x] 2.5 Handle multi-match / no-arg case: show interactive picker via `selectInstanceID`

## 3. Action menu

- [x] 3.1 After instance is selected, call `awsops.GetInstanceState` to determine current state
- [x] 3.2 Build action list: "Connect SSH", "Connect SSM", and dynamic "Start instance" / "Stop instance" based on state; "Toggle scheduler"
- [x] 3.3 Show action menu via `ui.Choose`; empty return (ESC) exits cleanly
- [x] 3.4 "Connect SSH" branch: call `connectSSH(instanceID)` (reuse from `ec2connect.go`)
- [x] 3.5 "Connect SSM" branch: call `connectSSM(instanceID)` (reuse from `ec2connect.go`)
- [x] 3.6 "Start instance" branch: call `startInstance(profile, instanceID)` (reuse from `ec2stopstart.go`)
- [x] 3.7 "Stop instance" branch: call `stopInstance(profile, instanceID, instances)` (reuse from `ec2stopstart.go`)
- [x] 3.8 "Toggle scheduler" branch: replicate scheduler toggle logic from `ec2scheduler.go` (confirm + `awsops.ToggleSchedulerTag`)
