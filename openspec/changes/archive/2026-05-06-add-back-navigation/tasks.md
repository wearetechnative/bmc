## 1. ui.Choose: distinguish ESC from Ctrl+C

- [x] 1.1 In `internal/ui/list.go`: add `wentBack bool` field to `listModel`
- [x] 1.2 In `listModel.Update`: separate ESC (`wentBack=true`, quit) from Ctrl+C (quit without setting `wentBack`)
- [x] 1.3 Export `var ErrBack = errors.New("ui: user navigated back")` in `internal/ui/list.go`
- [x] 1.4 In `Choose()`: after `p.Run()`, check `result.(listModel).wentBack` and return `("", ErrBack)` if true

## 2. profsel back navigation

- [x] 2.1 In `cmd/profsel.go` `selectProfileInteractive`: wrap the group→profile flow in a `for` loop
- [x] 2.2 After inner `ui.Choose` for profile: if `errors.Is(err, ui.ErrBack)` then `continue`, if `selectedName == ""` then return empty (Ctrl+C)
- [x] 2.3 Verify: run `./bmc profsel`, select a group, press ESC → returns to group list; press Ctrl+C → exits

## 3. ecsconnect back navigation

- [x] 3.1 In `cmd/ecsconnect.go`: restructure cluster→service→task→container selection into nested loops with `ErrBack` handling at each level
- [x] 3.2 Verify: run `bmc ecsconnect`, navigate to service level, press ESC → returns to cluster list

## 4. Update CHANGELOG

- [x] 4.1 Add entry under `## NEXT VERSION` documenting ESC-to-go-back in profsel and ecsconnect
