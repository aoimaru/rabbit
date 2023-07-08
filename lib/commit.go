package lib

import (
	"fmt"
	"time"
)

type Parent struct {
	Hash string
}

type Sign struct {
	Name      string
	Email     string
	TimeStamp time.Time
}

type Commit struct {
	Size      int
	Tree      string
	Parents   []Parent
	Author    Sign
	Committer Sign
	Message   string
	RepoPath  string
}

func (c *Client) CreateCommitObject(message string, hash string) Commit {
	var commit Commit
	commit.Size = 119
	commit.Tree = hash

	var sign Sign
	sign.Name = "aoi nakamura"
	sign.Email = "hello@world.com"
	sign.TimeStamp = time.Now()

	commit.Author = sign
	commit.Committer = sign

	commit.Message = message

	parent_hash, _ := c.GetHeadHash()
	commit.Parents = append(commit.Parents, Parent{Hash: parent_hash})
	commit.RepoPath = c.RepoPath

	return commit

}

func (commit *Commit) ToFile() {
	buffer := make([]byte, 0)
	buffer = append(buffer, []byte("commit 199")...)
	buffer = append(buffer, 0)

	for _, parent := range commit.Parents {
		parent_string := "parent " + parent.Hash
		buffer = append(buffer, []byte(parent_string)...)
		buffer = append(buffer, 0)
	}

	tree_string := "tree " + commit.Tree
	buffer = append(buffer, []byte(tree_string)...)
	buffer = append(buffer, 0)

	author_string := "author " + commit.Author.Name + " " + commit.Author.Email + " " + commit.Author.TimeStamp.String()
	buffer = append(buffer, []byte(author_string)...)
	buffer = append(buffer, 0)

	committer_string := "committer " + commit.Committer.Name + " " + commit.Committer.Email + " " + commit.Committer.TimeStamp.String()
	buffer = append(buffer, []byte(committer_string)...)
	buffer = append(buffer, 0)

	buffer = append(buffer, []byte(commit.Message)...)
	fmt.Println(buffer)

	compressed_buffer := Compress(buffer)

	hash := CreateHash(compressed_buffer)

	object_dir := commit.RepoPath + "/objects/" + hash[:2]
	object_path := commit.RepoPath + "/objects/" + hash[:2] + "/" + hash[2:]

	_ = CreateDir(object_dir)

	_, _ = CreateFile(object_path, compressed_buffer)

}
