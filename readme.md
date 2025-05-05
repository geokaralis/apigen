# apigen

An opinionated cli utility written in Go for generating type-safe api clients from openapi specifications.

## Installation
``` bash
git clone https://github.com/geokaralis/apigen.git
cd apigen
go build -o apigen cmd/apigen/main.go
```

## Usage
``` bash
apigen create client schema.json -o ./__generated__
```

Options:
- `<schema>` path to openapi schema file
- `--output, -o` directory where the generated files will be stored
