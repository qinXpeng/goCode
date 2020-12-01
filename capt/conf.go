package capt

import (
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。

func ReadFonts(dirPath, suffix string) error {
	file := make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	pathSqp := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			file = append(file, dirPath+pathSqp+fi.Name())
		}
	}
	SetFontFamily(file...)
	return nil
}

func SetFontFamily(fontPath ...string) {

	FontFamily = append(FontFamily, fontPath...)
}

//获取所及字体.
func RandFontFamily() (*truetype.Font, error) {
	fontfile := FontFamily[r.Intn(len(FontFamily))]

	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return &truetype.Font{}, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return &truetype.Font{}, err
	}
	return f, nil
}

//随机生成深色系.
func RandDeepColor() color.RGBA {

	randColor := RandColor()

	increase := float64(30 + r.Intn(255))

	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))

	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//随机生成浅色.
func RandLightColor() color.RGBA {

	red := r.Intn(55) + 200
	green := r.Intn(55) + 200
	blue := r.Intn(55) + 200

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//生成随机颜色.
func RandColor() color.RGBA {

	red := r.Intn(255)
	green := r.Intn(255)
	blue := r.Intn(255)
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//生成随机字体.
func RandText(num int) string {
	textNum := len(txtChars)
	text := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		text = text + string(txtChars[r.Intn(textNum)])
	}
	return text
}

// 颜色代码转换为RGB
//input int
//output int red, green, blue.
func ColorToRGB(colorVal int) color.RGBA {

	red := colorVal >> 16
	green := (colorVal & 0x00FF00) >> 8
	blue := colorVal & 0x0000FF

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(255),
	}
}
