package capt

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype"
	"github.com/qinXpeng/goCode/random"
	"golang.org/x/image/font"
)

const txtChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	dpi                 = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	r                   = rand.New(rand.NewSource(time.Now().UnixNano()))
	FontFamily []string = make([]string, 0)
)

const (
	//图片格式
	ImageFormatPng = iota
	ImageFormatJpeg
	ImageFormatGif
	//验证码噪点强度
	CaptchaComplexLower = iota
	CaptchaComplexMedium
	CaptchaComplexHigh
)

type CaptchaImage struct {
	nrgba   *image.NRGBA
	width   int
	height  int
	Complex int
}

func NewCaptchaImage(width, height int, bgColor color.RGBA) (*CaptchaImage, error) {
	im := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(im, im.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)
	return &CaptchaImage{
		nrgba:  im,
		width:  width,
		height: height,
	}, nil
}

func (this *CaptchaImage) SaveImage(w io.Writer, imageFormat int) error {

	switch imageFormat {
	case ImageFormatPng:
		return png.Encode(w, this.nrgba)
	case ImageFormatJpeg:
		return jpeg.Encode(w, this.nrgba, &jpeg.Options{Quality: 100})
	case ImageFormatGif:
		return gif.Encode(w, this.nrgba, &gif.Options{NumColors: 256})
	}
	return errors.New("Not support image format")
}

// 画粗的白条
func (this *CaptchaImage) DrawHollowLine() *CaptchaImage {
	first := this.width / 20
	end := first * 19
	lineColor := color.RGBA{
		R: 245,
		G: 250,
		B: 251,
		A: 255,
	}
	x1 := float64(r.Intn(first))
	x2 := float64(r.Intn(first) + end)
	multi := float64(r.Intn(5)+3) / float64(5)
	if int(multi*10)%3 == 0 {
		multi = multi * -1.0
	}
	w := this.height / 20
	for ; x1 < x2; x1++ {

		y := math.Sin(x1*math.Pi*multi/float64(this.width)) * float64(this.height/3)

		if multi < 0 {
			y = y + float64(this.height/2)
		}
		this.nrgba.Set(int(x1), int(y), lineColor)

		for i := 0; i <= w; i++ {
			this.nrgba.Set(int(x1), int(y)+i, lineColor)
		}
	}
	return this
}
func (captcha *CaptchaImage) DrawSineLine() *CaptchaImage {
	px := 0
	var py float64 = 0
	//振幅
	a := r.Intn(captcha.height / 2)
	//Y轴方向偏移量
	b := random.Random(int64(-captcha.height/4), int64(captcha.height/4))
	//X轴方向偏移量
	f := random.Random(int64(-captcha.height/4), int64(captcha.height/4))
	// 周期
	var t float64 = 0
	if captcha.height > captcha.width/2 {
		t = random.Random(int64(captcha.width/2), int64(captcha.height))
	} else {
		t = random.Random(int64(captcha.height), int64(captcha.width/2))
	}
	w := float64((2 * math.Pi) / t)
	// 曲线横坐标起始位置
	px1 := 0
	px2 := int(random.Random(int64(float64(captcha.width)*0.8), int64(captcha.width)))
	c := color.RGBA{R: uint8(r.Intn(150)), G: uint8(r.Intn(150)), B: uint8(r.Intn(150)), A: uint8(255)}
	for px = px1; px < px2; px++ {
		if w != 0 {
			py = float64(a)*math.Sin(w*float64(px)+f) + b + (float64(captcha.width) / float64(5))
			i := captcha.height / 5
			for i > 0 {
				captcha.nrgba.Set(px+i, int(py), c)
				i--
			}
		}
	}
	return captcha
}
func (captcha *CaptchaImage) Drawline(num int) *CaptchaImage {

	first := (captcha.width / 10)
	end := first * 9

	y := captcha.height / 3

	for i := 0; i < num; i++ {

		point1 := Point{X: r.Intn(first), Y: r.Intn(y)}
		point2 := Point{X: r.Intn(first) + end, Y: r.Intn(y)}

		if i%2 == 0 {
			point1.Y = r.Intn(y) + y*2
			point2.Y = r.Intn(y)
		} else {
			point1.Y = r.Intn(y) + y*(i%2)
			point2.Y = r.Intn(y) + y*2
		}

		captcha.drawBeeline(point1, point2, RandDeepColor())

	}
	return captcha
}
func (captcha *CaptchaImage) drawBeeline(point1, point2 Point, lineColor color.RGBA) {
	dx := math.Abs(float64(point1.X - point2.X))

	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy
	for {
		captcha.nrgba.Set(point1.X, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-1, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X+2, point1.Y, lineColor)
		captcha.nrgba.Set(point1.X-2, point1.Y, lineColor)
		if point1.X == point2.X && point1.Y == point2.Y {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

//画边框.
func (captcha *CaptchaImage) DrawBorder(borderColor color.RGBA) *CaptchaImage {
	for x := 0; x < captcha.width; x++ {
		captcha.nrgba.Set(x, 0, borderColor)
		captcha.nrgba.Set(x, captcha.height-1, borderColor)
	}
	for y := 0; y < captcha.height; y++ {
		captcha.nrgba.Set(0, y, borderColor)
		captcha.nrgba.Set(captcha.width-1, y, borderColor)
	}
	return captcha
}

//画噪点.
func (captcha *CaptchaImage) DrawNoise(complex int) *CaptchaImage {
	density := 18
	if complex == CaptchaComplexLower {
		density = 28
	} else if complex == CaptchaComplexMedium {
		density = 18
	} else if complex == CaptchaComplexHigh {
		density = 8
	}
	maxSize := (captcha.height * captcha.width) / density

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		captcha.nrgba.Set(rw, rh, RandColor())
		size := r.Intn(maxSize)
		if size%3 == 0 {
			captcha.nrgba.Set(rw+1, rh+1, RandColor())
		}
	}
	return captcha
}

//画文字噪点.
func (captcha *CaptchaImage) DrawTextNoise(complex int) error {
	density := 1500
	if complex == CaptchaComplexLower {
		density = 2000
	} else if complex == CaptchaComplexMedium {
		density = 1500
	} else if complex == CaptchaComplexHigh {
		density = 1000
	}

	maxSize := (captcha.height * captcha.width) / density

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)
	rawFontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(10))

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		text := RandText(1)
		fontSize := rawFontSize/2 + float64(r.Intn(5))

		c.SetSrc(image.NewUniform(RandLightColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			log.Println(err)
			return err
		}
		c.SetFont(f)
		pt := freetype.Pt(rw, rh)

		_, err = c.DrawString(text, pt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//写字.
func (captcha *CaptchaImage) DrawText(text string) error {
	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.nrgba.Bounds())
	c.SetDst(captcha.nrgba)
	c.SetHinting(font.HintingFull)

	fontWidth := captcha.width / len(text)

	for i, s := range text {

		fontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(9))

		c.SetSrc(image.NewUniform(RandDeepColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			log.Println(err)
			return err
		}
		c.SetFont(f)

		x := int(fontWidth)*i + int(fontWidth)/int(fontSize)

		y := 5 + r.Intn(captcha.height/2) + int(fontSize/2)

		pt := freetype.Pt(x, y)

		_, err = c.DrawString(string(s), pt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
