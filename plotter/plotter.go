package plotter

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/majdanrc/eventplotter/events"
)

type Plotter struct {
	start time.Time
}

func (p *Plotter) Plot(events []interface{}, examinedDate time.Time, desc string) {
	var img = image.NewRGBA(image.Rect(0, 0, 800, 3000))
	addLabel(img, 500, 20, fmt.Sprintf("id: %s", desc))

	p.choosePainter(img, events)
	p.plotExaminedDate(examinedDate, img)

	f, err := os.Create("plot.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func (p *Plotter) choosePainter(img *image.RGBA, items []interface{}) {
	for _, ev := range items {
		switch ev := ev.(type) {
		case nil:
			return
		case events.VerticalEvent:
			kv := ev.KeyEvents()
			keyvs := fmt.Sprintf(">>> %s", kv)
			fmt.Println(keyvs)
			p.plotVertical(img, kv)
		case events.ProgressingEvent:
			kv := ev.KeyEvents()
			keyvs := fmt.Sprintf(">>> %s", kv)
			fmt.Println(keyvs)
			info := ev.Description()
			p.plotProgressing(img, kv, info)
		case events.BasicEvent:
			kv := ev.KeyEvents()
			keyvs := fmt.Sprintf(">>> %s", kv)
			fmt.Println(keyvs)
			info := ev.Description()
			p.plotBasic(img, kv, info)
		default:
			return
		}
	}
}

func (p *Plotter) plotVertical(img *image.RGBA, items []time.Time) {
	colors := []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}

	fmt.Println(p.start)

	if p.start.IsZero() {
		p.start = items[0]
	}

	begin := daysRange(p.start, items[0])
	duration := daysRange(items[0], items[1])

	rect(10, begin, 80, begin+duration, img, colors[0])
}

func (p *Plotter) plotProgressing(img *image.RGBA, items []time.Time, info []string) {
	colors := []color.Color{
		color.RGBA{255, 128, 255, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{240, 240, 20, 255}}

	//fmt.Println(p.start)

	if p.start.IsZero() {
		p.start = items[0]
	}

	creation := daysRange(p.start, items[0])
	begin := daysRange(p.start, items[1])
	duration := daysRange(p.start, items[2])

	hLine(100, creation, 120, img, colors[1])
	addLabel(img, 100, creation+15, info[0])
	vLine(120, creation-10, creation+10, img, colors[2])
	vLine(120, creation, begin, img, colors[0])
	hLine(120, begin, 180, img, colors[1])
	vLine(180, begin-10, begin+10, img, colors[2])
	vLine(180, begin, duration, img, colors[0])
	hLine(180, duration, 220, img, colors[1])
}

func (p *Plotter) plotBasic(img *image.RGBA, items []time.Time, info []string) {
	colors := []color.Color{color.RGBA{128, 128, 255, 255}, color.RGBA{0, 255, 0, 255}}

	//fmt.Println(p.start)

	if p.start.IsZero() {
		p.start = items[0]
	}

	begin := daysRange(p.start, items[0])

	addLabel(img, 300, begin+15, info[0])
	hLine(300, begin, 400, img, colors[1])
}

func (p *Plotter) plotExaminedDate(examinedDate time.Time, img *image.RGBA) {
	begin := daysRange(p.start, examinedDate)
	hLine(0, begin, 700, img, color.RGBA{100, 100, 200, 255})
}

func daysRange(start, end time.Time) int {
	firstRounded := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	lastRounded := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

	hours := lastRounded.Sub(firstRounded).Hours()
	return int(hours) / 24
}

func hLine(x1, y, x2 int, img *image.RGBA, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

func vLine(x, y1, y2 int, img *image.RGBA, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

func rect(x1, y1, x2, y2 int, img *image.RGBA, col color.Color) {
	hLine(x1, y1, x2, img, col)
	hLine(x1, y2, x2, img, col)
	vLine(x1, y1, y2, img, col)
	vLine(x2, y1, y2, img, col)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
