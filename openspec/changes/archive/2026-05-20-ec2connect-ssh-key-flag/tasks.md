## 1. Flag and Logic

- [x] 1.1 Add `ec2connectKey` variable and `-k`/`--key` flag in `cmd/ec2connect.go` `init()`
- [x] 1.2 Extend auto-select SSH condition to also trigger when `ec2connectKey != ""`
- [x] 1.3 Pass `ec2connectKey` to `connectSSH` and append `-i <path>` to the ssh args when set

## 2. Documentation

- [x] 2.1 Update `docs/content/commands/ec2connect.md` to document the `-k`/`--key` flag
