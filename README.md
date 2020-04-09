## LwM2M registry

---

### Usage examples

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

Install openapi-code-generator:

```npm install @openapitools/openapi-generator-cli -g```

Generate client:

```npx openapi-generator generate -i api/openapi.yaml -g go -o api/client```

Client with XML tags:
```npx openapi-generator generate -i openapi.yaml -p withXml=true -g go -o client```