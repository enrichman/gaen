package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

var rootCmd = &cobra.Command{
	Use:   "gaen",
	Short: "gaen is a cli to interact with the Google Apple Exposure Notification",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the gaen version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gaen version " + version)
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a TEK export binary file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		teks, err := DecodeFromFile(args[0])
		if err != nil {
			return err
		}

		// print something
		b, err := json.MarshalIndent(teks, "", "    ")
		if err != nil {
			return err
		}

		fmt.Printf("%s", string(b))
		return nil
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a TEK export binary file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Download("out", args[0])
	},
}

func main() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.Execute()
}
