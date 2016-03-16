package tiles

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/qedus/osmpbf"
)

type ItemType int

const (
	ItemTypeNode = iota
	ItemTypeWay
	ItemTypeRelation
)

type OsmData struct {
	Nodes  map[int64]Node
	Ways   map[int64]Way
	Findex S2Index
}

type FeatureRef struct {
	Id    int64
	Type  ItemType
	FName string
}

type Node struct {
	Lon_, Lat_ float64
	Id         int64
}

func (p Node) Lon() float64 { return p.Lon_ }
func (p Node) Lat() float64 { return p.Lat_ }

func NodeFromPbf(n *osmpbf.Node) Node {
	return Node{
		Lon_: n.Lon,
		Lat_: n.Lat,
		Id:   n.ID,
	}
}

type Way struct {
	NodeIDs []int64
	Id      int64
	Tags    map[string]string
}

func WayFromPbf(w *osmpbf.Way) Way {
	return Way{
		NodeIDs: w.NodeIDs,
		Id:      w.ID,
		Tags:    w.Tags,
	}
}

func (w Way) Match(feature Feature) bool {
	for key, val := range w.Tags {
		for _, tag := range feature.Tags {
			if key == tag.Key && (val == tag.Val || tag.Val == "*") {
				return true
			}
		}
	}
	return false
}

func (w Way) MatchAny(features map[string]Feature) (string, bool) {
	for name, feature := range features {
		if w.Match(feature) {
			return name, true
		}
	}
	return "", false
}

func (w Way) GetNodes(nodes map[int64]Node) []Node {
	newNodes := []Node{}
	for _, id := range w.NodeIDs {
		newNodes = append(newNodes, nodes[id])
	}
	return newNodes
}

func ParsePbf(path string) (*OsmData, error) {
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

	data := &OsmData{
		Nodes:  map[int64]Node{},
		Ways:   map[int64]Way{},
		Findex: make(S2Index),
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				node := NodeFromPbf(v)
				data.Nodes[node.Id] = node
			case *osmpbf.Way:
				way := WayFromPbf(v)
				fName, ok := way.MatchAny(mapFeatures)
				if ok {
					data.Ways[way.Id] = way
					data.Findex.AddWay(way, fName, data)
				}
			case *osmpbf.Relation:
				// Ignore
			default:
				return nil, fmt.Errorf("unknown type %T", v)
			}
		}
	}

	log.Println("Num s2Cells", len(data.Findex))
	log.Println("Num ways", len(data.Ways))
	log.Println("Num nodes", len(data.Nodes))
	return data, nil
}
