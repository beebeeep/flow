package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Field [][]float64

type obj struct {
	x, y, m float64
}

var (
	_gravityC = 1.0
	resX      = 2000.0
	resY      = 2000.0
	scale     = 0.008 * resX
	step      = 0.01 * resX
	nR        = int(resY / scale)
	nC        = int(resX / scale)
	pSize     = scale / 10.0
	field     Field
)

func fillRect(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64, c color.RGBA) {
	gc.SetFillColor(c)
	gc.MoveTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.LineTo(x2, y2)
	gc.LineTo(x2, y1)
	gc.Close()
	gc.FillStroke()
}

func drawLine(gc *draw2dimg.GraphicContext, x0, y0 float64, iterations int, colors []color.RGBA) {
	d := math.Sqrt(math.Pow(resX/2-x0, 2.0)+math.Pow(resY/2-y0, 2.0)) / (resX)
	i := int(rand.Float64() * 2.5 * d * float64(len(colors)))
	if i >= len(colors) {
		i = len(colors) - 1
	}
	c := colors[i]

	gc.SetStrokeColor(c)
	x := x0
	y := y0
	gc.SetLineWidth(5)
	gc.BeginPath()
	gc.MoveTo(x, y)
	for i := 0; i < iterations; i++ {
		r := int(x / scale)
		c := int(y / scale)
		if r < 0 {
			r = 0
		}
		if c < 0 {
			c = 0
		}
		if r >= nR {
			r = nR - 1
		}
		if c >= nC {
			c = nC - 1
		}
		v := field[r][c]
		x += step * math.Cos(v)
		y += step * math.Sin(v)
		gc.LineTo(x, y)
	}
	gc.Stroke()
}

func (f Field) render(gc *draw2dimg.GraphicContext) {
	fillRect(gc, 0, 0, resX, resY, color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetStrokeColor(color.RGBA{0x30, 0xa0, 0, 0xff})

	for r := range field {
		for c := range field[r] {
			v := field[r][c]
			gc.BeginPath()
			x := float64(r) * scale
			y := float64(c) * scale
			fillRect(gc, x-pSize, y-pSize, x+pSize, y+pSize, color.RGBA{0x30, 0xa0, 0, 0xff})
			gc.MoveTo(x, y)
			gc.LineTo(x+scale*math.Cos(v), y+scale*math.Sin(v))
			gc.LineTo(x+scale*math.Cos(v), y+scale*math.Sin(v))
			gc.Close()
			gc.FillStroke()
		}
	}

}

func (o obj) render(gc *draw2dimg.GraphicContext) {
	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0, 0xff})
	gc.SetFillColor(color.RGBA{0xff, 0, 0, 0xff})
	//gc.MoveTo(o.x, o.y)
	gc.ArcTo(o.x, o.y, 10*o.m, 10*o.m, 0, 2.0*math.Pi)
	gc.FillStroke()
}

func gravityField(objects []obj, field Field) {
	for r := range field {
		x := float64(r) * scale
		for c := range field[r] {
			y := float64(c) * scale
			var gx, gy float64
			for _, obj := range objects {
				g := _gravityC * obj.m / math.Sqrt(math.Pow((obj.x-x), 2.0)+math.Pow((obj.y-y), 2.0))
				a := math.Atan2(obj.y-y, obj.x-x)
				gx += g * math.Cos(a)
				gy += g * math.Sin(a)
			}
			field[r][c] = math.Atan2(gy, gx)
		}
	}
}

func draw(img *image.RGBA, numLines, numIters int) {
	field = make([][]float64, nR)
	for i := range field {
		field[i] = make([]float64, nC)
	}

	/*
		for r := range field {
			for c := range field[r] {
				//field[r][c] = 2 * math.Pi * (1.0 - (math.Cos(1.0 * float64(c) / float64(r+1))))
				field[r][c] = math.Pi * float64(r) / float64(nR)
			}
		}
	*/
	var objects = make([]obj, 0)
	for i := 0; i < 50; i++ {
		objects = append(objects, obj{
			x: rand.Float64() * resX,
			y: rand.Float64() * resY,
			m: rand.Float64() * 3,
		})
	}
	gravityField(objects, field)

	gc := draw2dimg.NewGraphicContext(img)
	fillRect(gc, 0, 0, resX, resY, color.RGBA{0x00, 0x00, 0x00, 0xff})
	/*
		field.render(gc)
		for _, o := range objects {
			o.render(gc)
		}
	*/
	colors := []color.RGBA{
		{0xE8, 0xE0, 0x89, 0xff},
		{0x69, 0xD4, 0xF0, 0xff},
		{0xF0, 0xE3, 0x51, 0xff},
		{0xF0, 0x3A, 0x79, 0xff},
		{0xA3, 0x34, 0x5B, 0xff},
	}

	sort.Slice(colors, func(i, j int) bool {
		return rand.Intn(100) < 50
	})
	for i := 0; i < numLines; i++ {
		drawLine(gc, resX*rand.Float64(), resY*rand.Float64(), numIters, colors)
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	img := image.NewRGBA(image.Rect(0, 0, int(resX), int(resY)))
	numLines := flag.Int("lines", 10000, "numer of lines")
	numIters := flag.Int("iters", 5, "number of iterations")
	outFile := flag.String("out", "out.png", "output file")
	flag.Parse()
	draw(img, *numLines, *numIters)
	f, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	/*
		scene := img.SubImage(image.Rect(
			int(resX*0.25), int(resY*0.25),
			int(resX*0.75), int(resY*0.75),
		))
		if err := png.Encode(f, scene); err != nil {
			log.Fatal(err.Error())
		}
	*/
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("written %s\n", *outFile)

}
