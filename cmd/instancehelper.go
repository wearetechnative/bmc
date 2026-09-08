package cmd

import (
	"strings"

	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/ui"
)

// selectInstanceID presents an interactive table of instances and returns the
// selected instance's InstanceId. cols determines which columns are shown and
// in what order. The InstanceId is resolved by finding the "InstanceId" column
// index; if not present, index 0 is used as fallback.
//
// The picker supports real-time type-to-filter: each row's filter key is built
// from the instance's ID, name and IPs (the same fields as the positional
// [search] fragment), so filtering works even when those columns are hidden.
func selectInstanceID(instances []awsops.Instance, cols []string) (string, error) {
	rows := awsops.InstanceRows(instances, cols)
	filterKeys := make([]string, len(instances))
	for i, inst := range instances {
		filterKeys[i] = strings.ToLower(inst.InstanceID + " " + inst.Name + " " + inst.PrivateIP + " " + inst.PublicIP)
	}
	row, err := ui.SelectFromTable(cols, rows, filterKeys)
	if err != nil || row == nil {
		return "", err
	}

	idx := 0
	for i, col := range cols {
		if col == "InstanceId" {
			idx = i
			break
		}
	}
	return row[idx], nil
}
