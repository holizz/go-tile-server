package main

import (
	"io/ioutil"

	"code.google.com/p/freetype-go/freetype"
	"github.com/davecheney/profile"
	"github.com/holizz/go-tile-server"
)

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	nwPt := tiles.Point{-4.482421875, 54.162433968067795}
	sePt := tiles.Point{-4.471435546875, 54.156001090284924}
	scale := float64(15)

	// Read font
	font_, err := ioutil.ReadFile("../example/FiraSans-Regular.ttf")
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(font_)
	if err != nil {
		panic(err)
	}

	// Read PBF file
	data, err := tiles.ParsePbf("../example/isle-of-man-latest.osm.pbf")
	if err != nil {
		panic(err)
	}

	_, err = tiles.DrawTile(nwPt, sePt, scale, font, data, false)
	if err != nil {
		panic(err)
	}
}
