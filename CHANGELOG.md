# Release Notes for Nitro

## Unreleased

### Added
- Added a `self-update` command.
- Add a Multipass check to get.sh.

## 0.9.1 - 2020-04-18

### Added
- Added `edit` command to edit a nitro.yaml config.  ([#70](https://github.com/pixelandtonic/nitro/issues/70))
- Added `logs` command to check `nginx`, `database`, and `docker` logs.

### Fixed
- `apply` now checks if sites are setup in the machine and configures them if they are missing.

## 0.7.5 - 2020-04-12

### Added
- Added checksum support for `get.sh` when downloading and updating.  ([#56](https://github.com/pixelandtonic/nitro/issues/56))