package cmd

import (
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/ui"
)

// selectInstanceID presents an interactive table of instances and returns the
// selected instance's InstanceId. cols determines which columns are shown and
// in what order. The InstanceId is resolved by finding the "InstanceId" column
// index; if not present, index 0 is used as fallback.
func selectInstanceID(instances []awsops.Instance, cols []string) (string, error) {
	rows := awsops.InstanceRows(instances, cols)
	row, err := ui.SelectFromTable(cols, rows)
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
