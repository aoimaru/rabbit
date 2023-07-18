package lib

import (
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	currentDir, _ := os.Getwd()
	client := CreateClient(currentDir)

	t.Run("local path", func(t *testing.T) {
		want := "/Users/haradakanon/Desktop/rabbit/lib"
		assertLocalPath(t, client, want)
	})

	t.Run("local repository", func(t *testing.T) {
		want := "/Users/haradakanon/Desktop/rabbit/lib/.rabbit"
		assertLocalRepo(t, client, want)
	})
}

// assertLocalPath は実行者のローカルPCの絶対パスと一致しているかチェックする.
func assertLocalPath(t *testing.T, got Client, want string) {
	t.Helper()

	if got.GetWorkPath() != want {
		t.Errorf("Unexpected WorkPath value: %s", got.WorkPath)
	}
}

// assertLocalRepo は実行者のローカルPCの絶対パスと一致しているかチェックする.
func assertLocalRepo(t *testing.T, got Client, want string) {
	t.Helper()

	if got.GetRepoPath() != want {
		t.Errorf("Unexpected WorkPath value: %s", got.RepoPath)
	}
}
