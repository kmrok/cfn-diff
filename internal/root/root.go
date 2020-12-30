package root

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"

	"github.com/kmrok/cfn-diff/internal/changeset"
	"github.com/kmrok/cfn-diff/internal/drift"
)

func Execute() error {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		cmd.Println(err)
		return err
	}
	return nil
}

func newCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:           "cfn-diff",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	svc := cloudformation.New(session.Must(session.NewSession()))
	cmd.AddCommand(drift.NewCommand(svc), changeset.NewCommand(svc))

	return cmd
}
