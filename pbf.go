package tiles

import (
	"fmt"
	"io"
	"os"

	"github.com/qedus/osmpbf"
)

type OsmData struct {
	Nodes map[int64]Node
	Ways  []Way
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
}

func (w Way) GetNodes(nodes map[int64]Node) []Node {
	newNodes := []Node{}
	for _, id := range w.NodeIDs {
		newNodes = append(newNodes, nodes[id])
	}
	return newNodes
}

func (w Way) GetNodePairs(nodes map[int64]Node) [][]Node {
	pairs := [][]Node{}
	nodeList := w.GetNodes(nodes)
	for i := 0; i < len(nodeList)-1; i++ {
		pairs = append(pairs, []Node{nodeList[i], nodeList[i+1]})
	}
	return pairs
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
		Nodes: map[int64]Node{},
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
				if _, ok := v.Tags["highway"]; ok {
					data.Ways = append(data.Ways, Way{v.NodeIDs})
				}
			case *osmpbf.Relation:
				// Ignore
			default:
				return nil, fmt.Errorf("unknown type %T", v)
			}
		}
	}

	return data, nil
}
