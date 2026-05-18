## 1. profilehelper.go — shared variable and resolution logic

- [x] 1.1 Add `var globalProfile string` package-level variable to `cmd/profilehelper.go`
- [x] 1.2 Update `ensureAWSProfile()` to check `globalProfile` first, before `AWS_PROFILE` env var, and run MFA check on the flag-provided profile

## 2. Register -p/--profile flag on each AWS operation command

- [x] 2.1 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ec2connect` `init()`
- [x] 2.2 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ec2ls` `init()`
- [x] 2.3 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ec2` `init()`
- [x] 2.4 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ec2scheduler` `init()`
- [x] 2.5 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ec2stopstart` `init()`
- [x] 2.6 Add `Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use")` in `ecsconnect` `init()`

## 3. Verify

- [x] 3.1 Build the binary and run `bmc ec2ls -p <profile>` — confirm it uses the correct profile without interactive picker
- [x] 3.2 Confirm `bmc ec2connect -p <profile>` skips the picker and runs MFA check
- [x] 3.3 Confirm `-p` flag takes priority over a set `AWS_PROFILE` env var
- [x] 3.4 Confirm `bmc console -p` and `bmc profsel -p` still behave as before
- [x] 3.5 Confirm unknown profile name produces a clear error message
