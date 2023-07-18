package lib

import (
	"log"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	currentDir := "/Users/haradakanon/Desktop/Stub"
	client := CreateClient(currentDir)

	err := client.Init()
	if err != nil {
		t.Error("Failed to init")
	}

	t.Run(".rabbit dir", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit"
		assertDir(t, path)
	})

	t.Run("/objects dir", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit/objects"
		assertDir(t, path)
	})

	t.Run("/refs dir", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit/refs"
		assertDir(t, path)
	})

	t.Run("/refs/heads dir", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit/refs/heads"
		assertDir(t, path)
	})

	t.Run("HEAD file", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit/HEAD"
		want := "ref: refs/heads/main\n"
		assertFileContents(t, path, want)
	})

	t.Run("/refs/heads/main file", func(t *testing.T) {
		path := "/Users/haradakanon/Desktop/Stub/.rabbit/refs/heads/main"
		want := ""
		assertFileContents(t, path, want)
	})
}

// assertDir はディレクトリの存在をチェックする.
func assertDir(t *testing.T, path string) {
	if _, err := os.Stat(path); err != nil {
		t.Errorf("The directory does not exist: %s", path)
	}
}

// assertFileContents はファイルの中身をチェックする.
func assertFileContents(t *testing.T, path, want string) {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	got := string(contents)
	if got != want {
		t.Errorf("contents %s want %s", got, want)
	}
}
