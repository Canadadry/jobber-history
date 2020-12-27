package main

import (
	"app/parser"
	"app/svg"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error : ", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		in   string
		out  string
		min  int
		hour int
	)

	flag.StringVar(&in, "in", "", "input history.log file")
	flag.StringVar(&out, "out", "", "output svg file")
	flag.IntVar(&min, "min", 10, "minimal number of line to concider")
	flag.IntVar(&hour, "hour", 0, "take all line from the last X hour this date")

	flag.Parse()

	if len(in) == 0 || len(out) == 0 {
		return fmt.Errorf("param in and out must be set")
	}

	f, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("Can open %s : %w", in, err)
	}
	defer f.Close()

	lines, err := parser.ReadFile(f)
	if err != nil {
		return fmt.Errorf("Can process %s : %w", in, err)
	}

	filtered := map[string][]parser.Line{}

	for program, l := range lines {
		filtered[program] = parser.Filter(min, time.Now().Truncate(time.Hour*time.Duration(hour)))(l)
	}
	fOut, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Can open out file %s : %w", out, err)
	}
	defer f.Close()

	Green := "rgb(0,255,0)"
	Red := "rgb(255,0,0)"
	LineHeight := 10
	MarginTop := 10
	TextYOffset := 8

	return svg.Svg{Green, Red, LineHeight, MarginTop, TextYOffset}.Convert(filtered, fOut)

}
