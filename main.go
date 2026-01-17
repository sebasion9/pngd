package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"pngd/cli"
	"pngd/decoder"

	"github.com/spf13/cobra"
	"golang.org/x/image/bmp"
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


func run(opts *cli.Options) error {
	bytes, err := os.ReadFile(opts.Input)
	if err != nil {
		return err
	}


	decoder := decoder.NewDecoder(bytes)
	if err = decoder.ValidateSignature(); err != nil {
		return err
	}

	if err = decoder.Decode(); err != nil {
		return err
	}

	for _, w := range decoder.Warnings {
		fmt.Println(w)
	}

	return nil
}

