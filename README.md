# IPSO registry

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

---

## Usage examples

Initialize registry from OMA API:
```go
reg, err := registry.New(registry.DefaultConfiguration())
```

Export initialized registry to YAML file:
```go
err := regAPI.Export("registry.yaml")
```

Import previously exported registry from YAML file:
```go
err := reg.Import("registry.yaml")
```

Create registry with custom configuration:
```go
cfg := &registry.Configuration{
    InitOnNew:      false,
    SkipInitErrors: false,
}

reg, err := registry.New(cfg)
```

Compare two registries:
```go
comp := reg1.Compare(reg2)
```

---

## OpenAPI code generator usage

Install openapi-code-generator:

```npm install @openapitools/openapi-generator-cli -g```

Generate client:

```npx openapi-generator generate -i api/openapi.yaml -g go -o api/client```

Client with XML tags:
```npx openapi-generator generate -i openapi.yaml -p withXml=true -g go -o client```

# License
[MIT](LICENSE)

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg

[Stability-Status-Image]: http://badges.github.io/stability-badges/dist/experimental.svg

[Godoc-Url]: https://pkg.go.dev/mod/github.com/aliakseiz/ipso-registry
[Godoc-Image]: https://godoc.org/github.com/aliakseiz/ipso-registry?status.svg

[ReportCard-Url]: https://goreportcard.com/report/github.com/aliakseiz/ipso-registry
[ReportCard-Image]: https://goreportcard.com/badge/github.com/aliakseiz/ipso-registry