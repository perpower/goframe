// 二维码生成工具
package graphic

import (
	"image/color"

	qrcode "github.com/skip2/go-qrcode"
)

type Qrcode struct{}

// Encode 返回图片字符内容
// content: string 二维码包含的内容
// size: int 二维码图片的尺寸
// return: []byte
func (q *Qrcode) Encdoe(content string, size int) ([]byte, error) {
	return qrcode.Encode(content, qrcode.High, size)
}

// Writefile 保存为二维码图片
// content: string 二维码包含的内容
// size: int 二维码图片的尺寸
// fileName: string 文件名称
func (q *Qrcode) Writefile(content string, size int, fileName string) error {
	return qrcode.WriteFile(content, qrcode.High, size, fileName)
}

// Encode 保存为二维码图片
// content: string 二维码包含的内容
// size: int 二维码图片的尺寸
// background: color.Color 背景色
// foreground: color.Color 前景色
// fileName: string 文件名称
func (q *Qrcode) WriteColorFile(content string, size int, background color.Color, foreground color.Color, fileName string) error {
	return qrcode.WriteColorFile(content, qrcode.High, size, background, foreground, fileName)
}
