/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

type Vrfkey struct {
	CborHex string `json:"cborHex"`
}

var rootCmd = &cobra.Command{
	Use:   "cardanoleaderlogs current  -poolid xxx -blockfrostid xxx -poolticker xxx -local xxx",
	Short: "Check pool to see if scheduled for a block",
	Long: `Check pool to see if scheduled for a block
examples:
cardanoleaderlogs current  -poolid xxx -blockfrostid xxx -poolticker xxx -local xxx"

current - check current epoch
next    - check next epoch
previous - check previous epoch`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}