package changeset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
		Use:           "changeset",
		Short:         "describing change set",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configFilePath)
			if err != nil {
				return err
			}

			var isDetected bool
			for _, v := range cfg.StackTemplateMaps {
				describeStacksOutput, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
					StackName: &v.StackName,
				})
				if err != nil {
					return fmt.Errorf("DescribeStacksError: %w", err)
				}

				var parameters []*cloudformation.Parameter
				for _, stack := range describeStacksOutput.Stacks {
					for _, parameter := range stack.Parameters {
						parameters = append(parameters, &cloudformation.Parameter{
							ParameterKey:     parameter.ParameterKey,
							UsePreviousValue: aws.Bool(true),
						})
					}
				}

				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("GetwdError: %w", err)
				}
				f, err := os.Open(path.Join(dir, v.TemplateName))
				if err != nil {
					return err
				}
				defer f.Close()
				b, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}

				changeSetName := fmt.Sprintf("%s-%s", v.StackName, time.Now().Format("20060102150405"))
				if _, err := cfn.CreateChangeSet(&cloudformation.CreateChangeSetInput{
					StackName:     aws.String(v.StackName),
					ChangeSetName: aws.String(changeSetName),
					TemplateBody:  aws.String(string(b)),
					Parameters:    parameters,
					Capabilities:  []*string{aws.String(cloudformation.CapabilityCapabilityNamedIam)},
				}); err != nil {
					return fmt.Errorf("CreateChangeSetError: %w", err)
				}

				var describeChangeSetOutput *cloudformation.DescribeChangeSetOutput
				for {
					describeChangeSetOutput, err = cfn.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
						StackName:     aws.String(v.StackName),
						ChangeSetName: aws.String(changeSetName),
					})
					if err != nil {
						return fmt.Errorf("DescribeChangeSetError: %w", err)
					}

					if *describeChangeSetOutput.Status != cloudformation.ChangeSetStatusCreatePending &&
						*describeChangeSetOutput.Status != cloudformation.ChangeSetStatusCreateInProgress {
						break
					}

					if err := cfn.WaitUntilChangeSetCreateComplete(&cloudformation.DescribeChangeSetInput{
						StackName:     aws.String(v.StackName),
						ChangeSetName: aws.String(changeSetName),
					}); err != nil {
						return fmt.Errorf("WaitUntilChangeSetCreateCompleteError: %w", err)
					}
				}

				if len(describeChangeSetOutput.Changes) == 0 {
					cmd.Println(v.StackName, "NO_CHANGES_DETECTED")
				} else {
					isDetected = true
					cmd.Println(v.StackName, "CHANGES_DETECTED")
					cmd.Println(color.New(color.FgYellow).Sprint(describeChangeSetOutput.Changes))
				}

				if _, err := cfn.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{
					StackName:     aws.String(v.StackName),
					ChangeSetName: aws.String(changeSetName),
				}); err != nil {
					return fmt.Errorf("DeleteChangeSetError: %w", err)
				}
			}

			if isDetected && cfg.Run.EnableCIMode {
				return fmt.Errorf(color.New(color.FgRed).Sprint("detected diff between templates and resources"))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFilePath, "config", "f", "", "Configuration file")

	return cmd
}
