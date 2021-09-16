module github.com/churrodata/spacex-geo-extension

go 1.15

require (
	github.com/churrodata/churro v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.3.0
	github.com/ohler55/ojg v1.12.4
	github.com/rs/zerolog v1.23.0
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.27.1
)

replace github.com/churrodata/churro => ../churro
