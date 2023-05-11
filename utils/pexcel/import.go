// excel导入功能组件
package pexcel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

var Import = gimport{}

type gimport struct{}

// ImportExcel 导入Excel文件, 大多数场景下是从form表单接收文件
// file: io.Reader 文件
// skip: int 指定要跳过的标题行数
// return:  [][]string  将结果以一个二维数组返回
func (i *gimport) ImportExcel(file io.Reader, skip int) (arr [][]string, err error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return [][]string{}, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 默认读取第一个Sheet内容
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	for irow, row := range rows {
		if irow > (skip - 1) { // 跳过标题行
			var data []string
			arr = append(arr, append(data, row...))
		}
	}

	return arr, err
}
