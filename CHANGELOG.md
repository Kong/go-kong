# Table of Contents

- [v0.52.0](#v0520)
- [v0.51.0](#v0510)
- [v0.50.0](#v0500)
- [v0.49.0](#v0490)
- [v0.48.0](#v0480)
- [v0.47.0](#v0470)
- [v0.46.0](#v0460)
- [v0.45.0](#v0450)
- [v0.44.0](#v0440)
- [v0.43.0](#v0430)
- [v0.42.0](#v0420)
- [v0.41.0](#v0410)
- [v0.40.0](#v0400)
- [v0.39.2](#v0392)
- [v0.39.1](#v0391)
- [v0.39.0](#v0390)
- [v0.38.1](#v0381)
- [v0.38.0](#v0380)
- [v0.37.0](#v0370)
- [v0.36.0](#v0360)
- [v0.35.0](#v0350)
- [v0.34.1](#v0341)
- [v0.34.0](#v0340)
- [v0.33.0](#v0330)
- [v0.32.0](#v0320)
- [v0.31.1](#v0311)
- [v0.31.0](#v0310)
- [v0.30.0](#v0300)
- [v0.29.0](#v0290)
- [v0.28.1](#v0281)
- [v0.28.0](#v0280)
- [v0.27.0](#v0270)
- [v0.26.0](#v0260)
- [v0.25.1](#v0251)
- [v0.25.0](#v0250)
- [v0.24.0](#v0240)
- [v0.23.0](#v0230)
- [v0.22.0](#v0220)
- [v0.21.0](#v0210)
- [v0.20.0](#v0200)
- [v0.19.0](#v0190)
- [v0.18.0](#v0180)
- [v0.17.0](#v0170)
- [v0.16.0](#v0160)
- [v0.15.0](#v0150)
- [v0.14.0](#v0140)
- [v0.13.0](#v0130)
- [v0.12.0](#v0120)
- [v0.11.0](#v0110)
- [v0.10.0](#v0100)
- [v0.9.0](#v090)
- [v0.8.0](#v080)
- [v0.7.0](#v070)
- [v0.6.2](#v062)
- [v0.6.1](#v061)
- [v0.6.0](#v060)
- [v0.5.1](#v051)
- [v0.5.0](#v050)
- [v0.4.1](#v041)
- [v0.4.0](#v040)
- [0.3.0](#030)
- [0.2.0](#020)
- [0.1.0](#010)

## [v0.52.0]

> Release date: 2024/03/26

- Allow configuration of `User-Agent` for the client.
  [#427](https://github.com/Kong/go-kong/pull/427)

## [v0.51.0]

> Release date: 2024/01/23

- Added extract of `details` field in response when CRUD API fails to extract
  details from Konnect APIs.
  [#399](https://github.com/Kong/go-kong/pull/399)

## [v0.50.0]

> Release date: 2024/01/09

- Added `workspace` parameter to `FillID` method of entities to avoid entities
  with same key in different workspaces to have the same ID, which is unique
  cross all workspaces.
  [#395](https://github.com/Kong/go-kong/pull/395)

## [v0.49.0]

> Release date: 2024/01/04

- Added `Validate` method in vault service to validate vaults.
  [#392](https://github.com/Kong/go-kong/pull/392)
- Added `FillID()` method for vault.
  [#391](https://github.com/Kong/go-kong/pull/391)

## [v0.48.0]

> Release date: 2023/10/30

- `Route`'s `priority` field type is changed from `int` to `uint64`.
  [#378](https://github.com/Kong/go-kong/pull/378)
- Fix a bug preventing users to set fields to empty arrays when
  a default for those fields exist.
  [#374](https://github.com/Kong/go-kong/pull/374)

## [v0.47.0]

> Release date: 2023/08/29

- Added method `Validate` to `RouteService`
  [#368](https://github.com/Kong/go-kong/pull/368)

## [v0.46.0]

> Release date: 2023/07/17

- Added method `FillID()` method for Consumer Group.
  [#357](https://github.com/Kong/go-kong/pull/357)
- Added support for scoping plugins to Consumer Groups.
  [#352](https://github.com/Kong/go-kong/pull/352)

## [v0.45.0]

> Release date: 2023/07/03

- Added `Client.Config` handler, that can be used to obtain Kong's config.
  [#354](https://github.com/Kong/go-kong/pull/354)

## [v0.44.0]

> Release date: 2023/06/22

- Fix: handle empty array as nil when filling record defaults
  [#345](https://github.com/Kong/go-kong/pull/345)
- Fix leaking response body from `Client.Do()`. From now on `kong.Response` does
  not contain the `http.Response` field which was in fact leaking the implementation
  detail. It now contains the response headers, status and status code.
  If users want to get the response body they can still provide the `v` parameter
  to `Do()` to get the body unmarshalled into it (or copied into it if an `io.Writer`
  was to be provided).
  [#348](https://github.com/Kong/go-kong/pull/348)

## [v0.43.0]

> Release date: 2023/06/06

- Added missing `use_srv_name` and `healthchecks.active.headers` to `Upstream`
  entity.
  [#331](https://github.com/Kong/go-kong/pull/331)
- Fix regressions when filling defaults for set and record schema types.
  [#333](https://github.com/Kong/go-kong/pull/333)
  [#336](https://github.com/Kong/go-kong/pull/336)

## [v0.42.0]

> Release date: 2023/05/15

- Added `Details` and `SetDetails` methods to `APIError` struct to enable
  propagating status-specific details in there.
  Introduced `ErrTooManyRequestsDetails` struct that's going to be available
  as an `APIError`'s `Details` when Admin API returns status code 429 along
  with a `Retry-After` header. That should be useful for handling rate limiting
  on an application level.
  [#323](https://github.com/Kong/go-kong/pull/323)
- Use modern schema retrieval function for plugin schema retrieval. This
  provides compatibility with upcoming Kong releases.
  [#326](https://github.com/Kong/go-kong/pull/326)

## [v0.41.0]

> Release date: 2023/04/25

- Added `FillID` method for `Service`, `Route` and `Consumer` entities. It allows
  setting a deterministic ID to an entity.
  [#299](https://github.com/Kong/go-kong/pull/299)

## [v0.40.0]

> Release date: 2023/04/07

- Added support for the `/licenses` endpoint.
  [#304](https://github.com/Kong/go-kong/pull/304)

## [v0.39.2]

> Release date: 2023/04/05

This release reverts the fix from 0.39.1 while debugging [an
issue](https://github.com/Kong/go-kong/issues/307) not originally caught in
tests.

## [v0.39.1]

> Release date: 2023/04/04

- fix: default population now correctly handles arrays of records.
  [#302](https://github.com/Kong/go-kong/pull/302)

## [v0.39.0]

> Release date: 2023/03/17

- feat: add GraphQL-RateLimiting-Advanced and DeGraphQL DAOs (#288)
  [#288](https://github.com/Kong/go-kong/pull/288)

## [v0.38.1]

> Release date: 2023/02/22

- Omit response body from error when config push fails. The body can be quite
  large and will result in massive logs downstream if included in the error.
  The body is returned along with the error and downstream clients and log it
  separately if they so desire.

## [v0.38.0]

> Release date: 2023/02/13

- Added `instance_name` to Plugin schema.
  [#280](https://github.com/Kong/go-kong/pull/280)

## [v0.37.0]

> Release date: 2023/02/03

- **Breaking change:** the `ConfigService` is now directly embedded in the
  `kong.Client`. Configurations are collections of entities, not entities
  themselves, so they do not fit with other go-kong services.
- **Breaking change:** `ReloadDeclarativeRawConfig()` (formerly part of
  `ConfigService`, now part of `kong.Client`) now has the response signature
  `([]byte, error)` instead of `error`. The byte slice is the config response
  body. The error is unchanged.
- **Breaking change:** `ReloadDeclarativeRawConfig()` now requires a
  `flattenErrors` boolean argument. When `true`, requests will include
  `flatten_errors=1` in the query string, to activate the functionality added
  in <https://github.com/Kong/kong/pull/10161>.
  [#273](https://github.com/Kong/go-kong/pull/273)

## [v0.36.0]

> Release date: 2023/01/24

- Added `NewAPIErrorWithRaw()` to return a go-kong APIError along with the original raw error body.
  [#237](https://github.com/Kong/go-kong/pull/237)
  [#267](https://github.com/Kong/go-kong/pull/267)

## [v0.35.0]

> Release date: 2023/01/19

- Add consumer_groups tags support
  [#265](https://github.com/Kong/go-kong/pull/265)
- Add IsForbiddenErr error checking for 403s
  [#264](https://github.com/Kong/go-kong/pull/264)

## [v0.34.1]

> Release date: 2022/12/22

- Fix ingestion of entity defaults with arbitray map values
  [#258](https://github.com/Kong/go-kong/pull/258)

## [v0.34.0]

> Release date: 2022/12/19

- Add support to consumer_groups
  [#221](https://github.com/Kong/go-kong/pull/221)
- Add Keys and Key-sets as core entities
  [#238](https://github.com/Kong/go-kong/pull/238)
- Add `BaseRootURL()` to Client
  [#255](https://github.com/Kong/go-kong/pull/255)
- Add possibility to client to send declarative configs via `ReloadDeclarativeRawConfig()`
  [#252](https://github.com/Kong/go-kong/pull/252)
- Add OAuth2 client_type field
  [#235](https://github.com/Kong/go-kong/pull/235)
- Fix ingestion of entity defaults with map values
  [#244](https://github.com/Kong/go-kong/pull/244)
- Add support to filling entity defaults using JSON schemas.
  [#231](https://github.com/Kong/go-kong/pull/231)

## [v0.33.0]

> Release date: 2022/10/05

- Add support to Kong Vaults
  [#224](https://github.com/Kong/go-kong/pull/224)
- Add plugin to a route
  [#217](https://github.com/Kong/go-kong/pull/217)

## [v0.32.0]

> Release date: 2022/08/23

This release is not yet fully compatible with Kong 3.x. The target service `MarkHealthy()` and `MarkUnhealthy()`
functions do not yet work on 3.x. Other functionality supports 3.x, but may not yet support all new 3.x schema fields.

- Added support for 3.x Enterprise version strings.
  [#207](https://github.com/Kong/go-kong/pull/207)
- Added support for `expression` and `priority` route fields.
  [#210](https://github.com/Kong/go-kong/pull/210)
- Dropped support for 1.5 Enterprise, which exited sunset support in 2022-04.
- Added 3.x to test matrices.

## [v0.31.1]

> Release date: 2022/08/23

- Fixed support for plugins with IDs in `ForService` variants of plugin calls.
  [#205](https://github.com/Kong/go-kong/pull/205)

## [v0.31.0]

> Release date: 2022/08/16

- Updated to Go 1.19.
- Added `ForService` variants of plugin calls, which use
  `/service/<id>/plugins/` endpoints.
  [#192](https://github.com/Kong/go-kong/pull/192)
- Added support for Kong 3.0 upstream fields.
  [#200](https://github.com/Kong/go-kong/pull/200)

## [v0.30.0]

> Release date: 2022/07/11

- Updated to Go 1.18.
  [#164](https://github.com/Kong/go-kong/pull/164)
- Fixed an issue with the HTTP client wrapper not setting any transport when
  none was configured. It now uses the default transport.
  [#179](https://github.com/Kong/go-kong/pull/179)
- Added support for the new `ordering` field on Enterprise 3.x plugins.
  [#187](https://github.com/Kong/go-kong/pull/187)

## [v0.29.0]

> Release date: 2022/03/25

- Fixed a bug where `HTTPClientWithHeaders()` could not accept a `nil`
  (default) client.
  [#138](https://github.com/Kong/go-kong/pull/138)
- Fixed a bug where the plugin service's `Validate()` method could not handle
  errors associated with empty responses.
  [#142](https://github.com/Kong/go-kong/pull/142)
- Fixed a bug where plugin schema default retrieval included fields outside the
  plugin's `config` object.
  [#145](https://github.com/Kong/go-kong/pull/145)
- Fixed a typo in the `EntityCRUD.CRUDPath` JSON serialization.
  [#149](https://github.com/Kong/go-kong/pull/149)
- Added support for the new `configuration_hash` field in the status service.
  [#150](https://github.com/Kong/go-kong/pull/150)
- Fixed incorrect JSON marshalling/unmarshalling for RBAC endpoint permissions.
  [#148](https://github.com/Kong/go-kong/pull/148)
- Fixed a bug with RBAC endpoint permission path construction.
  [#148](https://github.com/Kong/go-kong/pull/148)

## [v0.28.1]

> Release date: 2022/02/22

- Export the `requiredFeatures` struct to be able to use `RunWhenEnterprise()`
  in other packages.
  [#135](https://github.com/Kong/go-kong/pull/135)

## [v0.28.0]

> Release date: 2022/02/14

- The `RunWhenKong()` and `RunWhenEnterprise()` functions are now exported for
  other packages that need to condition tests on Kong versions.
  [#132](https://github.com/Kong/go-kong/pull/132)
- RBAC permission actions are now an array of strings rather than a single
  comma-separated string to properly reflect their Kong type.
  [#131](https://github.com/Kong/go-kong/pull/132)
- Enterprise versions returned by `ParseSemanticVersion()` now include their
  supplemental Enterprise version as build info, rather than as a prerelease.
  This fixes an issue where semver comparisons saw the initial patch of a minor
  release as a lower version of that release.
  [#129](https://github.com/Kong/go-kong/pull/129)

## [v0.27.0]

> Release date: 2022/02/03

- Zero values (boolean `false` and integer `0`) now override default values
  when filling entity defaults.
  [#125](https://github.com/Kong/go-kong/pull/125)

## [v0.26.0]

> Release date: 2022/01/31

- Add missing entity fields from newer Kong releases.
  [#120](https://github.com/Kong/go-kong/pull/120)
- Add `SchemaService` and `FillEntityDefaults` utility supporting
  filling defaults for Services, Routes, Upstream and Targets.
  [#119](https://github.com/Kong/go-kong/pull/119)
- Add `GetFullSchema` to plugins to retrieve their complete schema.
  Also add `FillPluginsDefaults` utility to fill a plugin with its defaults.
  [#114](https://github.com/Kong/go-kong/pull/114)

## [v0.25.1]

> Release date: 2021/12/08

v0.25.1 reverts the k8s.io/code-generator version originally used in v0.25.0 to
avoid an unwanted upgrade of its github.com/go-logr/logr dependency.

## [v0.25.0]

> Release date: 2021/12/08

- Healthcheck types now include CRD validation annotations for
  [Kubebuilder](https://book.kubebuilder.io/reference/markers/crd-validation.html).
  [#104](https://github.com/Kong/go-kong/pull/104)

## [v0.24.0]

> Release date: 2021/11/05

### Added

- A `Listeners()` method was added to the Kong `Client` in order to retrieve
  the `proxy_listeners` and `stream_listeners` from the root conveniently.
  [#101](https://github.com/Kong/go-kong/pull/101)

## [v0.23.0]

> Release date: 2021/09/22

### Breaking changes

- The plugin service `Validate()` API's return signature is now `bool, string,
  err` rather than `bool, err`. The additional string value contains the
  validation failure reason if validation successfully determined that the
  proposed configuration was not valid. The err value is now only set if
  `Validate()` could not successfully retrieve validation information because
  of a timeout/DNS lookup failure/etc.

## [v0.22.0]

> Release date: 2021/09/22

### Added

- Test that target updates work as expected.
  [#85](https://github.com/Kong/go-kong/pull/85)
- Workspace service supports new `ExistsByName()` API. This is similar to
  `Exists()`, but it only accepts names, not IDs. This function can check if a
  workspace exists even if the RBAC user has access to that namespace only.
  `Exists()` requires access to the default namespace.
  [#90](https://github.com/Kong/go-kong/pull/90)

### Fixed

- The generic `exists()` client functions uses GETs instead of HEADs. Kong
  <=2.6 respond to HEAD requests with 200s as long as the endpoint has a valid
  form. They do not check if entities actually exist, and `exists()` would
  return true for entities that don't exist because of this.
  [#90](https://github.com/Kong/go-kong/pull/90)

## [v0.21.0]

> Release date: 2021/08/26

### Added

- oauth2 entities support `hash_secret`.
  [#74](https://github.com/Kong/go-kong/pull/74)

### Fixed

- Plugin validation checks against the correct status code.
  [#81](https://github.com/Kong/go-kong/pull/81)

## [v0.20.0]

> Release date: 2021/07/07

### Added

- FriendlyName is now defined for entities based on deck types.
  [#68](https://github.com/Kong/go-kong/pull/68)
- Added Info service for interacting with information exposed by the admin API root endpoint.
  [#65](https://github.com/Kong/go-kong/pull/65)
- Implemented client-level support for Kong workspaces.
  [#62](https://github.com/Kong/go-kong/pull/62)

### Changed

- Internally wrapped errors now use the standard Go library instead of a 3rd party
  wrapping lib and consequently several string versions of errors have changed.
  [#66](https://github.com/Kong/go-kong/pull/66)
- Various improvements to repository CI: caching go modules, updating codegen, e.t.c.

## [v0.19.0]

> Release date: 2021/05/05

### Added

- Client now allows reading the version of Kong. [#48](https://github.com/Kong/go-kong/pull/48)

## [v0.18.0]

> Release date: 2021/05/05

### Added

- Workspace and tag existence functions. [#56](https://github.com/Kong/go-kong/pull/56)
- Plugin schema retrieval. [#57](https://github.com/Kong/go-kong/pull/57)

### Changed

- The default client now sets a default timeout. [#51](https://github.com/Kong/go-kong/pull/51)
- Optimized HTTP header handling. [#49](https://github.com/Kong/go-kong/pull/49)
- Changed return type of `kong.HTTPClientWithHeaders(...)`: now returns `*http.Client` instead of `http.Client`. [#46](https://github.com/Kong/go-kong/pull/46)
- Now uses Go 1.16 module semantics. [#40](https://github.com/Kong/go-kong/pull/40)
- Now uses Dependabot to manage GitHub Actions. [#41](https://github.com/Kong/go-kong/pull/41), [#42](https://github.com/Kong/go-kong/pull/42)
- Now uses code-generator v0.21.0. [#43](https://github.com/Kong/go-kong/pull/43)
- CI no longer uses Bintray repositories. [#55](https://github.com/Kong/go-kong/pull/55)
- Dropped testing against older versions of Kong. [#58](https://github.com/Kong/go-kong/pull/58)

## [v0.17.0]

> Release date: 2021/04/05

### Added

- Added support for Developer entities. Thanks to @ChristianJacquot!
  [#27](https://github.com/Kong/go-kong/pull/27)
- Added support for Developer Role entities. Thanks to @mmorel-35!
  [#30](https://github.com/Kong/go-kong/pull/27)
- RBAC roles now support pagination and listing all entities. Thanks to
  @mmorel-35!
  [#30](https://github.com/Kong/go-kong/pull/27)
- Tests can now require the Portal. Added helpers to enable the Portal and
  related settings.
  [#30](https://github.com/Kong/go-kong/pull/27)
- Clients now use interfaces. Other libraries that use go-kong can define mock
  types that implement these interfaces for unit tests that do not require an
  actual Kong instance.
  [#24](https://github.com/Kong/go-kong/pull/27)

### Fixed

- RBAC roles now properly include their `negative` field in requests to Kong.
  [#32](https://github.com/Kong/go-kong/pull/27)

## [v0.16.0]

> Release date: 2021/03/03

### Added

- `Service` now includes `URL`.

## [v0.15.0]

> Release date: 2021/01/22

### Added

- `Route` now includes `RequestBuffering` and `ResponseBuffering`.

## [v0.14.0]

> Release date: 2021/01/12

### Breaking changes

- HTTP error format changed to feature HTTP codes

### Added

- RBACUser support
- Support for auto-expiring key-auth key TTL
- Support for RBAC roles
- Support for RBAC permissions
- DeepCopy annotations for enterprise entities

## [v0.13.0]

> Release date: 2020/08/04

### Summary

This release renames the package from `github.com/hbagdi/go-kong` to
`github.com/kong/go-kong`.

## [v0.12.0]

> Release date: 2020/07/30

### Summary

This release adds support for Kong 2.1 series and
a number of enterprise entities.

### Added

- Added `HTTPClientWithHeaders` helper function to inject authn/authz
  related HTTP headers in requests to kong.

- New fields in types:
  - `Service` struct has three new fields:
    - `TLSVerifyDepth`
    - `TLSVerify`
    - `CACertificates`
  - `ClientCertificate` field has been added to `Upstream` struct.
  - `Type` field has been added to `PassiveHealthcheck` struct.

- Following new services for Kong Enterprise have been introduced:
  - WorkspaceService to manage workspaces in kong
  - AdminService to manage users of Kong Manager
  - MTLSAuthService to manage MTLS credentials

### Misc

- Changed the branch name from `master` to `main`
- Introduced linters to improve code health

## [v0.11.0]

> Release date: 2020/01/17

### Summary

- This release adds support for Kong 2.0.0.

### Added

- `Threshold` field has been added to Upstream struct.
- `PathHandling` field has been added to Route struct.

## [v0.10.0]

> Release date: 2019/10/27

### Summary

- This release adds support for Kong 1.4.

### Added

- `HostHeader` field has been added to Upstream struct.
- `Tags` field has been added to the following types:
  - `KeyAuth`
  - `Basicauth`
  - `HMACAuth`
  - `Oauth2Credential`
  - `ACLGroup`
  - `JWTAuth`

## [v0.9.0]

> Release date: 2019/08/24

### Breaking changes

- `client.Do()` returns a response object even on errors so that
  clients can inspect the response directly when needed.
  The error condition has changed to include HTTP status codes 300 to 399
  as success and not failure.
  [b1f9eda31](https://github.com/hbagdi/go-kong/commit/b1f9eda311e1d9c9d6b0f5a5e33a3d399d85faf6)
- `String()` method has been dropped from all types defined in this package.

### Added

- `NewRequest()` method helping with creating HTTP requests is now exported
- URLs parsed inside the package are more robust.
- New method `GetByCustomID` has been introduced to fetch Consumers by
  `custom_id`.

## [v0.8.0]

> Release date: 2019/08/21

### Added

- Oauth2Credential type and service has been added
  which can be used to create Oauth2 credentials in Kong for some
  Oauth2 flows.

## [v0.7.0]

> Release date: 2019/08/13

### Summary

This release brings support for CRUD methods for
authentication credentials in Kong.

### Added

- The following credentials and corresponding services have been added:
  - `key-auth`
  - `basic-auth`
  - `hmac-auth`
  - `jwt`
  - `acl`

## [v0.6.2]

> Release date: 2019/08/09

### Fix

- Add missing omitempty tag to ClientCertificate field

## [v0.6.1]

> Release date: 2019/08/09

### Fix

- Fix a typo in Service struct definition for YAML tag

## [v0.6.0]

> Release date: 2019/08/09

### Summary

- This release adds support for Kong 1.3.

### Breaking change

- `Validator` Interface has been dropped and Valid() method from
  all entities is dropped.
  [#8](https://github.com/hbagdi/go-kong/pull/8)

### Added

- `Headers` field has been added to Route struct.
  [#5](https://github.com/hbagdi/go-kong/pull/5)
- `ClientCertificate` field has been added to Service struct.
  [#5](https://github.com/hbagdi/go-kong/pull/5)
- `CACertificate` is a new core entity in Kong. A struct to represent
  it and a corresponding new service is added.
  [#5](https://github.com/hbagdi/go-kong/pull/5)
- `Algorithm` field has been added to Upstream struct.
  [#9](https://github.com/hbagdi/go-kong/pull/9)

## [v0.5.1]

> Release date: 2019/08/05

### Fix

- Add missing healthchecks.active.unhealthy.interval field to Upstream
  [#6](https://github.com/hbagdi/go-kong/issues/6)

## [v0.5.0]

> Release date: 2019/06/07

### Summary

- This release adds support for Kong 1.2.

### Added

- Added HTTPSRedirectStatusCode property to Route struct.
  [#3](https://github.com/hbagdi/go-kong/pull/3)

### Breaking change

- `Create()` for Custom Entities now supports HTTP PUT method.
  If `id` is specified in the object, it will be used to PUT the entity.
  This was always POST previously.
  [#3](https://github.com/hbagdi/go-kong/pull/3)

## [v0.4.1]

> Release date: 2019/04/11

### Fix

- Add `omitempty` property to Upstream fields for Kong 1.0 compatibility

## [v0.4.0]

> Release date: 2019/04/06

### Summary

- This release adds support for features released in Kong 1.1.
  This version is compatible with Kong 1.0 and Kong 1.1.

### Breaking Change

- Please note that the version naming scheme for this library has changed from
  `x.y.z` to `vX.Y.Z`. This is to ensure compatibility with Go modules.

### Added

- `Tags` field has been added to all Kong Core entity structs.
- List methods now support tag based filtering introduced in Kong 1.1.
  Tags can be ANDed or ORed together. `ListOpt` struct can be used to
  specify the tags for filtering.
- `Protocols` field has been added to Plugin struct.
- New fields `Type`, `HTTPSSni` and `HTTPSVerifyCertificate` have been
  introduced for Active HTTPS healthchecks.
- `TargetService` has two new methods `MarkHealthy()` and `MarkUnhealthy()`
  to change the health of a target.

## [0.3.0]

> Release date: 2018/12/19

### Summary

- This release adds support for Kong 1.0.
  It is not compatible with 0.x.y  versions of Kong due to breaking
  Admin API changes as the deprecated API entity is dropped.
- The code and API for the library is same as 0.2.0, with the exception
  that struct defs and services related to `API` is dropped.

### Breaking changes

- `API` struct definition is no longer available.
- `APIService` is no longer available. Please ensure your code doesn't rely
  on these before upgrading.
- `Plugin` struct has dropped the `API` field.

## [0.2.0]

> Release date: 2018/12/19

### Summary

- This release adds support for Kong 0.15.x.
  It is not compatible with any other versions of Kong due to breaking
  Admin API changes in Kong for Plugins, Upstreams and Targets entities.

### Breaking changes

- `Target` struct now has an `Upstream` member in place of `UpstreamID`.
- `Plugin` struct now has `Consumer`, `API`, `Route`, `Service` members
  instead of `ConsumerID`, `APIID`, `RouteID` and `ServiceID`.

### Added

- `RunOn` property has been added to `Plugin`.
- New properties are added to `Route` for L4 proxy support.

## [0.1.0]

> Release date: 2018/12/01

### Summary

- Debut release of this library
- This release comes with support for Kong 0.14.x
- The library is not expected to work with previous or later
  releases of Kong since every release of Kong is introducing breaking changes
  to the Admin API.

[v0.51.0]: https://github.com/Kong/go-kong/compare/v0.50.0...v0.51.0
[v0.50.0]: https://github.com/Kong/go-kong/compare/v0.49.0...v0.50.0
[v0.49.0]: https://github.com/Kong/go-kong/compare/v0.48.0...v0.49.0
[v0.48.0]: https://github.com/Kong/go-kong/compare/v0.47.0...v0.48.0
[v0.47.0]: https://github.com/Kong/go-kong/compare/v0.46.0...v0.47.0
[v0.46.0]: https://github.com/Kong/go-kong/compare/v0.45.0...v0.46.0
[v0.45.0]: https://github.com/Kong/go-kong/compare/v0.44.0...v0.45.0
[v0.44.0]: https://github.com/Kong/go-kong/compare/v0.43.0...v0.44.0
[v0.43.0]: https://github.com/Kong/go-kong/compare/v0.42.0...v0.43.0
[v0.42.0]: https://github.com/Kong/go-kong/compare/v0.41.0...v0.42.0
[v0.41.0]: https://github.com/Kong/go-kong/compare/v0.40.0...v0.41.0
[v0.40.0]: https://github.com/Kong/go-kong/compare/v0.39.2...v0.40.0
[v0.39.2]: https://github.com/Kong/go-kong/compare/v0.39.1...v0.39.2
[v0.39.1]: https://github.com/Kong/go-kong/compare/v0.39.0...v0.39.1
[v0.39.0]: https://github.com/Kong/go-kong/compare/v0.38.1...v0.39.0
[v0.38.1]: https://github.com/Kong/go-kong/compare/v0.38.0...v0.38.1
[v0.38.0]: https://github.com/Kong/go-kong/compare/v0.37.0...v0.38.0
[v0.37.0]: https://github.com/Kong/go-kong/compare/v0.36.0...v0.37.0
[v0.36.0]: https://github.com/Kong/go-kong/compare/v0.35.0...v0.36.0
[v0.35.0]: https://github.com/Kong/go-kong/compare/v0.34.1...v0.35.0
[v0.34.1]: https://github.com/Kong/go-kong/compare/v0.34.0...v0.34.1
[v0.34.0]: https://github.com/Kong/go-kong/compare/v0.33.0...v0.34.0
[v0.33.0]: https://github.com/Kong/go-kong/compare/v0.32.0...v0.33.0
[v0.32.0]: https://github.com/Kong/go-kong/compare/v0.31.1...v0.32.0
[v0.31.1]: https://github.com/Kong/go-kong/compare/v0.31.0...v0.31.1
[v0.31.0]: https://github.com/Kong/go-kong/compare/v0.30.0...v0.31.0
[v0.30.0]: https://github.com/Kong/go-kong/compare/v0.29.0...v0.30.0
[v0.29.0]: https://github.com/Kong/go-kong/compare/v0.28.1...v0.29.0
[v0.28.1]: https://github.com/Kong/go-kong/compare/v0.28.0...v0.28.1
[v0.28.0]: https://github.com/Kong/go-kong/compare/v0.27.0...v0.28.0
[v0.27.0]: https://github.com/Kong/go-kong/compare/v0.26.0...v0.27.0
[v0.26.0]: https://github.com/Kong/go-kong/compare/v0.25.1...v0.26.0
[v0.25.1]: https://github.com/Kong/go-kong/compare/v0.24.0...v0.25.1
[v0.25.0]: https://github.com/Kong/go-kong/compare/v0.24.0...v0.25.0
[v0.24.0]: https://github.com/Kong/go-kong/compare/v0.23.0...v0.24.0
[v0.23.0]: https://github.com/Kong/go-kong/compare/v0.22.0...v0.23.0
[v0.22.0]: https://github.com/Kong/go-kong/compare/v0.21.0...v0.22.0
[v0.21.0]: https://github.com/Kong/go-kong/compare/v0.20.0...v0.21.0
[v0.20.0]: https://github.com/Kong/go-kong/compare/v0.19.0...v0.20.0
[v0.19.0]: https://github.com/Kong/go-kong/compare/v0.18.0...v0.19.0
[v0.18.0]: https://github.com/Kong/go-kong/compare/v0.17.0...v0.18.0
[v0.17.0]: https://github.com/Kong/go-kong/compare/v0.16.0...v0.17.0
[v0.16.0]: https://github.com/Kong/go-kong/compare/v0.15.0...v0.16.0
[v0.15.0]: https://github.com/Kong/go-kong/compare/v0.14.0...v0.15.0
[v0.14.0]: https://github.com/Kong/go-kong/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/Kong/go-kong/compare/v0.12.0...v0.13.0
[v0.12.0]: https://github.com/Kong/go-kong/compare/v0.11.0...v0.12.0
[v0.11.0]: https://github.com/Kong/go-kong/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/Kong/go-kong/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/Kong/go-kong/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/Kong/go-kong/compare/v0.7.0...v0.8.0
[v0.7.0]: https://github.com/Kong/go-kong/compare/v0.6.2...v0.7.0
[v0.6.2]: https://github.com/Kong/go-kong/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/Kong/go-kong/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/Kong/go-kong/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/Kong/go-kong/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/Kong/go-kong/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/Kong/go-kong/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/Kong/go-kong/compare/0.3.0...v0.4.0
[0.3.0]: https://github.com/Kong/go-kong/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/Kong/go-kong/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/Kong/go-kong/compare/87666c7fe73477d1874d35d690301241cd23059f...0.1.0
