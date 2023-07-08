/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

const NUM_OF_ADD_ARGS = 1

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		client := lib.CreateClient()
		if !client.IndexIsExist() {
			init_index := client.InitIndexObject()
			init_index.ToFile()
		}

		if name == "." {
			working_file_paths, _ := client.WalkingDir()
			index_buffer, _ := lib.GetFileBuffer(client.IndexPath)
			index, _ := client.GetIndexObject(index_buffer)
			for _, working_file_path := range working_file_paths {
				file_buffer, _ := lib.GetFileBuffer(working_file_path)
				size, _ := lib.GetFileSize(working_file_path)
				blob, _ := client.CreateBlobObject(file_buffer, size)
				_, hash, _ := blob.ToFile()
				index, _ = index.UpdateIndex(working_file_path, hash)
			}
			index.ToFile()
		} else {
			file_buffer, _ := lib.GetFileBuffer(name)
			size, _ := lib.GetFileSize(name)
			blob, _ := client.CreateBlobObject(file_buffer, size)
			_, hash, _ := blob.ToFile()
			index_buffer, _ := lib.GetFileBuffer(client.IndexPath)
			index, _ := client.GetIndexObject(index_buffer)
			index, _ = index.UpdateIndex(name, hash)
			index.ToFile()
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		/** 引数のバリデーションを行うことができる */
		if len(args) < NUM_OF_ADD_ARGS {
			return errors.New("requires args")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
