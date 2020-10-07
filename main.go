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

var teksFilter []string

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a TEK export binary file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		teks, err := DecodeFromFile(args[0])
		if err != nil {
			return err
		}

		filtered := make([]*TemporaryExposureKey, 0)
		for _, tek := range teks {
			if len(teksFilter) == 0 {
				tek.RPIs = nil
				filtered = append(filtered, tek)
			}

			for _, tf := range teksFilter {
				if tek.ID.ToBase64() == tf {
					filtered = append(filtered, tek)
				}
			}
		}

		// print something
		b, err := json.MarshalIndent(filtered, "", "    ")
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

	decodeCmd.Flags().StringArrayVarP(
		&teksFilter,
		"tek",
		"t",
		make([]string, 0),
		"Display the RPIs (Rolling Proximity Identifiers) of the specified TEKs",
	)

	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.Execute()
}
