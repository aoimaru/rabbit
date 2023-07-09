package lib

type RabbitObject interface {
}

func (c *Client) GetRabbitHeader(buffer []byte) string {
	// 責務 bufferを受け取り, headerを文字列で返す(tree, commit, blob)
	header := make([]byte, 1024)
	for _, buf := range buffer {
		if buf == 0 {
			break
		}
		header = append(header, buf)
	}
	return string(header)
}

func ToRabbitLines(buffer []byte) []string {
	/**　ここのcolumnは構造体Columnとは無関係*/
	columns := make([]string, 0)
	column := make([]byte, 0)

	for _, buf := range buffer {
		if buf == 0 {
			columns = append(columns, string(column))
			column = make([]byte, 0)
		} else {
			column = append(column, buf)
		}
	}
	columns = append(columns, string(column))
	return columns
}
