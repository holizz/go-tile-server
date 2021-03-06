package tiles

import (
	"io/ioutil"
	"testing"

	"code.google.com/p/freetype-go/freetype"
)

func BenchmarkDrawTile(b *testing.B) {
	nwPt := Point{-4.482421875, 54.162433968067795}
	sePt := Point{-4.471435546875, 54.156001090284924}
	scale := int64(15)

	// Read font
	font_, err := ioutil.ReadFile("example/FiraSans-Regular.ttf")
	if err != nil {
		b.Fatalf("Benchmark setup failed: %#v\n", err)
	}

	font, err := freetype.ParseFont(font_)
	if err != nil {
		b.Fatalf("Benchmark setup failed: %#v\n", err)
	}

	// Read PBF file
	data, err := ParsePbf("example/isle-of-man-latest.osm.pbf")
	if err != nil {
		b.Fatalf("Benchmark setup failed: %#v\n", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := DrawTile(nwPt, sePt, scale, font, data, false)
		if err != nil {
			b.Fatalf("Received error: %#v\n", err)
		}
	}
}
