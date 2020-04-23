# Release Notes for Craft Nitro

## Unreleased

> **Warning:** This release contains breaking changes. See the [upgrade notes](UPGRADE.md#upgrading-to-nitro-0100)
> for details.

## Added
- Added the `init` command, which initializes new machines.
- Added the `remove` command, which removes a site from a machine.

### Changed
- All machine configs are now stored saved in `~/.nitro/`.
- All commands now have an `-m` option, which can be used to specify which machine to work with. (The `-f` option
  has also been removed.)
- The `apply` command now creates any new database servers that it finds in the config file.
- The `machine destroy` command has been renamed to `destroy`, and it now permanently destroys the machine (as
  opposed to archiving it), and removes any hostnames added to your hosts file that point to its IP address.
- The `machine restart` command has been renamed to `restart`.
- The `machine start` command has been renamed to `start`.
- The `machine stop` command has been renamed to `stop`.
- Renamed `get.sh` to `install.sh`.

### Removed
- Removed the `machine create` command. Use the new `init` command to create new machines instead.

### Fixed
- Fixed a bug where users could get a segfault when adding a site. ([#78](https://github.com/craftcms/nitro/issues/78))
- Fixed a bug where it wasn’t possible to import databases using relative paths. ([#75](https://github.com/craftcms/nitro/issues/75))
- Fixed a bug where the `machine create` command listed MySQL 5.8 as an option.
- Fixed a bug where php-fpm wouldn’t restart after running the `xdebuf off` command.

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
