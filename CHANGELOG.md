# Release Notes for Nitro

## Unreleased

### Added
- Added `import` command to let users import a database backup from their system into nitro. ([#1](https://github.com/pixelandtonic/nitro/issues/1))

### Changed
- Nitro now installs the PHP SOAP extension by default.
- Nitro will now walk users through the creation of a machine when no config file is present ([#44](https://github.com/pixelandtonic/nitro/issues/44)).
- Nitro now prompts you to modify your hosts file after using `add` ([#40](https://github.com/pixelandtonic/nitro/issues/40)).

## 0.9.1 - 2020-04-18

### Added
- Added `edit` command to edit a nitro.yaml config.  ([#70](https://github.com/pixelandtonic/nitro/issues/70))
- Added `logs` command to check `nginx`, `database`, and `docker` logs.

### Fixed
- `apply` now checks if sites are setup in the machine and configures them if they are missing.

## 0.7.5 - 2020-04-12

### Added
- Added checksum support for `get.sh` when downloading and updating.  ([#56](https://github.com/pixelandtonic/nitro/issues/56))