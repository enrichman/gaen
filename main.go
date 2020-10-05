package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gaen/export"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var rootCmd = &cobra.Command{
	Use:   "gaen",
	Short: "gaen is a cli to interact with the Google Apple Exposure Notification",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the gaen version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a TEK export binary file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return decode(args[0])
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

func decode(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	in = in[16:]

	teke := &export.TemporaryExposureKeyExport{}
	if err := proto.Unmarshal(in, teke); err != nil {
		return err
	}

	for _, tek := range teke.Keys {
		fmt.Printf("\nTEK: [%s] - [%s]\n", base64.StdEncoding.EncodeToString(tek.KeyData), encodeToHexString(tek.KeyData))

		rpis := DecodeFromTEK(tek)
		for _, rpi := range rpis {
			b, _ := json.MarshalIndent(rpi, "", "\t")
			fmt.Printf("\nRPI:\n%v\n", string(b))
			break
		}
		break
	}

	return nil
}
