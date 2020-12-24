package main

import (
	"strings"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		In  string
		Out []string
	}{
		{
			In: `2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms`,
			Out: []string{
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
		},
	}

	for i, tt := range tests {
		result, err := readFile(strings.NewReader(tt.In))
		if err != nil {
			t.Fatalf("failed %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("[%d] failed expected  %d group of lines got  %d ", i, len(tt.Out), 1)
		}
		if len(result["test"]) != len(tt.Out) {
			t.Fatalf("[%d] failed expected  %d lines got  %d ", i, len(tt.Out), len(result))
		}
		for j, out := range tt.Out {
			if out != result["test"][j].Raw {
				t.Fatalf("[%d] failed expected line %d to be:  \n%s\n but got \n%s\n", i, j, out, result["test"][j].Raw)
			}
		}
	}
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		In  string
		Out Line
	}{
		{
			In: "2020/12/06 18:00:01 te-st-1607277601 ended succesfully after 3.29207ms",
			Out: Line{
				Date:    time.Date(2020, 12, 6, 18, 0, 1, 0, time.UTC),
				Program: "te-st",
				Success: true,
				// Duration: 3.2907,
				Raw: "2020/12/06 18:00:01 te-st-1607277601 ended succesfully after 3.29207ms",
			},
		},
		{
			In: "2020/12/08 19:22:34 prod-send-new-1607282401 ended succesfully after 2m32.861754073s",
			Out: Line{
				Date:    time.Date(2020, 12, 8, 19, 22, 34, 0, time.UTC),
				Program: "prod-send-new",
				Success: true,
				// Duration: ???,
				Raw: "2020/12/08 19:22:34 prod-send-new-1607282401 ended succesfully after 2m32.861754073s",
			},
		},
		{
			In: "2020/12/07 11:10:07 prod-send-uncompleted-1607339401 ended with error",
			Out: Line{
				Date:    time.Date(2020, 12, 7, 11, 10, 07, 0, time.UTC),
				Program: "prod-send-uncompleted",
				Success: false,
				// Duration: ???,
				Raw: "2020/12/07 11:10:07 prod-send-uncompleted-1607339401 ended with error",
			},
		},
	}

	for _, tt := range tests {
		result := readLine(tt.In)
		if result != tt.Out {
			t.Fatalf("failed expected %#v but got %#v", tt.Out, result)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		Min      int
		allAfter time.Time
		In       []string
		Out      []string
	}{
		{
			Min:      3,
			allAfter: time.Date(2020, 12, 13, 11, 10, 07, 0, time.UTC),
			In: []string{
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
			Out: []string{
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
		},
		{
			Min:      3,
			allAfter: time.Date(2020, 12, 05, 11, 10, 07, 0, time.UTC),
			In: []string{
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
			Out: []string{
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
		},
		{
			Min:      20,
			allAfter: time.Date(2020, 12, 05, 11, 10, 07, 0, time.UTC),
			In: []string{
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
			Out: []string{
				"2020/12/12 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/11 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/10 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/09 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
				"2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms",
			},
		},
	}

	for i, tt := range tests {
		lines := []Line{}
		for _, l := range tt.In {
			lines = append(lines, readLine(l))
		}
		f := filter(tt.Min, tt.allAfter)
		result := f(lines)
		if len(result) != len(tt.Out) {
			t.Fatalf("[%d] failed expected  %d lines got  %d ", i, len(tt.Out), len(result))
		}
		for j, out := range tt.Out {
			if out != result[j].Raw {
				t.Fatalf("[%d] failed expected line %d to be:  \n%s\n but got \n%s\n", i, j, out, result[j].Raw)
			}
		}
	}
}
