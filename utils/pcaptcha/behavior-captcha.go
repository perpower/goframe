// 行为式验证码
package pcaptcha

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/wenlng/go-captcha/captcha"
	"golang.org/x/image/font"
)

type CaptchaInfo struct {
	Dots             map[int]captcha.CharDot `json:"dots"`             // 图片位置数据
	ImageBase64      string                  `json:"imageBase64"`      // 主图base64
	ThumbImageBase64 string                  `json:"thumbImageBase64"` // 缩略图base64
	UniqKey          string                  `json:"uniqKey"`          // 验证码key
}

// 验证码图片尺寸配置
type ImageSize struct {
	Image captcha.Size // 主图尺寸 默认：300 * 300
	Thumb captcha.Size // 缩略图尺寸 默认：150 * 40
}

type gBehavior struct {
	catpacha *captcha.Captcha
}

var (
	Behavior = gBehavior{}
	// 设置验证码文本随机种子
	defaultChars = []string{"爱", "心", "梦", "情", "意", "思", "行", "动", "路", "歌", "舞", "风", "雨", "春", "夏", "秋", "冬", "日", "月", "星", "云", "水", "山", "田", "土", "金", "银", "钱", "宝", "贝", "叶", "花", "树", "果", "桃", "李", "柳", "竹", "草", "虫", "鸟", "兽", "鱼", "蛇", "龙", "虎", "牛", "羊", "马", "狗", "猫", "鼠", "人", "佛", "道", "德", "义", "仁", "智", "信", "善", "礼", "慈", "孝", "敬", "诚", "勇", "冠", "帅", "相", "侯", "伯", "卿", "子", "男", "女", "声", "影", "光", "彩", "红", "绿", "蓝", "紫", "丹", "橙", "黄", "灰", "白", "黑", "和", "平", "开", "快", "慢", "远", "近", "大", "小", "高", "低", "上", "下", "左", "右"}
)

// 实例化验证码生成器
// chars: []string 验证码文本随机种子,推荐设置 50+ 以上的随机种子,
// 每个单词不能超出2个字符，超出会影响位置的验证 ,不设置用默认的。
func (g *gBehavior) Instance(chars ...[]string) {
	// 单例模式的验证码实例
	capt := captcha.GetCaptcha()

	if len(chars) > 0 {
		capt.SetRangChars(chars[0])
	} else {
		capt.SetRangChars(defaultChars)
	}

	g.catpacha = capt
}

// 生成验证码数据
// size: ImageSzie  验证码图片尺寸, 不设置采用默认
func (g *gBehavior) Generate(size ...ImageSize) (*CaptchaInfo, error) {
	capt := g.catpacha
	capt = g.setFont(capt)

	if len(size) > 0 {
		if !reflect.DeepEqual(size[0].Image, captcha.Size{}) {
			capt = g.setImage(capt, size[0].Image)
		}
		if !reflect.DeepEqual(size[0].Thumb, captcha.Size{}) {
			capt = g.setThumb(capt, size[0].Thumb)
		}
	} else {
		capt = g.setImage(capt)
		capt = g.setThumb(capt)
	}

	// 生成验证码
	dots, imageBase64, thumbImageBase64, key, err := capt.Generate()
	if err != nil {
		return nil, err
	}

	return &CaptchaInfo{
		Dots:             dots,
		ImageBase64:      imageBase64,
		ThumbImageBase64: thumbImageBase64,
		UniqKey:          key,
	}, nil

}

// 校验验证码
// dots: verifyDots 基准点位置数据
// verify: []string 校验点位置数据
// return: bool
func (g *gBehavior) CheckCaptcha(verifyDots []string, dots map[int]captcha.CharDot) bool {
	checkRes := false
	if (len(dots) * 2) == len(verifyDots) {
		for i, dot := range dots {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", verifyDots[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", verifyDots[k]), 64)

			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			checkRes = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
			if !checkRes {
				break
			}
		}
	}

	return checkRes
}

// 设置字体
func (g *gBehavior) setFont(capt *captcha.Captcha) *captcha.Captcha {
	capt.SetFont([]string{
		"assets/fonts/fzshengsksjw_cu.ttf",
		"assets/fonts/fzssksxl.ttf",
		"assets/fonts/hyrunyuan.ttf",
	})

	return capt
}

