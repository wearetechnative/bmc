# Change: Add back navigation to TUI lists

## Why

When navigating multi-level TUI menus (profsel: group→profile, ecsconnect: cluster→service→task→container), pressing ESC or Ctrl+C cancels the entire flow. There is no way to go back to the previous menu level without restarting the command.

## What Changes

- Add `var ErrBack = errors.New("ui: user navigated back")` sentinel to `internal/ui`
- In `listModel`: distinguish ESC (go back) from Ctrl+C (cancel all)
- `ui.Choose()` returns `("", ErrBack)` on ESC; `("", nil)` on Ctrl+C
- `cmd/profsel.go` `selectProfileInteractive`: loop on `ErrBack` at profile level, returning to group selection
- `cmd/ecsconnect.go`: loop on `ErrBack` at service, task, and container levels

## Capabilities

### New Capabilities

- `tui-back-navigation`: ESC-to-go-back behavior in multi-level TUI list navigation

### Modified Capabilities

- `tui-list-display`: key binding semantics changed — ESC now means "go back" instead of "cancel"
- `profile-selection`: back navigation from profile list to group list
- `shell-integration`: no requirement changes (ecsconnect behavior change is internal)

## Impact

- `internal/ui/list.go`: add `wentBack bool` to `listModel`, sentinel error export
- `cmd/profsel.go`: loop in `selectProfileInteractive`
- `cmd/ecsconnect.go`: loops at service, task, container levels
- No impact on `ec2connect`, `ec2find`, `ec2stopstart` (single-level or fixed flows)
- No breaking changes for callers that treat `ErrBack` as an error (it is non-nil, so existing `if err != nil` guards still work — callers must opt in to back navigation)
