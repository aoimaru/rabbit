package lib

import (
	"os"
)

func (c *Client) Init() error {
	if _, err := os.Stat(c.RepoPath); err == nil {
		_ = os.RemoveAll(c.RepoPath)
	}

	// init時に生成するファイル名の配列.
	initArray := []string{"", "/objects", "/refs", "/refs/heads"}
	for _, value := range initArray {
		err := CreateDir(c.RepoPath + value)
		if err != nil {
			return err
		}
	}

	// TODO ここもまとめますか?
	initHeadBuffer := []byte("ref: refs/heads/main\n")
	_, _ = CreateFile(c.HeadPath, initHeadBuffer)

	initRefsHeadsMainBuffer := []byte("")
	_, _ = CreateFile(c.RepoPath+"/refs/heads/main", initRefsHeadsMainBuffer)

	return nil
}
