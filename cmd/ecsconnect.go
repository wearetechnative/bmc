package cmd

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wearetechnative/bmc/internal/awsops"
	"github.com/wearetechnative/bmc/internal/prereqs"
	"github.com/wearetechnative/bmc/internal/ui"
)

var ecsconnectCmd = &cobra.Command{
	Use:   "ecsconnect",
	Short: "Shell connect to an ECS container",
	RunE:  runECSConnect,
}

func init() {
	ecsconnectCmd.Flags().StringVarP(&globalProfile, "profile", "p", "", "AWS profile to use (omit value to force interactive selection)")
	ecsconnectCmd.Flags().Lookup("profile").NoOptDefVal = " "
	rootCmd.AddCommand(ecsconnectCmd)
}

func runECSConnect(cmd *cobra.Command, args []string) error {
	// Check prerequisites before starting TUI
	if err := prereqs.Check(prereqs.AWSCLI); err != nil {
		return err
	}
	if err := prereqs.Check(prereqs.SessionManagerPlugin); err != nil {
		return err
	}

	profile, err := ensureAWSProfile()
	if err != nil {
		return err
	}

	// Fetch clusters once.
	clusters, err := awsops.ListClusters(profile)
	if err != nil {
		return fmt.Errorf("error listing clusters: %w", err)
	}
	if len(clusters) == 0 {
		return fmt.Errorf("no clusters found in current region")
	}
	clusterItems := make([]ui.Item, len(clusters))
	for i, c := range clusters {
		clusterItems[i] = ui.Item{Title: c.Name}
	}

clusterLoop:
	for {
		cluster, err := ui.Choose("Select ECS cluster", clusterItems)
		if err != nil {
			return err
		}
		if cluster == "" {
			return nil
		}

		// Fetch services for selected cluster.
		services, err := awsops.ListServices(profile, cluster)
		if err != nil {
			return fmt.Errorf("error listing services: %w", err)
		}
		serviceItems := make([]ui.Item, len(services))
		for i, s := range services {
			serviceItems[i] = ui.Item{Title: s.Name}
		}

		for {
			service, err := ui.Choose("Select service ("+cluster+")", serviceItems)
			if errors.Is(err, ui.ErrBack) {
				continue clusterLoop
			}
			if err != nil {
				return err
			}
			if service == "" {
				return nil
			}

			// Fetch tasks for selected service.
			tasks, err := awsops.ListRunningTasks(profile, cluster, service)
			if err != nil {
				return fmt.Errorf("error listing tasks: %w", err)
			}
			taskItems := make([]ui.Item, len(tasks))
			for i, t := range tasks {
				taskItems[i] = ui.Item{Title: t.ShortID}
			}

			for {
				taskShortID, err := ui.Choose("Select task ("+cluster+" > "+service+")", taskItems)
				if errors.Is(err, ui.ErrBack) {
					break
				}
				if err != nil {
					return err
				}
				if taskShortID == "" {
					return nil
				}

				var taskARN string
				for _, t := range tasks {
					if t.ShortID == taskShortID {
						taskARN = t.ARN
						break
					}
				}

				// Fetch containers for selected task.
				containers, err := awsops.ListContainers(profile, cluster, taskARN)
				if err != nil {
					return fmt.Errorf("error listing containers: %w", err)
				}
				containerItems := make([]ui.Item, len(containers))
				for i, c := range containers {
					containerItems[i] = ui.Item{Title: c.Name}
				}

				for {
					container, err := ui.Choose("Select container ("+cluster+" > "+service+" > "+taskShortID+")", containerItems)
					if errors.Is(err, ui.ErrBack) {
						break
					}
					if err != nil {
						return err
					}
					if container == "" {
						return nil
					}

					awsBin := prereqs.FindPath("aws")
					fmt.Fprintf(os.Stderr, "-- Connecting to %s > %s > %s > %s\n", cluster, service, taskShortID, container)
					return syscall.Exec(awsBin, []string{
						"aws", "ecs", "execute-command",
						"--cluster", cluster,
						"--interactive",
						"--container", container,
						"--command", "/bin/sh",
						"--task", taskARN,
					}, os.Environ())
				}
			}
		}
	}
}
