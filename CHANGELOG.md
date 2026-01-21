# BMC Changelog

## NEXT VERSION
- Feature: `bmc console` respects AWS_PROFILE environment variable when set
- Feature: `bmc console -p` (without value) forces profile selection even when AWS_PROFILE is set
- Enhancement: Reduced friction when AWS_PROFILE is already configured
- Fix: `bmc profsel` no longer exits the shell when profile selection is cancelled (Ctrl-C) or no profile is chosen
- Feature: `bmc ec2connect` automatically selects SSH connection when `-u` (username) or `-i` (identity file) flags are provided, eliminating unnecessary connection type prompt

## 2.7.0 - 18 sept 2025
- open profile selection when AWS_PROFILE is not set
- use filter in stead of table/choose
- cleanups

## 0.2.6.7
- cleanups

## 0.2.6.6
- Add -s flag to console option. user bmc console -s <service-name> to directly open the console with the prefered service.\

## 0.2.6.5
- Add -p flag to console option. user bmc console -p <profile-name> to directly open the console with the profile.\

## 0.2.6.4
- Set session duration to 3600s

## 0.2.6.3
- Fix e2stopstart function, rewriting function call

## 0.2.6.2
- Fix spinner ec2stopstart function

## 0.2.6.0
- Refactor ec2ls function, integrated in main library
- ec2find option. Search for keyword in selected aws-profile

## 0.2.5.3
- Fix: table height to fix items not being visible

## 0.2.5.2
- Fix: table height to fix items not being visible

## 0.2.5.1
- Feature: Search option for ec2ls. Now you can search through the output list for a string
- Fix bugs

## 0.2.4
- fix renaming function error ec2connect

## 0.2.3
- rename ec2ssh option to ec2connect
- add options ssm and ssh connection method connecting to ec2

## 0.2.2
- fix usage
- Feature add usage in help
- make VERSION unique for flake distribution

## 0.2.1
- another fix usage

## 0.2.0
- Breaking: renamed cli command
- Feature: add more commands to bmc
- Feature: more refactoring
- Fix: new sourcing fix

## 0.1.1
- Fix: sourcing, detect being sourced 

## 0.1.0

- new official project
