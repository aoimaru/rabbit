/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		client := lib.CreateClient()
		index_buffer, _ := lib.GetFileBuffer(client.IndexPath)
		index, _ := client.GetIndexObject(index_buffer)
		node, _ := index.CreateNodes()
		hash := client.WriteTree(node)
		commit := client.CreateCommitObject(message, hash)
		commit.ToFile()
		refs := client.GetHeadRef()
		err := client.UpdateRef(refs, hash)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringP("message", "m", "", "set file message")
}
