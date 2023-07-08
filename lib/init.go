package lib

import (
	"fmt"
	"os"
)

func (c *Client) Init() error {
	if _, err := os.Stat(c.RepoPath); err == nil {
		_ = os.RemoveAll(c.RepoPath)
	}
	if err := os.MkdirAll(c.RepoPath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(c.RepoPath+"/objects", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(c.RepoPath+"/refs", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(c.RepoPath+"/refs/heads", os.ModePerm); err != nil {
		return err
	}
	empty_buffer := []byte("ref: ref/heads/master\n")
	_, err := CreateFile(c.HeadPath, empty_buffer)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
