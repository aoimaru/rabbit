/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "writeTree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.CreateClient()
		index_buffer, _ := lib.GetFileBuffer(client.IndexPath)
		index, _ := client.GetIndexObject(index_buffer)
		node, _ := index.CreateNodes()
		hash := client.WriteTree(node)
		fmt.Println("hash:", hash)

	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
