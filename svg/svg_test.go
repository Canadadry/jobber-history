package svg

import (
	"app/parser"
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestToSvg(t *testing.T) {

	Green := "rgb(0,255,0)"
	Red := "rgb(255,0,0)"
	LineHeight := 10
	MarginTop := 10
	TextYOffset := 8

	tests := []struct {
		in  string
		out string
	}{
		{
			in: ``,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
</svg>`,
		},
		{
			in: `2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms`,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
	<rect x="150" y="10" width="150" height="10" style="fill:rgb(0,255,0)" />
	<text x="0" y="18" font-family="Verdana" font-size="10">
		test
	</text>
</svg>`,
		},
		{
			in: `2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/07 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms`,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
	<rect x="150" y="10" width="75" height="10" style="fill:rgb(0,255,0)" />
	<rect x="225" y="10" width="75" height="10" style="fill:rgb(0,255,0)" />
	<text x="0" y="18" font-family="Verdana" font-size="10">
		test
	</text>
</svg>`,
		},
		{
			in: `2020/12/06 18:00:01 test-1607277601 ended failed after 3.29207ms
2020/12/07 18:00:01 test-1607277601 ended failed after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms`,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
	<rect x="150" y="10" width="75" height="10" style="fill:rgb(0,255,0)" />
	<rect x="225" y="10" width="75" height="10" style="fill:rgb(255,0,0)" />
	<text x="0" y="18" font-family="Verdana" font-size="10">
		test
	</text>
</svg>`,
		},
		{
			in: `2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/07 18:00:01 test-1607277601 ended failed after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms`,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
	<rect x="150" y="10" width="75" height="10" style="fill:rgb(0,255,0)" />
	<rect x="225" y="10" width="75" height="10" style="fill:rgb(255,0,0)" />
	<text x="0" y="18" font-family="Verdana" font-size="10">
		test
	</text>
</svg>`,
		},
		{
			in: `2020/12/06 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/08 18:00:01 test-1607277601 ended succesfully after 3.29207ms
2020/12/06 18:00:01 test2-1607277601 ended succesfully after 3.29207ms
2020/12/08 18:00:01 test2-1607277601 ended succesfully after 3.29207ms`,
			out: `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 300 150" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
	<rect x="150" y="10" width="150" height="10" style="fill:rgb(0,255,0)" />
	<text x="0" y="18" font-family="Verdana" font-size="10">
		test
	</text>
	<rect x="150" y="20" width="150" height="10" style="fill:rgb(0,255,0)" />
	<text x="0" y="28" font-family="Verdana" font-size="10">
		test2
	</text>
</svg>`,
		},
	}

	for i, tt := range tests {

		lines, err := parser.ReadFile(strings.NewReader(tt.in))
		if err != nil {
			t.Fatalf("[%d] failed %v", i, err)
		}
		filtered := map[string][]parser.Line{}
		for program, l := range lines {
			filtered[program] = parser.Filter(100, time.Now())(l)
		}

		buf := bytes.Buffer{}
		err = Svg{Green, Red, LineHeight, MarginTop, TextYOffset}.Convert(filtered, &buf)
		if err != nil {
			t.Fatalf("[%d] failed %v", i, err)
		}

		if buf.String() != tt.out {
			t.Fatalf("[%d] failed got \n%s\n expected \n%s\n", i, buf.String(), tt.out)
		}
	}
}

func TestRemap(t *testing.T) {
	tests := []struct {
		fromMin  float64
		fromMax  float64
		toMin    float64
		toMax    float64
		in       []float64
		expected []float64
	}{
		{
			fromMin:  0,
			fromMax:  100,
			toMin:    150,
			toMax:    300,
			in:       []float64{0, 50, 100},
			expected: []float64{150, 225, 300},
		},
		{
			fromMin:  12678,
			fromMax:  1567922,
			toMin:    0,
			toMax:    100,
			in:       []float64{12678, 401489, 790300, 1567922},
			expected: []float64{0, 25, 50, 100},
		},
		{
			fromMin:  1607277601,
			fromMax:  1607450401,
			toMin:    150,
			toMax:    300,
			in:       []float64{1607277601, 1607450401},
			expected: []float64{150, 300},
		},
		{
			fromMin:  1607277601,
			fromMax:  1607450401,
			toMin:    300,
			toMax:    150,
			in:       []float64{1607277601, 1607450401},
			expected: []float64{300, 150},
		},
	}

	for i, tt := range tests {
		conv := remap(tt.fromMin, tt.fromMax, tt.toMin, tt.toMax)

		for j, v := range tt.in {
			result := conv(v)
			if result != tt.expected[j] {
				t.Fatalf("[%d:%d] failed got %g expected %g", i, j, result, tt.expected[j])
			}
		}
	}
}
