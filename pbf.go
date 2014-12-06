package tiles

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/holizz/go-tile-server/osmpbf"
)

type Pbf struct {
	Ways []Way
}

type Way struct {
	Nodes []Node
}

type Node struct {
	Lon float64
	Lat float64
}

func NewPbf(pbfPath string) *Pbf {
	pbf := &Pbf{}

	f, err := os.Open(pbfPath)
	if err != nil {
		panic(err)
	}

	size := readInt4(f)

	blobHeader := &OSMPBF.BlobHeader{}

	pbfData := make([]byte, size)
	_, err = f.Read(pbfData)
	if err != nil {
		panic(err)
	}

	err = proto.Unmarshal(pbfData, blobHeader)
	if err != nil {
		panic(err)
	}

	fmt.Println(blobHeader)
	fmt.Println(*blobHeader.Type)
	fmt.Println(*blobHeader.Datasize)

	os.Exit(1)

	return pbf
}

func readInt4(r io.Reader) uint32 {
	b := make([]byte, 4)
	_, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	var val uint32 = 0
	for _, i := range b {
		val = (val << 8) + uint32(i)
	}

	return val
}
