# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-10-14

### BREAKING CHANGES
- `Sanitize` and `DefaultSanitizer()`: whitespace-only sanitizer strings (`"\n"`, `"\t"`, `"  "`) are replaced by a single space instead of removed, and runs of spaces are collapsed, so text no longer gets glued
- `Import` takes an `io.Reader` and `Export` takes an `io.Writer`. Use the new `ImportFile(filename)` and `ExportFile(filename)` for the previous file based behaviour
- `GetObjects` returns copies of the objects instead of the internal slice

### Added
- `OperationTypeEmpty` (`""`) for resources listed without operations
- `ResourceTypeEmpty` (`""`) for resources without data type
- Package `embedded` with `embedded.New()`: the sanitized registry compiled into the binary (about 3 MB when imported)
- `(*Reg).FindResourceByOIRAndVer`: OIR lookup in a given object version
- `(*Reg).FindObjectsByName` and `(*Reg).FindResourcesByName`: exact, case-insensitive name lookup, ordered by object ID and version
- `NewContext` and `(*Reg).ImportFromAPIContext` take a `context.Context`
- `Configuration.HTTPClient` for a custom HTTP client, and `Configuration.BaseURL` for the root URL of the objects list
- Exported sentinel errors for `errors.Is`
- Package overview in `doc.go` (sources, lookups, errors, concurrency)

### Changed
- `FindResourceByOIR` accepts a leading slash (`/3303/0/5700`)
- Godoc on every exported symbol 
- Errors returned by `Import`, `Export`, `ImportFile`, `ExportFile`, `New` and the API fetch carry context
- `Resource` declares its YAML keys explicitly
- `example/main.go` runs offline on the embedded registry
- `New` returns `*Reg` instead of the `Registry` interface
- `ImportFromAPI` returns an error when no objects were fetched
- OMA API requests use `https`, a shared HTTP client and `http.NewRequestWithContext`
- Updated registry.yaml and registry_sanitized.yaml
- Updated `yaml` and `testify` dependencies
- Upgraded go module from v1.16 to v1.27

### Fixed
- `FindResourceByOIR` doc was misleading about the used object version
- `FindObjectsByID` and `FindResourcesByObjResIDs` returned results in random order; they are now ordered by object version ascending, numerically (`1.2` before `1.10`)
- `New` ignored `Configuration.Sanitizer` and always used `DefaultSanitizer()`
- API responses with a non-2xx status were parsed as JSON/XML
- `Sanitize` did not modify object `Description1` and `Description2`
- `Export` produced YAML that `Import` could not parse
- `FindResourceByOIR` accepted IDs that overflow `int32` and negative IDs
- `Import` left the registry partly overwritten when the file failed to parse

## [0.2.0] - 2022-04-13

### Changed
- Package renamed to `ipso`

## [0.1.0] - 2022-04-13

### Added
- Find object by ID and version
- Find resource by object ID, object version and resource ID
- Find resource by OIR string ("3303/0/5700")
- Registry interface definition
- YAML registries exported from OMA API

### Changed
- Fixed golangci-linter complaints (v1.38.0)
- Moved registry to top-level directory

## [0.0.3] - 2021-03-07

### Added
- Empty Object's `ObjectVersion` and `LwM2MVersion` fields initialization with default value "1.0"

### Changed
- Replaced references in structures with objects

### Fixed
- Incorrect registry initialization from YAML file
- Slow object searching by ID

## [0.0.2] - 2021-01-15

### Added
- Added this changelog

### Changed
- Replaced OpenAPI client with native "net/http" calls 

## [0.0.1] - 2020-06-11
Initial release
