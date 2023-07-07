package lib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Entry struct {
	CTime time.Time
	MTime time.Time
	Dev   uint32
	Inode uint32
	Mode  uint32
	Uid   uint32
	Gid   uint32
	Size  uint32
	Hash  string
	Name  string
}

type Index struct {
	Dirc     string
	Version  uint32
	Number   uint32
	Entries  []Entry
	Path     string
	WorkPath string
}

func (c *Client) InitIndexObject() Index {
	var index Index
	index.Dirc = "DIRC"
	index.Version = 2
	index.Number = uint32(0)
	index.Path = c.IndexPath
	index.WorkPath = c.WorkPath

	return index
}

func (c *Client) GetIndexObject(buffer []byte) (Index, error) {

	dirc := string(buffer[0:4])
	if dirc != "DIRC" {
		return Index{}, errors.New("Is Not IndexFile")
	}
	version := BufferToUint32(buffer[4:8])
	if version != 2 {
		return Index{}, errors.New("Invalid Version")
	}

	number_of_entry := BufferToUint32(buffer[8:12])

	var index Index

	index.Dirc = dirc
	index.Version = version
	index.Number = number_of_entry
	index.Path = c.IndexPath
	index.WorkPath = c.WorkPath

	buffer = buffer[12:]

	var current_entry_number uint32
	current_entry_number = 0

	for {
		if current_entry_number >= number_of_entry {
			break
		}
		current_entry_number++

		c_time, _ := BufferToUnixTimeStamp(buffer[0:4])
		m_time, _ := BufferToUnixTimeStamp(buffer[8:12])
		dev := BufferToUint32(buffer[16:20])
		inode := BufferToUint32(buffer[20:24])
		mode, _ := BufferToMode(buffer[24:28])
		uid := BufferToUint32(buffer[28:32])
		gid := BufferToUint32(buffer[32:36])
		size := BufferToUint32(buffer[36:40])
		hash := hex.EncodeToString(buffer[40:60])
		name_size := BufferToUint64(buffer[60:62])
		name := string(buffer[62 : 62+name_size])

		entry := Entry{
			CTime: c_time,
			MTime: m_time,
			Dev:   dev,
			Inode: inode,
			Mode:  mode,
			Uid:   uid,
			Gid:   gid,
			Size:  size,
			Hash:  hash,
			Name:  name,
		}
		index.Entries = append(index.Entries, entry)

		padding := GetPaddingSize(62 + name_size)
		offset := 62 + name_size + padding
		buffer = buffer[offset:]
	}
	return index, nil
}

func (index *Index) UpdateIndex(name string, hash string) (Index, error) {
	file_path := index.WorkPath + "/" + name
	// fmt.Println("<---debug::lib/UpdateIndex::file_path--->")
	// fmt.Println("index.WorkPath:", index.WorkPath)
	// fmt.Println("file_path:", file_path)
	var system_call syscall.Stat_t
	syscall.Stat(file_path, &system_call)

	file_info, err := os.Stat(file_path)
	if err != nil {
		return Index{}, err
	}

	oct := fmt.Sprintf("%o", uint32(system_call.Mode))
	mode_number, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return Index{}, err
	}
	mode := uint32(mode_number)

	new_entry := Entry{
		CTime: file_info.ModTime(),
		MTime: file_info.ModTime(),
		Dev:   uint32(system_call.Dev),
		Inode: uint32(system_call.Ino),
		Mode:  mode,
		Uid:   system_call.Uid,
		Gid:   system_call.Gid,
		Size:  uint32(system_call.Size),
		Hash:  hash,
		Name:  name,
	}

	fmt.Println("<---debug::lib/index.go::UpdateIndex::145--->")
	fmt.Println("new_entry:", new_entry)

	var new_index Index
	new_index.WorkPath = index.WorkPath
	new_index.Path = index.Path

	for _, entry := range index.Entries {
		if entry.Name == name {
			continue
		}
		if entry.Hash == hash {
			continue
		}
		fmt.Println("old:", entry)
		new_index.Entries = append(new_index.Entries, entry)
	}
	new_index.Entries = append(new_index.Entries, new_entry)

	new_index.Dirc = "DIRC"
	new_index.Version = 2
	new_index.Number = uint32(len(new_index.Entries))

	fmt.Println("<---debug::lib/index.go::UpdateIndex::167--->")
	fmt.Println("new_index:", new_index)

	return new_index, nil
}

func (index *Index) ToFile() error {

	// fmt.Println("<---debug::lib/ToFile--->")
	// fmt.Printf("%+v\n", index)

	buffer := make([]byte, 0)
	dirc := []byte(index.Dirc)
	version := EntryFieldToBuffer(index.Version)
	number := EntryFieldToBuffer(index.Number)

	// fmt.Println("number:", number)

	buffer = append(buffer, dirc...)
	buffer = append(buffer, version...)
	buffer = append(buffer, number...)

	for _, entry := range index.Entries {
		// fmt.Println("entry:", entry)

		c_time_buffer := EntryFieldToBuffer(uint32(entry.CTime.Unix()))
		buffer = append(buffer, c_time_buffer...)
		buffer = append(buffer, c_time_buffer...)

		m_time_buffer := EntryFieldToBuffer(uint32(entry.MTime.Unix()))
		buffer = append(buffer, m_time_buffer...)
		buffer = append(buffer, m_time_buffer...)

		dev_buffer := EntryFieldToBuffer(entry.Dev)
		buffer = append(buffer, dev_buffer...)

		inode_buffer := EntryFieldToBuffer(entry.Inode)
		buffer = append(buffer, inode_buffer...)

		mode_buffer := EntryFieldToBuffer(entry.Mode)
		buffer = append(buffer, mode_buffer...)

		uid_buffer := EntryFieldToBuffer(entry.Uid)
		buffer = append(buffer, uid_buffer...)

		gid_buffer := EntryFieldToBuffer(entry.Gid)
		buffer = append(buffer, gid_buffer...)

		size_buffer := EntryFieldToBuffer(entry.Size)
		buffer = append(buffer, size_buffer...)

		hash_buffer, _ := hex.DecodeString(entry.Hash)
		buffer = append(buffer, hash_buffer...)

		name_size_buffer := make([]byte, 2)
		length_of_name := len(entry.Name)
		binary.BigEndian.PutUint16(name_size_buffer, uint16(length_of_name))
		buffer = append(buffer, name_size_buffer...)

		name_buffer := []byte(entry.Name)
		buffer = append(buffer, name_buffer...)

		var name_padding uint64
		name_padding = 62

		length_of_name_buffer := len(name_buffer)
		padding_size := GetPaddingSize(name_padding + uint64(length_of_name_buffer))
		padding_buffer := make([]byte, padding_size)
		buffer = append(buffer, padding_buffer...)

	}

	// fmt.Println("buffer:", buffer)

	index_path := index.Path
	_, err := CreateFile(index_path, buffer)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
