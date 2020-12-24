package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error : ", err)
		os.Exit(1)
	}
}

type Line struct {
	Date     time.Time
	Program  string
	Success  bool
	Duration time.Duration
	Raw      string
}

func ByDecrisingDate(l1, pl2 *Line) bool {
	return l1.Date.After(pl2.Date)
}

type lineSorter struct {
	Lines []Line
	By    func(l1, l2 *Line) bool
}

func (s *lineSorter) Len() int {
	return len(s.Lines)
}

func (s *lineSorter) Swap(i, j int) {
	s.Lines[i], s.Lines[j] = s.Lines[j], s.Lines[i]
}

func (s *lineSorter) Less(i, j int) bool {
	return s.By(&s.Lines[i], &s.Lines[j])
}

func readLine(line string) Line {
	part := strings.Split(line, " ")
	if len(part) <= 4 {
		return Line{Raw: line}
	}
	t, _ := time.Parse("2006/01/02 15:04:05", part[0]+" "+part[1])
	p := strings.Split(part[2], "-")

	return Line{
		Date:    t,
		Program: strings.Join(p[:len(p)-1], "-"),
		Success: strings.Contains(strings.Join(part[3:], " "), "succesfully"),
		Raw:     line,
	}
}

func readFile(f io.Reader) (map[string][]Line, error) {
	lines := map[string][]Line{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := readLine(scanner.Text())
		lines[line.Program] = append(lines[line.Program], line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func filter(min int, allAfter time.Time) func([]Line) []Line {
	return func(in []Line) []Line {
		fmt.Println("filter")
		sort.Sort(&lineSorter{
			Lines: in,
			By:    ByDecrisingDate,
		})
		last := 0
		for i, l := range in {
			if l.Date.After(allAfter) {
				last = i + 1
			}
		}
		fmt.Println(last)
		if last < min {
			last = min
		}
		fmt.Println(last)
		if last > len(in) {
			last = len(in)
		}
		fmt.Println(last)
		return in[:last]
	}
}

func run() error {
	var (
		in  string
		out string
	)

	flag.StringVar(&in, "in", "", "input history.log file")
	flag.StringVar(&out, "out", "", "output svg file")

	flag.Parse()

	if len(in) == 0 || len(out) == 0 {
		return fmt.Errorf("param in and out must be set")
	}

	f, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("Can open %s : %w", in, err)
	}
	defer f.Close()

	lines, err := readFile(f)
	if err != nil {
		return fmt.Errorf("Can process %s : %w", in, err)
	}

	filtered := map[string][]Line{}

	for program, l := range lines {
		filtered[program] = filter(10, time.Now().Truncate(time.Hour*24))(l)
	}

	fmt.Println(len(lines))

	return nil
}
