package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ec2findCmd = &cobra.Command{
	Use:   "ec2find <search-string>",
	Short: "Find EC2 instances across all profiles in a group",
	Args:  cobra.ExactArgs(1),
	RunE:  runEC2Find,
}

func init() {
	rootCmd.AddCommand(ec2findCmd)
}

func runEC2Find(cmd *cobra.Command, args []string) error {
	searchString := strings.ToLower(args[0])
	if searchString == "" {
		return fmt.Errorf("search string cannot be empty")
	}

	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return err
	}

	groups := awsconfig.Groups(profiles)
	if len(groups) == 0 {
		return fmt.Errorf("no profile groups found")
	}

	groupItems := make([]ui.Item, len(groups))
	for i, g := range groups {
		groupItems[i] = ui.Item{Title: g}
	}

	selectedGroup, err := ui.Choose("Select AWS account group to search in", groupItems)
	if err != nil {
		return err
	}
	if selectedGroup == "" {
		return nil
	}

	groupProfiles := awsconfig.ByGroup(profiles, selectedGroup)
	profileNames := make([]string, len(groupProfiles))
	for i, p := range groupProfiles {
		profileNames[i] = p.Name
	}

	fmt.Fprintf(os.Stderr, "Searching for %q across %d profiles...\n", args[0], len(profileNames))
	allInstances := awsops.ListInstancesForProfiles(profileNames)

	// Filter by search string
	var matches []awsops.Instance
	for _, inst := range allInstances {
		combined := strings.ToLower(inst.InstanceID + inst.Name + inst.PrivateIP + inst.PublicIP + inst.Profile)
		if strings.Contains(combined, searchString) {
			matches = append(matches, inst)
		}
	}

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "No instances found matching %q in group %q\n", args[0], selectedGroup)
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cols := cfg.EC2.Columns
	hasProfile := false
	for _, c := range cols {
		if c == "Profile" {
			hasProfile = true
			break
		}
	}
	if !hasProfile {
		cols = append(cols, "Profile")
	}

	return ui.ShowTable(cols, awsops.InstanceRows(matches, cols))
}
