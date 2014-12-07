# go-tile-server

In-memory OSM tile server written in Go

So far it's very basic - it just renders roads and coastlines, and it takes forever to parse big PBF files.

## Example

    go get -u github.com/holizz/go-tile-server
    cd $GOPATH/src/github.com/holizz/go-tile-server/example
    go run main.go

And visit http://localhost:3000/

## License

MIT

Things in example/ directory may have different licenses. widthOfString function may also be under a different license.
