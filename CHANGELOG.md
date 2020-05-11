# Release Notes for Craft Nitro

### Unreleased

## Added
- Added the `--skip-hosts` option to the `apply` command to skip editing the hosts file. ([#110](https://github.com/craftcms/nitro/issues/110))
- The `add` command will now prompt to create a machine if it does not exist. ([#103](https://github.com/craftcms/nitro/issues/103))
- The `init` command will edit the hosts file if sites are present in the config. ([#123](https://github.com/craftcms/nitro/issues/123))

## Changed
- The `destroy` command now prompts for verification before machine destruction. ([#116](https://github.com/craftcms/nitro/issues/116))
- The `self-update` command no longer prompts for installation upon completion. ([#115](https://github.com/craftcms/nitro/issues/115))
- The `import` command now shows the requirement for the first argument of a database backup. ([#124](https://github.com/craftcms/nitro/issues/124))
- When a new machine is created the `DB_USER` and `DB_PASSWORD` are automatically set in the environment. ([#119](https://github.com/craftcms/nitro/issues/119)) 

### 1.0.0-beta.2 - 2020-05-06

## Fixed
- Fixed an error when using the select prompt. ([#104](https://github.com/craftcms/nitro/issues/104))

### 1.0.0-beta.1 - 2020-05-05

## Changed
- Improved the `init` command workflow.
- Changed the input package to use [pixelandtonic/prompt](https://github.com/pixelandtonic/prompt).

## Removed
- Removed MySQL 8.0 support for now.
- Removed PHP 7.0 and 7.1 support.

## Fixed
- Fixed a potential permission error when installing/updating.

### 0.11.4 - 2020-05-04

## Fixed
- Fixed a broken test which prevented a release.

### 0.11.3 - 2020-05-04

## Changed
- The [GMP](https://www.php.net/manual/en/book.gmp.php) and [BCMath](https://www.php.net/manual/en/book.bc.php) PHP extensions are now installed by default.
- Composer is now installed globally on machines. ([#92](https://github.com/craftcms/nitro/issues/92))

## Fixed
- Fixed a permission error when provisioning a PostgreSQL database.
- Fix a bug where the `import` command wasn’t importing.
- Fixed an issue where the machine DNS was not resolving in some environments. ([#91](https://github.com/craftcms/nitro/issues/91))
- Fixed an error when trying to create a database during PostgreSQL import. ([#94](https://github.com/craftcms/nitro/issues/94))

### 0.11.2 - 2020-04-09

## Changed
- The `init` command now prompts for how many CPU cores should be assigned to the machine.

### 0.11.0 - 2020-04-29

## Added
- Added the `rename` command to allow users to quickly rename sites.

## Changed
- The `destroy` command now has a `--clean` option which will delete a config file after destroying the machine.
- The `nitro` database user now has root privileges for `mysql` and `postgres` databases. ([#79](https://github.com/craftcms/nitro/issues/79))
- Added the `php` option back to the config file.
- All commands that perform config changes (e.g. `add`, `remove`, and `rename`) now use the same logic as the `apply` command.
- When importing a database using the `import` command, users will be prompted for the database name which will be created if it does not exist.
- The `apply` command will automatically update the machine's hosts file.
- The `destroy` command will now remove any sites in the machine config from the hosts file.
- The `init` command will use an existing config file and recreate the entire environment.
- Commands now output more _statuses_ where possible to provide the user more feedback.

## Fixed
- When using the `add` command, the config file checks for duplicate sites and mounts. ([#86](https://github.com/craftcms/nitro/issues/86))
- Fixed an issue when using some commands on Windows. ([#88](https://github.com/craftcms/nitro/issues/88))
- Fixed an issue in the `apply` command that would not detect new changes to the config file.

## 0.10.0 - 2020-04-23

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
