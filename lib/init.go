package lib

import (
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

	return nil
}