// 主图配置
// size: captcha.Size  图片尺寸
func (g *gBehavior) setImage(capt *captcha.Captcha, size ...captcha.Size) *captcha.Captcha {
	// Method: SetBackground(color []string);
	// Desc: 设置验证码背景图，自动仅读取一次并加载到内存中缓存，如需重置可清除缓存
	// ====================================================
	capt.SetBackground([]string{
		"assets/images/1.jpg",
		"assets/images/2.jpg",
	})

	// ====================================================
	// Method: SetImageSize(size Size);
	// Desc: 设置验证码主图的尺寸
	// ====================================================
	if len(size) > 0 && !reflect.DeepEqual(size[0], captcha.Size{}) {
		capt.SetImageSize(size[0])
	} else {
		capt.SetImageSize(captcha.Size{Width: 300, Height: 300})
	}

	// ====================================================
	// Method: SetImageQuality(val int);
	// Desc: 设置验证码主图清晰度，压缩级别范围 QualityCompressLevel1 - 5，QualityCompressNone无压缩，默认为最低压缩级别
	// ====================================================
	capt.SetImageQuality(captcha.QualityCompressNone)

	// ====================================================
	// Method: SetFontHinting(val font.Hinting);
	// Desc: 设置字体Hinting值 (HintingNone,HintingVertical,HintingFull)
	// ====================================================
	capt.SetFontHinting(font.HintingFull)

	// ====================================================
	// Method: SetTextRangLen(val captcha.RangeVal);
	// Desc: 设置验证码文本显示的总数随机范围
	// ====================================================
	capt.SetTextRangLen(captcha.RangeVal{Min: 6, Max: 7})

	// ====================================================
	// Method: SetRangFontSize(val captcha.RangeVal);
	// Desc: 设置验证码文本的随机大小
	// ====================================================
	capt.SetRangFontSize(captcha.RangeVal{Min: 32, Max: 42})

	// ====================================================
	// Method: SetTextRangFontColors(colors []string);
	// Desc: 设置验证码文本的随机十六进制颜色
	// ====================================================
	capt.SetTextRangFontColors([]string{
		"#1d3f84",
		"#3a6a1e",
	})

	// ====================================================
	// Method: SetImageFontAlpha(val float64);
	// Desc:设置验证码字体的透明度
	// ====================================================
	capt.SetImageFontAlpha(0.5)

	// ====================================================
	// Method: SetTextShadow(val bool);
	// Desc:设置字体阴影
	// ====================================================
	capt.SetTextShadow(true)

	// ====================================================
	// Method: SetTextShadowColor(val string);
	// Desc:设置字体阴影颜色
	// ====================================================
	capt.SetTextShadowColor("#101010")

	// ====================================================
	// Method: SetTextShadowPoint(val captcha.Point);
	// Desc:设置字体阴影偏移位置
	// ====================================================
	capt.SetTextShadowPoint(captcha.Point{X: 1, Y: 1})

	// ====================================================
	// Method: SetTextRangAnglePos(pos []captcha.RangeVal);
	// Desc:设置验证码文本的旋转角度
	// ====================================================
	capt.SetTextRangAnglePos([]captcha.RangeVal{
		{Min: 1, Max: 15},
		{Min: 15, Max: 30},
		{Min: 30, Max: 45},
		{Min: 315, Max: 330},
		{Min: 330, Max: 345},
		{Min: 345, Max: 359},
	})

	// ====================================================
	// Method: SetImageFontDistort(val int);
	// Desc:设置验证码字体的扭曲程度
	// ====================================================
	capt.SetImageFontDistort(captcha.DistortLevel2)

	return capt
}

// 缩略图配置
// size: captcha.Size  图片尺寸
func (g *gBehavior) setThumb(capt *captcha.Captcha, size ...captcha.Size) *captcha.Captcha {
	// Method: SetThumbSize(size Size);
	// Desc: 设置缩略图的尺寸
	// ====================================================
	if len(size) > 0 && !reflect.DeepEqual(size[0], captcha.Size{}) {
		capt.SetThumbSize(size[0])
	} else {
		capt.SetThumbSize(captcha.Size{Width: 150, Height: 40})
	}

	// ====================================================
	// Method: SetRangCheckTextLen(val captcha.RangeVal);
	// Desc:设置缩略图校验文本的随机长度范围
	// ====================================================
	capt.SetRangCheckTextLen(captcha.RangeVal{Min: 2, Max: 4})

	// ====================================================
	// Method: SetRangCheckFontSize(val captcha.RangeVal);
	// Desc:设置缩略图校验文本的随机大小
	// ====================================================
	capt.SetRangCheckFontSize(captcha.RangeVal{Min: 24, Max: 30})

	// ====================================================
	// Method: SetThumbTextRangFontColors(colors []string);
	// Desc: 设置缩略图文本的随机十六进制颜色
	// ====================================================
	capt.SetThumbTextRangFontColors([]string{
		"#1d3f84",
		"#3a6a1e",
	})

	// ====================================================
	// Method: SetThumbBgColors(colors []string);
	// Desc: 设置缩略图的背景随机十六进制颜色
	// ====================================================
	capt.SetThumbBgColors([]string{
		"#1d3f84",
		"#3a6a1e",
	})

	// ====================================================
	// Method: SetThumbBackground(colors []string);
	// Desc:设置缩略图的随机图像背景，自动仅读取一次并加载到内存中缓存，如需重置可清除缓存
	// ====================================================
	capt.SetThumbBackground([]string{
		"assets/images/thumb/r1.jpg",
		"assets/images/thumb/r2.jpg",
	})

	// ====================================================
	// Method: SetThumbBgDistort(val int);
	// Desc:设置缩略图背景的扭曲程度
	// ====================================================
	capt.SetThumbBgDistort(captcha.DistortLevel2)

	// ====================================================
	// Method: SetThumbFontDistort(val int);
	// Desc:设置缩略图字体的扭曲程度
	// ====================================================
	capt.SetThumbFontDistort(captcha.DistortLevel2)

	// ====================================================
	// Method: SetThumbBgCirclesNum(val int);
	// Desc:设置缩略图背景的圈点数
	// ====================================================
	capt.SetThumbBgCirclesNum(20)

	// ====================================================
	// Method: SetThumbBgSlimLineNum(val int);
	// Desc:设置缩略图背景的线条数
	// ====================================================
	capt.SetThumbBgSlimLineNum(3)

	return capt
}
