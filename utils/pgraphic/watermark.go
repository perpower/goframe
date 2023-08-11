// 图片加水印
package pgraphic

import (
	"image"

	"github.com/disintegration/imaging"
)

var WaterMark = gwaterMark{}

type gwaterMark struct{}

// 图片式水印
// sourceImg: string 源图片地址
// waterImg: string 水印图片地址
// destImgs: *string 目标图片地址
func (w *gwaterMark) CreatePicMark(sourceImg, waterImg string, destImg *string) error {
	// Open the original image file
	originalImage, err := imaging.Open(sourceImg, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	// Load the watermark image
	watermarkImage, err := imaging.Open(waterImg)
	if err != nil {
		return err
	}

	// Resize the watermark image to match the size of the original image
	resizedWatermark := imaging.Resize(watermarkImage, originalImage.Bounds().Dx(), originalImage.Bounds().Dy(), imaging.Lanczos)

	// Add the watermark to the original image
	result := imaging.Overlay(originalImage, resizedWatermark, image.Pt(0, 0), 1)

	// Save the result to a new file
	err = imaging.Save(result, *destImg)
	if err != nil {
		return err
	}
	return nil
}
