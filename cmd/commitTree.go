/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// commitTreeCmd represents the commitTree command
var commitTreeCmd = &cobra.Command{
	Use:   "commitTree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		hash, _ := cmd.Flags().GetString("hash")

		client := lib.CreateClient()
		commit := client.CreateCommitObject(message, hash)
		commit.ToFile()
	},
}

func init() {
	rootCmd.AddCommand(commitTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitTreeCmd.Flags().StringP("message", "m", "", "set file message")
	commitTreeCmd.Flags().StringP("hash", "s", "", "set file hash")
}
