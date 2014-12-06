# go-tile-server

An in-memory OSM tile server written in Go

## Example

(The default map position is in Georgia - if you download a different pbf file expect to zoom out before you find your map)

    cd example
    wget http://download.geofabrik.de/north-america/us/georgia-latest.osm.pbf
    go run main.go georgia-latest.osm.pbf

## Development

You will need these things installed:

    % apt-get install protobuf-compiler
    % go get -u github.com/golang/protobuf/protoc-gen-go

To rebuild the protobuf .go files:

    % cd osmpbf && go generate
