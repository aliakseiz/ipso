# IPSO registry

[![License][License-Image]][License-Url]
[![Godoc][Godoc-Image]][Godoc-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

What is IPSO: https://omaspecworks.org/develop-with-oma-specworks/ipso-smart-objects/

---

## Functionality

- Import registry from [OMA API](https://technical.openmobilealliance.org/OMNA/LwM2M/LwM2MRegistry.html)
- Export registry to YAML-file
- Import registry from file
- Compare two registries
- Find objects and resources by ID and version
- Find resources by object ID and resource ID
- Find resources by OIR string i.e. "3303/0/5700"
- Sanitize objects and resources text fields
---

## Usage examples

Initialize a registry from [OMA API](https://technical.openmobilealliance.org/OMNA/LwM2M/LwM2MRegistry.html):
```go
reg, err := registry.New(ipso.DefaultConfiguration())
```

Export initialized registry to YAML file:
```go
err := reg.Export("registry.yaml")
```

Import a previously exported registry from YAML file:
```go
err := reg.Import("registry.yaml")
```

Create a registry with custom configuration:
```go
cfg := ipso.Configuration{
    InitOnNew:      false,
    SkipInitErrors: false,
    Sanitize: false,
}

reg, err := ipso.New(cfg)
```

Compare two registries:
```go
comp := reg1.Compare(reg2.GetRegistry())
```
Remove unwanted strings from objects and resources description:
```go
reg.Sanitize(ipso.DefaultSanitizer())
```

---
## Linter

```shell
go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0
golangci-lint run --enable-all
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
