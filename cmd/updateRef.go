/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// updateRefCmd represents the updateRef command
var updateRefCmd = &cobra.Command{
	Use:   "updateRef",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		refs, _ := cmd.Flags().GetString("refs")
		hash, _ := cmd.Flags().GetString("hash")
		client := lib.CreateClient()
		_ = client.UpdateRef(refs, hash)
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateRefCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateRefCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateRefCmd.Flags().StringP("refs", "r", "", "set refs")
	updateRefCmd.Flags().StringP("hash", "s", "", "set hash")
}
