package lib

func CreateClient(currentDir string) Client {
	return Client{
		WorkPath:  currentDir,
		RepoPath:  currentDir + "/.rabbit",
		IndexPath: currentDir + "/.rabbit/index",
		HeadPath:  currentDir + "/.rabbit/HEAD",
	}
}

type Client struct {
	WorkPath  string
	RepoPath  string
	IndexPath string
	HeadPath  string
}

func (c *Client) GetWorkPath() string {
	return c.WorkPath
}

func (c *Client) GetRepoPath() string {
	return c.RepoPath
}
