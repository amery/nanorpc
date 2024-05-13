module github.com/amery/nanorpc/pkg/nanorpc

go 1.21.9

replace (
	darvaza.org/core => ../../../../../darvaza.org/core
	darvaza.org/sidecar => ../../../../../darvaza.org/sidecar
)

require (
	darvaza.org/core v0.13.1
	darvaza.org/sidecar v0.4.0
	darvaza.org/slog v0.5.7
	darvaza.org/slog/handlers/discard v0.4.11
	darvaza.org/x/config v0.3.2
	github.com/amery/defaults v0.1.0 // indirect
	github.com/amery/nanorpc/pkg/nanopb v0.0.0
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

replace github.com/amery/nanorpc/pkg/nanopb => ../nanopb
