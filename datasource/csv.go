package datasource

import (
	"encoding/csv"
	"os"
)

type csvHandler struct {
}

func (c *csvHandler) ReadFile(fileName string) []map[string]any {

	// 打开CSV文件
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建CSV阅读器
	reader := csv.NewReader(file)

	// 读取第一行作为标题行
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}

	// 读取所有CSV内容
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// 创建一个slice，用于存储map
	var dataSlice []map[string]interface{}

	// 遍历CSV的每一行
	for _, record := range records {
		// 创建一个map来存储每一行的数据
		rowData := make(map[string]interface{})
		for i, value := range record {
			// 使用headers作为key
			rowData[headers[i]] = value
		}
		// 将map添加到slice中
		dataSlice = append(dataSlice, rowData)
	}

	return dataSlice
}
