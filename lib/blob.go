package lib

import (
	"strconv"
)

type Blob struct {
	Size     int64
	Content  []byte
	RepoPath string
}

const (
	HEADER_FORMAT_INT_SIZE = 10
)

func (c *Client) CreateBlobObject(buffer []byte, size int64) (Blob, error) {
	header := []byte("blob" + " " + strconv.FormatInt(size, HEADER_FORMAT_INT_SIZE))
	header = append(header, 0)
	buffer = append(header, buffer...)
	return Blob{Size: size, Content: buffer, RepoPath: c.RepoPath}, nil
}

func (b *Blob) ToFile() (int, string, error) {
	compressed_buffer := Compress(b.Content)

	hash := CreateHash(compressed_buffer)

	object_dir := b.RepoPath + "/objects/" + hash[:2]
	object_path := b.RepoPath + "/objects/" + hash[:2] + "/" + hash[2:]

	_ = CreateDir(object_dir)

	byte_size, _ := CreateFile(object_path, compressed_buffer)

	return byte_size, hash, nil
}
