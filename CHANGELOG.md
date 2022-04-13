# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Find object by ID and version
- Find resource by object ID, object version and resource ID
- Find resource by OIR string ("3303/0/5700")
- Registry interface definition
- YAML registries exported from OMA API

### Changed
- Fixed golangci-linter complaints (v1.38.0)
- Moved registry to top-level directory

### Fixed

### Deprecated

### Removed

### Security

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
