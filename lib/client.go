package lib

import "os"

func CreateClient() Client {
	curent_dir, _ := os.Getwd()
	return Client{
		WorkPath:  curent_dir,
		RepoPath:  curent_dir + "/.rabbit",
		IndexPath: curent_dir + "/.rabbit/index",
	}
}

type Client struct {
	WorkPath  string
	RepoPath  string
	IndexPath string
}

func (c *Client) GetWorkPath() string {
	return c.WorkPath
}

func (c *Client) GetRepoPath() string {
	return c.RepoPath
}
