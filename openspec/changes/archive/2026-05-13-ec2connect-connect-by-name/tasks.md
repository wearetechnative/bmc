## 1. Update ec2connect command

- [x] 1.1 Add optional positional argument support to `ec2connectCmd` in `cmd/ec2connect.go` (`Args: cobra.MaximumNArgs(1)`)
- [x] 1.2 In `runEC2Connect`, after loading instances, check if `args[0]` is provided and apply case-insensitive substring filter on `InstanceID + Name + PrivateIP + PublicIP`
- [x] 1.3 If `-i` flag is set and a positional arg is also provided, print a warning to stderr and proceed with `-i` (ignore the positional arg)
- [x] 1.4 If 0 instances match the filter, return an error with a clear message
- [x] 1.5 If exactly 1 instance matches, use its `InstanceID` directly (skip `selectInstanceID`)
- [x] 1.6 If 2+ instances match, pass the filtered slice to `selectInstanceID` (existing function, no changes needed)

## 2. Verify existing behaviour unchanged

- [x] 2.1 Confirm `bmc ec2connect` (no args) still shows the full instance picker
- [x] 2.2 Confirm `bmc ec2connect -i <id>` still connects directly without showing a picker
