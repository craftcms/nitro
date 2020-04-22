# Release Notes for Nitro

## Unreleased
 
 ### Changes
 - `remove` will not prompt for sites to remove if it cannot find any sites in the config file.
 - `remove` will now remove the machines sites from the hosts file. 
 - `apply` will now create any new databases that it finds in the config file.
 - `destroy` command is now always permanent and the `--permanent` flag has been removed. `destroy` is no longer nested under the machine command.
 - Fixed an error where a user could select MySQL version 5.8.
 
### Fixed
- Fixed and issue where users could get a segfault when adding a site ([#78](https://github.com/craftcms/nitro/issues/78))
- Fixed (again) an issue when importing database backups using relative paths and added more tests ([#75](https://github.com/craftcms/nitro/issues/75))
- Fixed an error when running `xdebuf off` where php-fpm would not restart
- Renamed `get.sh` to `install.sh`

## 0.9.3 - 2020-04-20

### Fixed
- Fixed an issue when importing database backups using relative paths ([#75](https://github.com/craftcms/nitro/issues/75))

## 0.9.2 - 2020-04-20

### Added
- Added `import` command to let users import a database backup from their system into nitro. ([#1](https://github.com/craftcms/nitro/issues/1))

### Changed
- Nitro now installs the PHP SOAP extension by default.
- Nitro will now walk users through the creation of a machine when no config file is present ([#44](https://github.com/craftcms/nitro/issues/44)).
- Nitro now prompts you to modify your hosts file after using `add` ([#40](https://github.com/craftcms/nitro/issues/40)).

## 0.9.1 - 2020-04-18

### Added
- Added `edit` command to edit a nitro.yaml config.  ([#70](https://github.com/craftcms/nitro/issues/70))
- Added `logs` command to check `nginx`, `database`, and `docker` logs.

### Fixed
- `apply` now checks if sites are setup in the machine and configures them if they are missing.

## 0.7.5 - 2020-04-12

### Added
- Added checksum support for `get.sh` when downloading and updating.  ([#56](https://github.com/craftcms/nitro/issues/56))
