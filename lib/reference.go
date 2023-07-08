package lib

import (
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

func (c *Client) GetHeadRef() (string, error) {
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
