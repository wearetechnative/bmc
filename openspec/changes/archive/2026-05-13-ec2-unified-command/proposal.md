## Why

Users frequently want to act on an EC2 instance immediately after finding it, but currently must repeat the instance selection in a separate command (`ec2connect`, `ec2stopstart`, `ec2scheduler`). A unified `bmc ec2` command eliminates this friction by combining selection and action in one flow. `ec2ls` intentionally stays display-only (like Unix `ls`) and will not be modified.

Linked to bean: `.beans/bmc-n5k1--ec2ls-context-menu.md`

## What Changes

- New command `bmc ec2 [search]` — interactive instance picker followed by an action menu
- Accepts an optional positional search argument (case-insensitive substring match on InstanceID, Name, PrivateIP, PublicIP — identical behaviour to `ec2connect [search]`)
- Single match → skips picker, shows action menu directly
- No match → clear error message
- Action menu items: **Connect SSH**, **Connect SSM**, **Start/Stop**, **Toggle scheduler**
- All actions delegate to existing internal functions already used by `ec2connect`, `ec2stopstart`, and `ec2scheduler`
- Existing commands (`ec2ls`, `ec2connect`, `ec2stopstart`, `ec2scheduler`, `ec2find`) remain unchanged

## Capabilities

### New Capabilities

- `ec2-unified`: Single-entry-point EC2 command that combines instance selection with an action menu

### Modified Capabilities

_(none — no existing specs change their requirements)_

## Impact

- New file: `cmd/ec2.go` (new Cobra subcommand)
- `cmd/root.go`: register `ec2Cmd`
- Reuses internal helpers from `cmd/ec2connect.go`, `cmd/ec2stopstart.go`, `cmd/ec2scheduler.go` — these may need minor refactoring to export shared logic or the new command can call the same `awsops` / `ui` functions directly
