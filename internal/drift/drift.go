package drift

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kmrok/cfn-diff/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCommand(cfn *cloudformation.CloudFormation) *cobra.Command {
	var configFilePath string
	var cmd = &cobra.Command{
		Use:           "drift",
		Short:         "detecting drift on a stack",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			listStacksOutput, err := cfn.ListStacks(nil)
			if err != nil {
				return fmt.Errorf("ListStacksError: %w", err)
			}

			stacks := make([]string, 0, len(listStacksOutput.StackSummaries))
			for _, v := range listStacksOutput.StackSummaries {
				if *v.StackStatus == cloudformation.StackStatusDeleteComplete {
					continue
				}
				stacks = append(stacks, *v.StackName)
			}

			cfg, err := config.Load(configFilePath)
			if err != nil {
				return err
			}

			var isDetected bool
			stacks = excludeConfigStacks(stacks, cfg)
			for _, stack := range stacks {
				status, err := describeDriftDetectionStatus(cfn, stack)
				if err != nil {
					return err
				}

				cmd.Println(stack, status)
				if status == cloudformation.StackDriftStatusDrifted {
					describeStackResourceDriftsOutput, err := cfn.DescribeStackResourceDrifts(&cloudformation.DescribeStackResourceDriftsInput{
						StackName: aws.String(stack),
						StackResourceDriftStatusFilters: []*string{
							aws.String(cloudformation.StackResourceDriftStatusModified),
							aws.String(cloudformation.StackResourceDriftStatusDeleted),
						},
					})
					if err != nil {
						return err
					}
					cmd.Println(color.New(color.FgYellow).Sprint(describeStackResourceDriftsOutput))
					isDetected = true
				}
			}

			if isDetected && cfg.Run.EnableCIMode {
				return fmt.Errorf(color.New(color.FgRed).Sprint("detected drift on stacks"))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFilePath, "config", "f", "", "Configuration file")

	return cmd
}

func excludeConfigStacks(stacks []string, cfg config.Config) []string {
	shouldExclude := func(stackName string) bool {
		for _, pattern := range cfg.StackWithDriftDetection {
			re := regexp.MustCompile(pattern)
			if !re.MatchString(stackName) {
				return true
			}
		}
		for _, pattern := range cfg.StackWithoutDriftDetection {
			re := regexp.MustCompile(pattern)
			if re.MatchString(stackName) {
				return true
			}
		}
		return false
	}

	s := make([]string, 0, len(stacks))
	for _, stack := range stacks {
		if shouldExclude(stack) {
			continue
		}
		s = append(s, stack)
	}

	return s
}

func describeDriftDetectionStatus(cfn *cloudformation.CloudFormation, stackName string) (string, error) {
	detectStackDriftOutput, err := cfn.DetectStackDrift(&cloudformation.DetectStackDriftInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return "", fmt.Errorf("DetectStackDriftError: %w", err)
	}

	var status string
	for {
		time.Sleep(3 * time.Second)

		describeStackDriftDetectionStatusOutput, err := cfn.DescribeStackDriftDetectionStatus(&cloudformation.DescribeStackDriftDetectionStatusInput{
			StackDriftDetectionId: detectStackDriftOutput.StackDriftDetectionId,
		})
		if err != nil {
			return "", fmt.Errorf("DescribeStackDriftDetectionStatusError: %w", err)
		}

		if *describeStackDriftDetectionStatusOutput.DetectionStatus != cloudformation.StackDriftDetectionStatusDetectionInProgress {
			status = *describeStackDriftDetectionStatusOutput.StackDriftStatus
			break
		}
	}

	return status, nil
}
