package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsconfig"
	"golang.org/x/term"
	"github.com/wearetechnative/bmc/internal/config"
	"github.com/wearetechnative/bmc/internal/mfa"
	"github.com/wearetechnative/bmc/internal/ui"
)

var (
	profselPreferred string
	profselList      bool
	profselJSON      bool
)

var profselCmd = &cobra.Command{
	Use:   "profsel",
	Short: "Set AWS_PROFILE by evaluating this command's output",
	Long: `Select an AWS profile interactively.

To set AWS_PROFILE in your current shell, use the shell wrapper:
  eval "$(bmc profsel)"

Or install the wrapper permanently:
  bmc install-shell-integration`,
	RunE: runProfsel,
}

func init() {
	profselCmd.Flags().StringVarP(&profselPreferred, "profile", "p", "", "Pre-select a profile by name")
	profselCmd.Flags().BoolVarP(&profselList, "list", "l", false, "List all profiles in tabular format")
	profselCmd.Flags().BoolVar(&profselJSON, "json", false, "Output JSON {source_profile, profile_name, profile_arn}")
	rootCmd.AddCommand(profselCmd)
}

func runProfsel(cmd *cobra.Command, args []string) error {
	if profselList && profselJSON {
		return fmt.Errorf("flags -l and --json are incompatible")
	}

	profiles, err := awsconfig.LoadProfiles()
	if err != nil {
		return err
	}

	if profselList {
		return printProfiles(profiles)
	}

	var selectedProfile awsconfig.Profile

	if profselPreferred != "" {
		p, ok := awsconfig.FindProfile(profiles, profselPreferred)
		if !ok {
			return fmt.Errorf("profile %q not found", profselPreferred)
		}
		selectedProfile = p
	} else {
		selectedProfile, err = selectProfileInteractive(profiles)
		if err != nil {
			return err
		}
		if selectedProfile.Name == "" {
			if profselJSON {
				fmt.Println(`{"error": "no profile selected"}`)
				os.Exit(1)
			}
			return nil
		}
	}

	sourceProfile, err := awsconfig.ResolveSourceProfile(selectedProfile)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := mfa.EnsureValid(sourceProfile, cfg, os.Stderr); err != nil {
		return err
	}

	if profselJSON {
		out := map[string]string{
			"source_profile": sourceProfile,
			"profile_name":   selectedProfile.Name,
			"profile_arn":    selectedProfile.RoleARN,
		}
		data, _ := json.Marshal(out)
		fmt.Println(string(data))
		return nil
	}

	// Normal output: eval-able export statement
	fmt.Printf("export AWS_PROFILE=%s\n", selectedProfile.Name)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprintln(os.Stderr, "Tip: run 'bmc install-shell-integration' to set AWS_PROFILE automatically")
	}
	return nil
}

func selectProfileInteractive(profiles []awsconfig.Profile) (awsconfig.Profile, error) {
	groups := awsconfig.Groups(profiles)
	if len(groups) == 0 {
		return awsconfig.Profile{}, fmt.Errorf("no profile groups found in ~/.aws/config")
	}

	groupItems := make([]ui.Item, len(groups))
	for i, g := range groups {
		groupItems[i] = ui.Item{Title: g}
	}

	for {
		selectedGroup, err := ui.Choose("Select AWS account group", groupItems)
		if err != nil {
			return awsconfig.Profile{}, err
		}
		if selectedGroup == "" {
			return awsconfig.Profile{}, nil
		}

		groupProfiles := awsconfig.ByGroup(profiles, selectedGroup)
		profileItems := make([]ui.Item, len(groupProfiles))
		for i, p := range groupProfiles {
			desc := p.AccountID
			if p.RoleName != "" {
				desc += " / " + p.RoleName
			}
			profileItems[i] = ui.Item{Title: p.Name, Desc: desc}
		}

		selectedName, err := ui.Choose("Select profile  (group: "+selectedGroup+")", profileItems)
		if errors.Is(err, ui.ErrBack) {
			continue
		}
		if err != nil {
			return awsconfig.Profile{}, err
		}
		if selectedName == "" {
			return awsconfig.Profile{}, nil
		}

		p, _ := awsconfig.FindProfile(profiles, selectedName)
		return p, nil
	}
}

func printProfiles(profiles []awsconfig.Profile) error {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Group\tName\tARN Number")
	for _, p := range profiles {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Group, p.Name, p.AccountID)
	}
	return w.Flush()
}
