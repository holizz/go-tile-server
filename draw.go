package tiles

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const tileSize = 256

func drawTile(nwPt, sePt Pointer, scale float64, font *truetype.Font, data *OsmData, debug bool) (image.Image, error) {
	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Plot some nodes
	for _, node := range data.Nodes {
		if node.Lon() > nwPt.Lon() && node.Lon() < sePt.Lon() &&
			node.Lat() < nwPt.Lat() && node.Lat() > sePt.Lat() {
			x, y := getRelativeXY(nwPt, node, scale)
			img.Set(round(x), round(y), image.Black)
		}
	}

	// Debugging
	if debug {
		red := color.RGBA{0xff, 0x00, 0x00, 0xff}

		// Top/left border
		for i := 0; i < tileSize; i++ {
			img.Set(0, i, red)
			img.Set(i, 0, red)
		}

		// Tile location
		err := drawText(img, font, red, tileSize/2, 20, fmt.Sprintf("%f, %f", nwPt.Lon(), nwPt.Lat()))
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

func drawText(img *image.RGBA, font *truetype.Font, color color.Color, x, y int, s string) error {
	var ptSize float64 = 12

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(ptSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.NewUniform(color))
	ctx.SetHinting(freetype.FullHinting)

	width := int(widthOfString(font, ptSize, s))
	pt := freetype.Pt(x-width/2, y+int(ctx.PointToFix32(ptSize)>>8)/2)
	_, err := ctx.DrawString(s, pt)
	if err != nil {
		return err
	}

	return nil
}

// https://code.google.com/p/plotinum/source/browse/vg/font.go#160
func widthOfString(font *truetype.Font, size float64, s string) float64 {
	// scale converts truetype.FUnit to float64
	scale := size / float64(font.FUnitsPerEm())

	width := 0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := font.Index(rune)
		if hasPrev {
			width += int(font.Kerning(font.FUnitsPerEm(), prev, index))
		}
		width += int(font.HMetric(font.FUnitsPerEm(), index).AdvanceWidth)
		prev, hasPrev = index, true
	}

	return float64(width) * scale
}

func round(n float64) int {
	//TODO: this is incorrect
	return int(math.Floor(n))
}
