package tiles

import (
	"fmt"
	"io"
	"os"

	"github.com/qedus/osmpbf"
)

type OsmData struct {
	Nodes []Node
}

type Node struct {
	Lon, Lat float64
}

func parsePbf(path string) (*OsmData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(10)
	if err != nil {
		return nil, err
	}

	data := &OsmData{}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				data.Nodes = append(data.Nodes, Node{
					Lon: v.Lon,
					Lat: v.Lat,
				})
			case *osmpbf.Way:
				// Ignore
			case *osmpbf.Relation:
				// Ignore
			default:
				return nil, fmt.Errorf("unknown type %T", v)
			}
		}
	}

	return data, nil
}
