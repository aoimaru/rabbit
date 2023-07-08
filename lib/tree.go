package lib

import (
	"errors"
	"os"
	"strings"
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
