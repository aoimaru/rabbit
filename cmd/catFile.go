/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "catFile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.CreateClient()
		// tree, _ := client.GetTreeObject("cbddda25bf86bb669870dd9a0a740e443a26d20c")
		// fmt.Printf("%+v\n", tree)
		hash := "797e800b94ac0b595db4b8259d88ab23addc3fee"
		init_columns := make([]lib.Column, 0)
		blob_columns := client.WalkingTree(hash, init_columns)
		buffer, _ := lib.GetFileBuffer(client.IndexPath)
		index, _ := client.GetIndexObject(buffer)
		roll_back_index := index.RollBackIndex(blob_columns)
		// roll_back_index.ToFile()
		roll_back_index.RollBackWorkingTree()
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
