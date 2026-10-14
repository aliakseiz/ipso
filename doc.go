// Package ipso holds the OMA LwM2M IPSO registry: objects and their resources, as published by the
// Open Mobile Alliance.
//
// A registry ([Reg]) is filled from one of three sources:
//   - the OMA API, with [New] and [Configuration.InitOnNew] or [Reg.ImportFromAPI];
//   - a YAML file or stream written by [Reg.Export], with [Reg.ImportFile] or [Reg.Import];
//   - the sanitized registry compiled into the binary, with the package github.com/aliakseiz/ipso/embedded.
//
// Fetching from the OMA API makes about 500 requests one after another and takes minutes, so export the
// registry once and load the file afterwards.
//
// # Lookups
//
// Objects are found by ID and version ([Reg.FindObjectByIDAndVer], [Reg.FindObjectsByID]) or by name
// ([Reg.FindObjectsByName]). Resources are found by ID ([Reg.FindResourcesByID]), by object and resource ID
// ([Reg.FindResourcesByObjResIDs], [Reg.FindResourceByObjIDObjVerResID]), by an OIR string such as
// "3303/0/5700" ([Reg.FindResourceByOIR], [Reg.FindResourceByOIRAndVer]) or by name ([Reg.FindResourcesByName]).
// Objects and resources returned by lookups are copies: changing them does not change the registry.
//
// [Reg.Compare] lists the objects that exist in only one of two registries, and [Reg.Sanitize] cleans the free
// text of descriptions.
//
// # Errors
//
// Errors carry context and wrap the sentinel errors ([ErrObjectNotFound],
// [ErrResourceNotFound], [ErrInvalidOIR], [ErrEmptyFilename], [ErrNoObjects], [ErrHTTPStatus],
// [ErrBodyTooLarge]). Match them with [errors.Is], not by message text.
//
// # Concurrency
//
// A registry is safe for concurrent reads. [Reg.Import], [Reg.ImportFile], [Reg.Sanitize] and [New] with
// [Configuration.InitOnNew] modify it and must not run concurrently with other calls on the same registry.
package ipso
