package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand"
	"strconv"
)

type Img struct {
	Code   string
	Height int
	Width  int
	// 计算点的大小
	r int
	p *image.Paletted
}

// 数字写入图片
func (img *Img) drawFont(index string, sx, sy int) {
	val, ok := font[index]
	if !ok {
		return
	}

	for i, v := range val { //draw point
		if v < 1 {
			//zero fill
			continue
		}
		ax := (i % fontWidth) * img.r
		ay := (i / fontWidth) * img.r
		for ri := 1; ri <= img.r; ri++ {
			if (ri % 2) == 0 {
				img.p.SetColorIndex(sx+ax, sy+ay+ri, 1)
			} else {
				for i := 1; i < img.r; i++ {
					if (i % 2) == 0 {
						img.p.SetColorIndex(sx+ax-i, sy+ay+ri, 1)
					}
				}
			}
		}
	}
}

// 设置图片随机扭曲
func (img *Img) distort() {
	mixedx := 2.0 * math.Pi / (rand.Float64()*100 + 100)
	mixedz := rand.Float64()*5 + 5

	newP := image.NewPaletted(image.Rect(0, 0, img.Width, img.Height), img.p.Palette)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			xo := mixedz * math.Sin(float64(y)*mixedx)
			yo := mixedz * math.Cos(float64(x)*mixedx)
			newP.SetColorIndex(x, y, img.p.ColorIndexAt(x+int(xo), y+int(yo)))
		}
	}
	img.p = newP
}

// 增加干扰线
func (img *Img) through() {

	amplitude := rand.Float64()*15 + 7

	y := rand.Intn(img.Height-((img.Height/3)*2)) + img.Height/3
	dx := (math.Pi * 2) / (rand.Float64()*100 + 80)

	for x := 0; x < img.Width; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < img.r; yn++ {
			r := rand.Intn(img.r)
			img.p.SetColorIndex(x+int(xo), y+int(yo)+(yn*img.r), uint8(r))
		}
	}
}

// 增加噪点
func (img *Img) circle() {
	for i := 0; i < 20; i++ {
		colorIdx := uint8(rand.Intn(19)) + 1
		r := rand.Intn(img.r) + 1
		img.p.SetColorIndex(rand.Intn(img.Width-2*r)+r, rand.Intn(img.Height-2*r)+r, colorIdx)
	}
}

func NewImg(code int) *Img {
	codes := strconv.FormatInt(int64(code), 10)

	img := Img{
		Width:  240,
		Height: 80,
		Code:   codes,
	}

	palette := []color.Color{color.White, color.Black}
	rect := image.Rect(0, 0, img.Width, img.Height)

	img.p = image.NewPaletted(rect, palette)
	img.r = img.Height/fontHeight - 1

	// 计算开始绘制的xy轴。
	weight := fontWidth*img.r + 5
	x := (img.Width - weight*len(img.Code)) / 2
	y := (img.Height - fontHeight*img.r) / 2
	// 此步骤写入数字
	for _, i := range img.Code {
		img.drawFont(string(i), x, y)
		x += weight
	}
	img.distort()
	img.through()
	img.circle()
	return &img
}

func (img *Img) EncodePng(w io.Writer) {
	png.Encode(w, img.p)
}

func (img *Img) EncodeBase64() string {
	buf := new(bytes.Buffer)
	png.Encode(buf, img.p)
	data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + data
}
