# Implementation Tasks

## 1. Update console function in bmc script
- [x] 1.1 Modify getopts pattern from 'lp:s:' to 'lp::s:' to make -p optional
- [x] 1.2 Add forceProfileSelection flag when -p is provided without argument
- [x] 1.3 Add conditional logic to check AWS_PROFILE before selectAWSProfile call
- [x] 1.4 When AWS_PROFILE is set and -p not used, assign it to selectedProfileName variable
- [x] 1.5 Ensure setMFA is called with correct sourceProfile when using AWS_PROFILE

## 2. Testing
- [x] 2.1 Test `bmc console` when AWS_PROFILE is set (should use existing profile)
- [x] 2.2 Test `bmc console` when AWS_PROFILE is not set (should prompt)
- [x] 2.3 Test `bmc console -p` when AWS_PROFILE is set (should prompt)
- [x] 2.4 Test `bmc console -p <profile>` (should use specified profile)
- [x] 2.5 Test `bmc console -s <service>` with AWS_PROFILE set
- [x] 2.6 Test `bmc console -l` (should list profiles)
- [x] 2.7 Test combinations: `bmc console -p -s ec2`
- [x] 2.8 Verify on both bash and zsh shells

## 3. Documentation
- [x] 3.1 Update CHANGELOG.md with new behavior
- [x] 3.2 Consider updating README or docs if console command is documented
