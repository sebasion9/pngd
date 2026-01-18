package main

import (
	"fmt"
	"os"
	"pngd/cli"
	dc "pngd/decoder"

	"github.com/spf13/cobra"
)

func main() {
	cmd, opts := cli.NewRootCmd()
	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := run(opts)
		if err != nil {
			fmt.Println("[!] pngd encountered an error")
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}


// example usage
func run(opts *cli.Options) error {
	bytes, err := os.ReadFile(opts.Input)
	if err != nil {
		return err
	}


	decoder := dc.NewDecoder(bytes)
	if err = decoder.ValidateSignature(); err != nil {
		return err
	}

	// decoded, err := decoder.Decode()
	_, err = decoder.Decode()
	if err != nil {
		return err
	}

	decoder.Info()
	for _, w := range decoder.Warnings() {
		fmt.Println(w)
	}

	return nil
}


