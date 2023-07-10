/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/aoimaru/rabbit/lib"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		branch_name := args[0]
		client := lib.CreateClient()

		// 指定してブランチのリファレンスから, そのブランチのハッシュを取得
		hash, err := client.GetBranchHash(branch_name)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 上記のコミットオブジェクトから最新のコミットに紐づいているtreeオブジェクトのハッシュを取得
		commit := client.GetCommitObject(hash)
		fmt.Printf("commit:%+v\n", commit)
		tree_hash := commit.Tree

		// ツリーオブジェクトのハッシュから紐づいているblobオブジェクトのリストを作成(blob_columns)
		init_columns := make([]lib.Column, 0)
		blob_columns := client.WalkingTree(tree_hash, init_columns)

		// インデックスファイルを取得して, オブジェクトに変換
		buffer, _ := lib.GetFileBuffer(client.IndexPath)
		index, _ := client.GetIndexObject(buffer)

		// コミットtree(blobオブジェクト)に基づいてインデックスを更新
		roll_back_index := index.RollBackIndex(blob_columns)
		roll_back_index.ToFile()

		// あるブランチでは, 管理していないファイルがあった場合は削除を行う
		working_paths, _ := client.WalkingDir()
		roll_back_index.RollBackWorkingTree(working_paths)

		// HEADを指定したブランチに書き換える
		err = client.SwitchRef(branch_name)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
