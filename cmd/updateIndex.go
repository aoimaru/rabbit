/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "updateIndex",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		hash, _ := cmd.Flags().GetString("hash")
		client := lib.CreateClient()
		if !client.IndexIsExist() {
			init_index := client.InitIndexObject()
			init_index.ToFile()
		}
		buffer, _ := lib.GetFileBuffer(client.IndexPath)
		index, _ := client.GetIndexObject(buffer)
		index, _ = index.UpdateIndex(name, hash)
		index.ToFile()
	},
}

func init() {
	rootCmd.AddCommand(updateIndexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateIndexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateIndexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateIndexCmd.Flags().StringP("name", "n", "", "set name")
	updateIndexCmd.Flags().StringP("hash", "s", "", "set hash")
}
