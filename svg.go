package main

import (
	"fmt"
	"io"
	"text/template"
	"time"
)

const (
	tmpl = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
{{ range $key, $value := . }}{{ range $value }}	<rect x="{{.Start}}" y="10" width="{{.Len}}" height="10" style="fill:{{.Color}}" />
{{ end }}	<text x="0" y="18" font-family="Verdana" font-size="10">
		{{ $key }}
	</text>
{{ end }}</svg>`
)

type Svg struct {
	Green string
	Red   string
}

func (s Svg) Convert(lines map[string][]Line, out io.Writer) error {
	tmpl, err := template.New("svg").Parse(tmpl)
	if err != nil {
		return err
	}
	data := map[string][]Rect{}
	for program, ln := range lines {
		data[program] = s.convert(ln)
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}

type Rect struct {
	Start float64
	Len   float64
	Color string
}

func (s Svg) convert(lines []Line) []Rect {
	greenRed := ternary(s.Green, s.Red)

	timestamps := []int64{}
	for _, l := range lines {
		fmt.Printf("%s : %d\n", l.Date.Format(time.RFC3339), l.Date.Unix())
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
			Start: start,
			Len:   end - start,
			Color: greenRed(lines[i].Success),
		})
	}
	fmt.Printf("convert %#v\n", floatCoord)
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
	fmt.Printf("from %g - %g \nto %g - %g\n", fromMin, fromMax, toMin, toMax)
	return func(in float64) float64 {
		absolute := (in - fromMin) / (fromMax - fromMin)
		return absolute*(toMax-toMin) + toMin
	}
}
