module github.com/ORKL/library/olm/olm

go 1.19

replace github.com/ORKL/library/olm/utils => ../utils

require (
	github.com/ORKL/library/olm/utils v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/barasher/go-exiftool v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
)
