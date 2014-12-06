package tiles

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const tileSize = 256

func drawTile(lonWest, latNorth, lonEast, latSouth, scale float64, font *truetype.Font, data *OsmData) (image.Image, error) {
	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Top/left border
	for i := 0; i < tileSize; i++ {
		img.Set(0, i, image.Black)
		img.Set(i, 0, image.Black)
	}

	// Dots
	for i := 0; i < tileSize; i += 16 {
		for j := 0; j < tileSize; j += 16 {
			img.Set(i, j, image.Black)
		}
	}

	// Tile location
	err := drawText(img, font, tileSize/2, 20, fmt.Sprintf("%f, %f", lonWest, latNorth))
	if err != nil {
		return nil, err
	}
	err = drawText(img, font, tileSize/2, tileSize-20, fmt.Sprintf("%f, %f", lonEast, latSouth))
	if err != nil {
		return nil, err
	}

	// Plot some nodes
	for _, node := range data.Nodes {
		if node.Lon > lonWest && node.Lon < lonEast &&
			node.Lat < latNorth && node.Lat > latSouth {
			x, y := getRelativeXY(lonWest, latNorth, node.Lon, node.Lat, scale)
			// fmt.Println(x, y)
			img.Set(round(x), round(y), image.Black)
		}
	}

	return img, nil
}

func drawText(img *image.RGBA, font *truetype.Font, x, y int, s string) error {
	var ptSize float64 = 12

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(ptSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.Black)
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
