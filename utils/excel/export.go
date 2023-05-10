package excel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// maxCharCount 最多26个字符A-Z
const maxCharCount = 26

// SaveAsFile 保存为指定Excel文件
// f: *excelize.File
// fileName: string 文件路径（包含文件名）
func SaveAsFile(f *excelize.File, fileName string) (err error) {
	err = f.SaveAs(fileName)
	return err
}

// SaveAsWeb 保存为Excel流媒体文件
// fileName: string 文件名
func SaveAsWeb(f *excelize.File, fileName string, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/vnd.ms-excel") // 需要让浏览器识别为excel文件
	w.Header().Set("Content-Type", "application/octet-stream") // 默认让浏览器下载文件
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	if err := f.Write(w); err != nil {
		http.Error(w, "导出失败", http.StatusInternalServerError)
	}
}

// ExportExcel 导出Excel文件
// sheetName: string 工作表名称, 注意这里不要取sheet1这种名字,否则导致文件打开时发生部分错误。
// cols: []string 列名切片， 表头
// rows: [][]interface{} 数据切片，是一个二维数组
func ExportExcel(sheetName string, cols []string, rows [][]interface{}) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if sheetName == "" {
		sheetName = "Sheet1" // 设置默认的工作簿名称
	}
	sheetIndex, _ := f.NewSheet(sheetName)
	maxColumnRowNameLen := 1 + len(strconv.Itoa(len(rows)))
	columnCount := len(cols)
	if columnCount > maxCharCount {
		maxColumnRowNameLen++
	} else if columnCount > maxCharCount*maxCharCount {
		maxColumnRowNameLen += 2
	}
	columnNames := make([][]byte, 0, columnCount)
	for i, header := range cols {
		columnName := getColumnName(i, maxColumnRowNameLen)
		columnNames = append(columnNames, columnName)
		// 初始化excel表头，这里的index从1开始要注意
		curColumnName := getColumnRowName(columnName, 1)
		err := f.SetCellValue(sheetName, curColumnName, header)
		if err != nil {
			return nil, err
		}
	}
	for rowIndex, row := range rows {
		for columnIndex, columnName := range columnNames {
			// 从第二行开始
			err := f.SetCellValue(sheetName, getColumnRowName(columnName, rowIndex+2), row[columnIndex])
			if err != nil {
				return nil, err
			}
		}
	}
	f.SetActiveSheet(sheetIndex)
	return f, nil
}

// getColumnName 生成列名
// Excel的列名规则是从A-Z往后排;超过Z以后用两个字母表示，比如AA,AB,AC;两个字母不够以后用三个字母表示，比如AAA,AAB,AAC
// 这里做数字到列名的映射：0 -> A, 1 -> B, 2 -> C
// maxColumnRowNameLen 表示名称框的最大长度，假设数据是10行，1000列，则最后一个名称框是J1000(如果有表头，则是J1001),是4位
// 这里根据 maxColumnRowNameLen 生成切片，后面生成名称框的时候可以复用这个切片，而无需扩容
func getColumnName(column, maxColumnRowNameLen int) []byte {
	const A = 'A'
	if column < maxCharCount {
		// 第一次就分配好切片的容量
		slice := make([]byte, 0, maxColumnRowNameLen)
		return append(slice, byte(A+column))
	} else {
		// 递归生成类似AA,AB,AAA,AAB这种形式的列名
		return append(getColumnName(column/maxCharCount-1, maxColumnRowNameLen), byte(A+column%maxCharCount))
	}
}

// getColumnRowName 生成名称框
// Excel的名称框是用A1,A2,B1,B2来表示的，这里需要传入前一步生成的列名切片，然后直接加上行索引来生成名称框，就无需每次分配内存
func getColumnRowName(columnName []byte, rowIndex int) (columnRowName string) {
	l := len(columnName)
	columnName = strconv.AppendInt(columnName, int64(rowIndex), 10)
	columnRowName = string(columnName)
	// 将列名恢复回去
	columnName = columnName[:l]
	return
}
