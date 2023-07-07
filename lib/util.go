package lib

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

const (
	PADDING_BYTE_TO_UINT64_IS_8 = 8
	PADDING_BYTE_TO_UINT32_IS_4 = 4
	PADDING_BYTE_TO_UINT16_IS_2 = 2
)

func BufferToUint64(buffer []byte) uint64 {
	padding := make([]byte, PADDING_BYTE_TO_UINT64_IS_8-len(buffer))
	source := append(padding, buffer...)
	artifact := binary.BigEndian.Uint64(source)
	return artifact
}

func BufferToUint32(buffer []byte) uint32 {
	padding := make([]byte, PADDING_BYTE_TO_UINT32_IS_4-len(buffer))
	source := append(padding, buffer...)
	artifact := binary.BigEndian.Uint32(source)
	return artifact

}

func BufferToUint16(buffer []byte) uint16 {
	padding := make([]byte, PADDING_BYTE_TO_UINT16_IS_2-len(buffer))
	source := append(padding, buffer...)
	artifact := binary.BigEndian.Uint16(source)
	return artifact
}

func BufferToUnixTimeStamp(buffer []byte) (time.Time, error) {
	padding := make([]byte, PADDING_BYTE_TO_UINT64_IS_8-len(buffer))
	source := append(padding, buffer...)
	artifact := binary.BigEndian.Uint64(source)
	times := int64(artifact)
	var offsetHour, offsetMinute int
	if _, err := fmt.Sscanf("+0900", "+%02d%02d", &offsetHour, &offsetMinute); err != nil {
		return time.Time{}, err
	}
	location := time.FixedZone(" ", 3600*offsetHour+60*offsetMinute)
	timestamp := time.Unix(times, 0).In(location)
	return timestamp, nil
}

func BufferToMode(buffer []byte) (uint32, error) {
	dec := binary.BigEndian.Uint32(buffer)
	oct := fmt.Sprintf("%o", dec)
	num, err := strconv.ParseUint(oct, 10, 32)
	if err != nil {
		return 0, err
	}
	mode := uint32(num)
	return mode, nil
}

func GetPaddingSize(had uint64) uint64 {
	Rem := had % 8
	return 8 - Rem
}

func EntryFieldToBuffer(entry_field uint32) []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, entry_field)
	return buffer
}
