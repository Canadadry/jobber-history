package svg

import (
	"github.com/canadadry/jobber-history/parser"
	"io"
	"text/template"
)

const (
	tmpl = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
{{ range . }}{{ range .Rects }}	<rect x="{{.Start}}" y="{{.Y}}" width="{{.Len}}" height="{{.Height}}" style="fill:{{.Color}}" />
{{ end }}	<text x="0" y="{{.TextY}}" font-family="Verdana" font-size="10">
		{{ .Name }}
	</text>
{{ end }}</svg>`
)

type SvgLine struct {
	Rects []Rect
	Name  string
	TextY int
}

type Rect struct {
	Start  float64
	Len    float64
	Color  string
	Height int
	Y      int
}

type Svg struct {
	Green       string
	Red         string
	LineHeight  int
	MarginTop   int
	TextYOffset int
}

func (s Svg) Convert(lines map[string][]parser.Line, out io.Writer) error {
	tmpl, err := template.New("svg").Parse(tmpl)
	if err != nil {
		return err
	}
	data := []SvgLine{}
	y := s.MarginTop
	for program, ln := range lines {
		l := SvgLine{
			Rects: s.convert(ln, y),
			Name:  program,
			TextY: y + s.TextYOffset,
		}
		data = append(data, l)
		y = y + s.LineHeight
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}

func (s Svg) convert(lines []parser.Line, y int) []Rect {
	greenRed := ternary(s.Green, s.Red)

	timestamps := []int64{}
	for _, l := range lines {
		timestamps = append(timestamps, l.Date.Unix())
	}
	fromMin := float64(timestamps[len(timestamps)-1])
	fromMax := float64(timestamps[0])
	toMin := 300.0
	toMax := 150.0
	toSvgCoord := remap(fromMin, fromMax, toMin, toMax)

	floatCoord := []Rect{}
	for i := 0; i < len(timestamps)-1; i++ {
		start := toSvgCoord(float64(timestamps[i+0]))
		end := toSvgCoord(float64(timestamps[i+1]))
		floatCoord = append(floatCoord, Rect{
			Start:  start,
			Len:    end - start,
			Color:  greenRed(lines[i].Success),
			Height: s.LineHeight,
			Y:      y,
		})
	}
	return floatCoord
}

func ternary(ifTrue, ifFalse string) func(bool) string {
	return func(s bool) string {
		if s {
			return ifTrue
		}
		return ifFalse
	}
}

func remap(fromMin, fromMax, toMin, toMax float64) func(float64) float64 {
	return func(in float64) float64 {
		absolute := (in - fromMin) / (fromMax - fromMin)
		return absolute*(toMax-toMin) + toMin
	}
}
