package tiles

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const tileSize = 256

func DrawTile(nwPt, sePt Pointer, zoom int64, font *truetype.Font, data *OsmData, debug bool) (image.Image, error) {
	t := time.Now()

	// Create white image
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Plot some features
	for fName, features := range data.Features {
		if mapFeatures[fName].MinZoom > zoom {
			continue
		}

		for _, feature := range features {
			switch feature.Type {
			case ItemTypeNode:
				//TODO
			case ItemTypeWay:
				way := data.Ways[feature.Id]
				//TODO: this ignores ways crossing over the tile
				withinBounds := false
				for _, node := range way.GetNodes(data.Nodes) {
					if node.Lon() > nwPt.Lon() && node.Lon() < sePt.Lon() &&
						node.Lat() < nwPt.Lat() && node.Lat() > sePt.Lat() {
						withinBounds = true
					}
				}

				if !withinBounds {
					continue
				}

				for _, pair := range way.GetNodePairs(data.Nodes) {
					x1, y1 := getRelativeXY(nwPt, pair[0], float64(zoom))
					x2, y2 := getRelativeXY(nwPt, pair[1], float64(zoom))

					drawLine(img, color.Black, x1, y1, x2, y2)
				}
			case ItemTypeRelation:
				//TODO
			}
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

		// Tile location
		err = drawText(img, font, red, tileSize/2, tileSize-20, time.Since(t).String())
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

func drawLine(img *image.RGBA, color color.Color, x1, y1, x2, y2 float64) {
	path := draw2d.NewPathStorage().MoveTo(x1, y1).LineTo(x2, y2)
	gc := draw2d.NewGraphicContext(img)
	gc.SetStrokeColor(color)
	gc.Stroke(path)

	slope := (y1 - y2) / (x1 - x2)
	yInt := slope*x1 - y1

	for x := 0; x < tileSize; x++ {
		if (float64(x) < x1 || float64(x) < x2) &&
			(float64(x) > x1 || float64(x) > x2) {
			img.Set(x, round(float64(x)*slope+yInt), color)
		}
	}
}
