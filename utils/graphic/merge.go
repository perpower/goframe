// 图片合成工具
package graphic

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/perpower/goframe/funcs/convert"
)

type Merge struct{}

// Generate 图片合成
// destFile: string 目标图片地址
// bgfile: string 背景图地址
// pics: [][3]interface{fileName, X, Y} 待合成的图片,支持多图片
func (m *Merge) Generate(destFile, bgFile string, pics ...[3]interface{}) error {
	// 引入背景图片
	file, err := os.Open(bgFile)
	if err != nil {
		return err
	}
	// 图片解码
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 开始绘制
	// image.NewNRGBA(图像的大小) 这里因为要把二维码放在海报上，所以传入海报的大小
	posterImg := image.NewNRGBA(img.Bounds())

	// draw.Draw(被绘制的图片, 绘制框的大小, 要绘制的图片, 绘制的位置, 绘制类型)
	// 先把背景图画上
	draw.Draw(posterImg, posterImg.Bounds(), img, image.Pt(0, 0), draw.Over)

	// 再把需要合并的图片画上，需要注意的是坐标.
	for _, val := range pics {
		picfile, err := os.Open(convert.String(val[0]))
		if err != nil {
			continue
		}
		// 图片解码
		decimg, _, err := image.Decode(picfile)
		if err != nil {
			continue
		}
		draw.Draw(posterImg, posterImg.Bounds(), decimg, image.Pt(convert.Int(val[1]), convert.Int(val[2])), draw.Over)
	}

	// 绘制好后保存到文件中
	posterFile, _ := os.Create(destFile)
	return jpeg.Encode(posterFile, posterImg, nil)
}
