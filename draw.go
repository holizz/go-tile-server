package tiles

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	img_font "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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

				coords := [][]float64{}
				for _, node := range way.GetNodes(data.Nodes) {
					x, y := getRelativeXY(nwPt, node, float64(zoom))
					coords = append(coords, []float64{x, y})
				}

				drawPolyLine(img, color.Black, coords)

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
	ctx.SetHinting(img_font.HintingFull)

	width := int(widthOfString(font, ptSize, s))
	pt := freetype.Pt(x-width/2, y+int(int32(ctx.PointToFixed(ptSize))>>8)/2)
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
			width += int(font.Kern(fixed.Int26_6(font.FUnitsPerEm()), prev, index))
		}
		width += int(font.HMetric(fixed.Int26_6(font.FUnitsPerEm()), index).AdvanceWidth)
		prev, hasPrev = index, true
	}

	return float64(width) * scale
}

func round(n float64) int {
	//TODO: this is incorrect
	return int(math.Floor(n))
}

func drawPolyLine(img *image.RGBA, color color.Color, coords [][]float64) {
	path := new(draw2d.Path)
	for i, coord := range coords {
		if i == 0 {
			path.MoveTo(coord[0], coord[1])
		} else {
			path.LineTo(coord[0], coord[1])
		}
	}

	gc := draw2dimg.NewGraphicContext(img)
	gc.SetStrokeColor(color)
	gc.Stroke(path)
}
