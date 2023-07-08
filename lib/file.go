package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetFileBuffer(file_path string) ([]byte, error) {
	f, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		return buffer, err
	}

	return buffer, nil
}

func GetFileSize(file_path string) (int64, error) {
	f, err := os.Stat(file_path)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

func (c *Client) HashToPath(hash string) (string, string) {
	object_dir := c.RepoPath + "/objects/" + hash[:2]
	object_path := c.RepoPath + "/objects/" + hash[:2] + "/" + hash[2:]
	return object_dir, object_path
}

func Compress(buffer []byte) []byte {
	var compressed bytes.Buffer
	zlib_writer := zlib.NewWriter(&compressed)
	zlib_writer.Write(buffer)
	zlib_writer.Close()
	return compressed.Bytes()
}

func CreateHash(buffer []byte) string {
	sha1 := sha1.New()
	sha1.Write(buffer)
	return hex.EncodeToString(sha1.Sum(nil))
}

func CreateDir(dir_path string) error {
	if _, err := os.Stat(dir_path); err != nil {
		if err := os.MkdirAll(dir_path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CreateFile(file_path string, buffer []byte) (int, error) {
	w, err := os.Create(file_path)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	byte_size, err := w.Write(buffer)
	if err != nil {
		return 0, err
	}
	return byte_size, nil
}

func (c *Client) WalkingDir() ([]string, error) {
	if _, err := os.Stat(c.WorkPath); err != nil {
		return nil, err
	}
	paths := make([]string, 0)
	err := filepath.Walk(c.WorkPath, func(path string, info os.FileInfo, err error) error {
		rel_path, err := filepath.Rel(c.WorkPath, path)
		if info.IsDir() {
			if strings.HasPrefix(rel_path, ".rabbit") {
				return filepath.SkipDir
			}
			return nil
		}
		paths = append(paths, rel_path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil

}
