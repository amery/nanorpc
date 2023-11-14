module github.com/amery/nanorpc

go 1.19

require (
	github.com/amery/nanorpc/pkg/reconnect v0.0.0-00010101000000-000000000000
	github.com/amery/protogen v0.3.10
)

require (
	darvaza.org/core v0.10.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/amery/nanorpc/pkg/reconnect => ./pkg/reconnect
