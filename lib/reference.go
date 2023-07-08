package lib

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Head interface {
}

type DetachedHead struct {
	Head string
}

type TatchedHead struct {
	Head string
}

func (c *Client) GetHeadHash() (string, error) {
	buffer, _ := GetFileBuffer(c.HeadPath)
	reference := string(buffer)
	re := regexp.MustCompile(`ref: refs/heads/(\w+)`)
	if re.MatchString(reference) {
		reference = strings.Replace(reference, "\n", "", -1)
		reference = strings.Replace(reference, "ref: ", "", 1)
		reference = strings.Replace(reference, ":", "", 1)
		tatched_reference_path := c.RepoPath + reference
		tatched_buffer, _ := GetFileBuffer(tatched_reference_path)
		return string(tatched_buffer), nil
	} else {
		return reference, nil
	}
}

func (c *Client) GetHeadRef() string {
	buffer, _ := GetFileBuffer(c.HeadPath)
	reference := string(buffer)
	re := regexp.MustCompile(`ref: refs/heads/(\w+)`)
	if re.MatchString(reference) {
		reference = strings.Replace(reference, "\n", "", -1)
		reference = strings.Replace(reference, "ref: ", "", 1)
		reference = strings.Replace(reference, ":", "", 1)
		return reference
	} else {
		return reference
	}

}

func (c *Client) UpdateRef(refs string, hash string) error {
	re := regexp.MustCompile(`refs/heads/(\w+)`)
	if !re.MatchString(refs) {
		return errors.New("Invalid refs path")
	}
	refs_path := c.RepoPath + "/" + refs
	if _, err := os.Stat(refs_path); err != nil {
		return err
	}
	_, _ = CreateFile(refs_path, []byte(hash))

	return nil
}
