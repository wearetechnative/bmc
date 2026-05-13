## Why

`ec2connect` requires either knowing an exact instance ID or scrolling through a full interactive list. Users working with named instances (e.g. `ec2-nixhost-prod-992382728492`) should be able to type a partial name fragment to jump directly to the right instance — removing friction for common scripted and interactive workflows.

## What Changes

- `ec2connect` accepts an optional positional argument (partial search string) to pre-filter instances before selection
- If the search matches exactly one instance, connection proceeds immediately without a picker
- If the search matches multiple instances, the existing interactive table picker is shown with the filtered results
- If the search matches zero instances, an error is returned
- If `-i` (instance ID flag) is provided alongside the positional argument, `-i` takes precedence and a warning is printed to stderr
- Search is case-insensitive substring match across: `InstanceID`, `Name`, `PrivateIP`, `PublicIP`
- All instance states are included in results (existing stop→start mechanism handles stopped instances)

## Capabilities

### New Capabilities

- `ec2connect-filter-by-name`: Filter EC2 instance selection via a partial name/id/ip argument passed to `ec2connect`

### Modified Capabilities

- `ec2connect`: The command gains an optional positional argument for pre-filtering; existing `-i` flag behavior is unchanged

## Impact

- `cmd/ec2connect.go`: add positional arg handling and filter logic
- No changes to `internal/awsops/ec2.go`, `cmd/instancehelper.go`, or any other file
- No breaking changes — all existing invocations (`bmc ec2connect`, `bmc ec2connect -i <id>`) continue to work as before
- Linked bean: `.beans/bmc-cjnq--connect-by-name.md`
