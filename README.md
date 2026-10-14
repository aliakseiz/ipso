# IPSO registry

[![Go Reference][Godoc-Image]][Godoc-Url]
[![License][License-Image]][License-Url]


Go library for the OMA LwM2M IPSO object and resource registry: load it from the OMA API, a YAML file or the
copy compiled into your binary, then search, compare and sanitize it.

```shell
go get github.com/aliakseiz/ipso
```

Requires Go 1.27 or newer.

---

## Functionality

- Import registry from OMA API
- Export registry to YAML, to a file or any `io.Writer`
- Import registry from YAML, from a file or any `io.Reader` (`embed.FS`, HTTP body, buffer)
- Compare two registries
- Find objects and resources by ID and version
- Find resources by object ID and resource ID
- Find objects and resources by name, ignoring letter case
- Find resources by OIR string i.e. "3303/0/5700" or "/3303/0/5700"
- Sanitize objects and resources text fields
- Sanitized registry compiled into the binary (`ipso/embedded`)

---

## Data model

- `ipso.Object`: an OMA object such as 3303 Temperature, identified by `ObjectID` and `ObjectVersion`. It holds
  its `Resources.Item` list.
- `ipso.Resource`: a value of an object such as 5700 Sensor Value, with `Operations`, `Type`, `Units` and so on.
  Its ID is unique within one object version.
- Enumerations (`InstanceType`, `MandatoryType`, `OperationType`, `ResourceType`, `DifferenceType`) are strings with
  constants for the values in the OMA registry.
- The YAML written by `Export` uses the field names as keys for objects (`ObjectID`, `Description1`) and lowercase
  keys for resources (`multipleinstances`, `rangeenumeration`). Keep that in mind when editing a file by hand.

Registry files:

| File                            | Content                                                               |
|---------------------------------|-----------------------------------------------------------------------|
| `share/registry.yaml`           | Registry as fetched from the OMA API, raw descriptions                |
| `share/registry_sanitized.yaml` | `share/registry.yaml` sanitized with `DefaultSanitizer()`             |
| `embedded/registry.yaml`        | Copy of `share/registry_sanitized.yaml` compiled into `ipso/embedded` |

---

## Usage examples

Initialize a registry from OMA API:

```go
reg, err := ipso.New(ipso.DefaultConfiguration())
```

The fetch requests the objects list and then every object definition one after another, to keep the load on the
public OMA and GitHub hosts low. It makes about 500 requests and takes minutes (2m19s for 513 objects when measured).
For regular use, export the registry once and load it with `ImportFile`, or load a file from `share/`.

A runnable demo is in `example/`: `go run ./example` works offline on the embedded registry,
`go run ./example -fetch -out /tmp/ipso` downloads from the OMA API and writes YAML files.

Errors carry context and wrap the exported sentinels, so test them with `errors.Is(err, ipso.ErrObjectNotFound)`.

Export initialized registry to YAML file:

```go
err := reg.ExportFile("registry.yaml")
```

Import a previously exported registry from YAML file:

```go
err := reg.ImportFile("registry.yaml")
```

Read and write any stream, for example a registry compiled into the binary:

```go
//go:embed registry.yaml
var registryYAML []byte

err := reg.Import(bytes.NewReader(registryYAML))
err = reg.Export(os.Stdout)
```

Create a registry with custom configuration. `Sanitizer` lists the strings cleaned from descriptions by sanitizing (removed; whitespace-only ones such
as `"\t"` become a space); `DefaultSanitizer()` is used when it is empty:

```go
cfg := ipso.Configuration{
InitOnNew:      true,
SkipInitErrors: true,
Sanitize:       true,
Sanitizer:      []string{"\t", "•"},
}

reg, err := ipso.New(cfg)
```

Use the sanitized registry compiled into the binary, without files or network (adds about 3 MB to the executable):

```go
import "github.com/aliakseiz/ipso/embedded"

reg, err := embedded.New() // parses the data on every call, keep and reuse the registry
```

Cancel the OMA API fetch with a context, and use your own HTTP client (proxy, tracing, retries):

```go
cfg := ipso.DefaultConfiguration()
cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second} // used for the objects list and every object definition

ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

reg, err := ipso.NewContext(ctx, cfg) // or reg.ImportFromAPIContext(ctx) on an existing registry
```

Compare two registries:

```go
comp := reg1.Compare(reg2.GetRegistry())
```

Remove unwanted strings from objects and resources description:

```go
reg.Sanitize(ipso.DefaultSanitizer())
```

Find a resource by OIR string (object/instance/resource); a leading slash is accepted:

```go
res, err := reg.FindResourceByOIR("/3303/0/5700") // object version 1.0
res, err = reg.FindResourceByOIRAndVer("/3303/0/5700", "1.1") // another object version
```

Find objects and resources by name when the ID is unknown (the match is exact, letter case is ignored):

```go
objects, err := reg.FindObjectsByName("temperature") // every version of object 3303
resources, err := reg.FindResourcesByName("Sensor Value") // across all objects and versions
```

---

## Behavior notes

- OMA API requests use `https` and fail on a non-2xx status. The default HTTP client times out each request after
  5 seconds; a `Configuration.HTTPClient` replaces it. Response bodies are limited to 16 MiB.
- `Configuration.BaseURL` changes only where the objects list is requested (`BaseURL + "/Object"`). Object definitions
  are fetched from the `ObjectLink` URLs in that list (the official API links to `raw.githubusercontent.com`).
- `ImportFromAPI` returns an error when no objects were fetched, for example when the API is down
  and `SkipInitErrors` is set. It returns the objects and does not modify the registry.
- `Import` and `ImportFile` replace the registry contents only when the whole YAML is parsed. On error the registry is unchanged.
- `Import` and `Export` work on streams. `ImportFile` and `ExportFile` take a file name and return
  `ErrEmptyFilename` for an empty one; `ExportFile` creates the file with mode 0600 or truncates an existing one.
- OIR segments must be non-negative numbers that fit `int32`.
- `FindObjectsByID` and `FindResourcesByObjResIDs` return results ordered by object version ascending,
  compared numerically (`1.2` before `1.10`). The last element is the latest version.
- `FindObjectsByName` and `FindResourcesByName` match the whole name, ignoring letter case; an empty name matches
  nothing. Results are ordered by object ID, then by object version ascending. Descriptions are free text and are not searchable.
- Objects returned by `Find`, `FindObjectByIDAndVer`, `FindObjectsByID`, `FindObjectsByName`, `GetObjects` and
  `Compare` are copies. Changing them does not change the registry.
- A registry is safe for concurrent reads. `Import`, `ImportFile`, `Sanitize` and `New` with `InitOnNew` modify it
  and must not run concurrently with other calls.
- `Compare` lists objects present in only one registry (matched by ID and version); it does not compare the content
  of objects present in both.
- `Sanitize` cleans `Description1` and `Description2` of objects as well as `Description`, `RangeEnumeration` and
  `Units` of resources.
  Without a whitespace-only string inner spaces are left as they are.
- `share/registry_sanitized.yaml` is `share/registry.yaml` sanitized with `DefaultSanitizer()`, and
  `embedded/registry.yaml` is a copy of it; tests fail when either is out of date.

---

# License

[MIT](LICENSE)

[License-Url]: LICENSE
[License-Image]: https://img.shields.io/github/license/aliakseiz/ipso

[Godoc-Url]: https://pkg.go.dev/github.com/aliakseiz/ipso
[Godoc-Image]: https://pkg.go.dev/badge/github.com/aliakseiz/ipso.svg
