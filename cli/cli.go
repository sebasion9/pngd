package cli

import "github.com/spf13/cobra"


type Options struct {
	Input string
	Output string
}

func NewRootCmd() (*cobra.Command, *Options) {
	opts := &Options{}
	cmd := &cobra.Command {
		Use: "pngd",
		Short: "png decoder",
	}

	cmd.Flags().StringVarP(&opts.Input, "input", "i", "", "source file")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "output file")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")

	return cmd, opts

}

