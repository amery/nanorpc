module github.com/amery/nanorpc

go 1.21.9

replace (
	darvaza.org/core => ../../../darvaza.org/core
	darvaza.org/sidecar => ../../../darvaza.org/sidecar
)

require (
	darvaza.org/core v0.13.1 // indirect
	darvaza.org/sidecar v0.4.0
	darvaza.org/slog v0.5.7
	darvaza.org/slog/handlers/discard v0.4.11 // indirect
	darvaza.org/slog/handlers/filter v0.4.9 // indirect
	darvaza.org/slog/handlers/zerolog v0.4.9 // indirect
	darvaza.org/x/config v0.3.2 // indirect
	github.com/amery/defaults v0.1.0 // indirect
	github.com/amery/nanorpc/pkg/nanorpc v0.0.0-00010101000000-000000000000
)

require (
	github.com/amery/nanorpc/pkg/nanopb v0.0.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace (
	github.com/amery/nanorpc/pkg/generator => ./pkg/generator
	github.com/amery/nanorpc/pkg/nanopb => ./pkg/nanopb
	github.com/amery/nanorpc/pkg/nanorpc => ./pkg/nanorpc
)
