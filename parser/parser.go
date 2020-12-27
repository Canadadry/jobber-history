package parser

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"time"
)

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

func ReadFile(f io.Reader) (map[string][]Line, error) {
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

func Filter(min int, allAfter time.Time) func([]Line) []Line {
	return func(in []Line) []Line {
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
		if last < min {
			last = min
		}
		if last > len(in) {
			last = len(in)
		}
		return in[:last]
	}
}
