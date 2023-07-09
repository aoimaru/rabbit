package lib

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const NUM_OF_MIN_CHILD_CHILDREN = 0

type Node struct {
	ID       string
	Name     string
	Hash     string
	Type     string
	Children []*Node
}

type Unique struct {
	/** ID: root/ABC/123/A.pyなど Name ABC/123/A.py Hash <- treeの場合は空 Type Blob */
	ID   string
	Name string
	Hash string
	Type string
}

func GetUpperLayerID(tree *Node) string {
	tmp := strings.Split(tree.ID, "/")
	return strings.Join(tmp[:len(tmp)-1], "/")
}

func (index *Index) CreateNodes() (*Node, error) {
	var unique_layer_names []string
	var uniques []Unique

	for _, entry := range index.Entries {
		file_path := index.WorkPath + "/" + entry.Name
		if _, err := os.Stat(file_path); err != nil {
			continue
		}
		file_path_layers := strings.Split("root/"+entry.Name, "/")
		for i := 0; i <= len(file_path_layers); i++ {
			layer_name := strings.Join(file_path_layers[:i], "/")
			isUnique := true
			for _, unique_layer_name := range unique_layer_names {
				if layer_name == unique_layer_name {
					isUnique = false
					break
				}
			}
			if isUnique {
				unique_layer_names = append(unique_layer_names, layer_name)
				unique := Unique{
					ID:   layer_name,
					Name: strings.Replace(layer_name, "root/", "", 1),
				}
				layer_name = strings.Replace(layer_name, "root/", "", 1)
				// fmt.Println("layer_name:", layer_name)
				if layer_name == entry.Name {
					unique.Hash = entry.Hash
					unique.Type = "blob"
				} else {
					unique.Type = "tree"
				}
				uniques = append(uniques, unique)
			}
		}
	}
	var nodes []*Node
	for _, unique := range uniques {
		node := Node{
			ID:   unique.ID,
			Name: unique.Name,
			Type: unique.Type,
		}
		if unique.Type == "blob" {
			node.Hash = unique.Hash
		}
		// fmt.Printf("node:%+v\n", node)
		nodes = append(nodes, &node)
	}

	for _, node := range nodes {
		parent_layer_id := GetUpperLayerID(node)

		for _, parent_node := range nodes {
			if parent_node.ID == parent_layer_id {
				parent_node.Children = append(parent_node.Children, node)
			}
		}
	}
	for _, node := range nodes {
		if node.ID == "root" {
			return node, nil
		}
	}
	return nil, errors.New("fail")
}

func (c *Client) WriteTree(node *Node) string {
	if len(node.Children) <= NUM_OF_MIN_CHILD_CHILDREN {
		return node.Hash
	}

	buffer := make([]byte, 0)
	// TODO: ハードコーディング
	header := []byte{116, 114, 101, 101, 32, 51, 53, 51}
	buffer = append(buffer, header...)

	for _, child_node := range node.Children {
		child_node_buffer := make([]byte, 0)
		child_node_buffer = append(child_node_buffer, 0)
		fmt.Printf("child_node:%+v\n", child_node)
		if child_node.Type == "blob" {
			child_node_buffer = append(child_node_buffer, []byte("blob"+" ")...)
			child_node_buffer = append(child_node_buffer, []byte(child_node.Name+" ")...)
			child_node_buffer = append(child_node_buffer, []byte(child_node.Hash)...)
		} else {
			child_node_buffer = append(child_node_buffer, []byte("tree"+" ")...)
			child_node_buffer = append(child_node_buffer, []byte(child_node.Name+" ")...)
			child_node_buffer = append(child_node_buffer, []byte(c.WriteTree(child_node))...)
		}
		buffer = append(buffer, child_node_buffer...)
	}
	compressed_buffer := Compress(buffer)
	hash := CreateHash(compressed_buffer)

	object_dir := c.RepoPath + "/objects/" + hash[:2]
	object_path := c.RepoPath + "/objects/" + hash[:2] + "/" + hash[2:]
	_ = CreateDir(object_dir)
	_, _ = CreateFile(object_path, compressed_buffer)
	return hash
}

type Column struct {
	Type string
	Name string
	Hash string
	Path string
}

type Tree struct {
	Size    int
	Columns []Column
}

func (c *Client) GetTreeObject(hash string) (Tree, error) {
	// fmt.Println("Recursion")
	object_path := c.RepoPath + "/objects/" + hash[:2] + "/" + hash[2:]
	buffer, _ := GetFileBuffer(object_path)
	extracted_buffer, _ := Extract(buffer)
	lines := ToRabbitLines(extracted_buffer)
	tree := Tree{}
	for _, line := range lines[1:] {
		tmp := strings.Split(line, " ")
		column := Column{
			Type: tmp[0],
			Name: tmp[1],
			Hash: tmp[2],
			Path: c.WorkPath + "/" + tmp[1],
		}
		tree.Columns = append(tree.Columns, column)
	}
	return tree, nil
}

func (c *Client) WalkingTree(hash string, blob_columns []Column) []Column {
	tree, _ := c.GetTreeObject(hash)
	// fmt.Printf("tree:%+v\n", tree)

	for _, column := range tree.Columns {
		if column.Type == "tree" {
			blob_columns = append(blob_columns, c.WalkingTree(column.Hash, blob_columns)...)
		} else {
			blob_column := Column{Type: column.Type, Name: column.Name, Hash: column.Hash, Path: column.Path}
			blob_columns = append(blob_columns, blob_column)
		}
	}
	return blob_columns
}

func (col *Column) ToEntry() (Entry, error) {
	file_path := col.Path
	var system_call syscall.Stat_t
	syscall.Stat(file_path, &system_call)

	file_info, err := os.Stat(file_path)
	if err != nil {
		return Entry{}, err
	}

	oct := fmt.Sprintf("%o", uint32(system_call.Mode))
	mode_number, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return Entry{}, err
	}
	mode := uint32(mode_number)

	entry := Entry{
		CTime: file_info.ModTime(),
		MTime: file_info.ModTime(),
		Dev:   uint32(system_call.Dev),
		Inode: uint32(system_call.Ino),
		Mode:  mode,
		Uid:   system_call.Uid,
		Gid:   system_call.Gid,
		Size:  uint32(system_call.Size),
		Hash:  col.Hash,
		Name:  col.Name,
	}

	return entry, nil
}
