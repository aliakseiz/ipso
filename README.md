# IPSO registry

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

---

## Functionality

- Import registry from OMA API
- Export registry to YAML-file
- Import registry from file
- Compare two registries
- Find objects and resources by ID
- Find resources by object ID and resource ID
- Sanitize objects and resources text fields
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
    Sanitize: false,
    Sanitizer: nil,
}

reg, err := registry.New(cfg)
```

Compare two registries:
```go
comp := reg1.Compare(reg2)
```
Remove unwanted strings from objects and resources description:
```go
reg.Config.Sanitizer = registry.DefaultSanitizer()
reg.Sanitize()
```

---

# License
[MIT](LICENSE)

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg

[Stability-Status-Image]: http://badges.github.io/stability-badges/dist/experimental.svg

[Godoc-Url]: https://pkg.go.dev/mod/github.com/aliakseiz/ipso-registry
[Godoc-Image]: https://godoc.org/github.com/aliakseiz/ipso-registry?status.svg

[ReportCard-Url]: https://goreportcard.com/report/github.com/aliakseiz/ipso-registry
[ReportCard-Image]: https://goreportcard.com/badge/github.com/aliakseiz/ipso-registry