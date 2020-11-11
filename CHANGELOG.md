# Release Notes for Craft Nitro

## 1.1.1 - 2020-11-11

### Added
- Added support for PHP 8.0.

## 1.1.0 - 2020-10-06

### Added
- Added the `nitro create` command, which will set up a new Craft installation without PHP or Composer installed locally. ([#101](https://github.com/craftcms/nitro/issues/101))
- Added the `--silent` flag to the `xon`, `xoff`, and `php iniset` commands.

### Changed
- The `db import` command can now import zip and gzip files.
- The `db import` command will now detect the database backup type and automatically select the appropriate database engine, if the backup file was uncompressed. ([#132](https://github.com/craftcms/nitro/issues/132))
- The `db import` and `db add` commands now replaces dashes (`-`) with underscores (`_`) in the specified database name, to work around a SQL error. ([#212](https://github.com/craftcms/nitro/issues/212))
- The `php iniset` command can now set the `display_errors` config setting. ([#172](https://github.com/craftcms/nitro/issues/172))
- The `display_errors` config setting is now set to `On` by default.

### Fixed
- Fixed a bug where `php iniset memory_limit` was listed as an available command.
- Fixed a bug where `php iniset` commands weren’t always setting the correct values. ([#207](https://github.com/craftcms/nitro/issues/207))
- Fixed a bug where `php iniget` commands weren’t always returning the correct values.
- Fixed a bug where the system’s `hosts` file wasn’t getting updated on Linux machines. ([#213](https://github.com/craftcms/nitro/issues/213))

## 1.0.1 - 2020-08-12

### Changed
- Composer is now installed by default on new machines.
- The default Nginx site is now a friendly landing page with helpful links.

### Fixed
- Fixed an error that occurred when checking for the latest version of Nitro. ([#199](https://github.com/craftcms/nitro/issues/199))

## 1.0.0 - 2020-08-11

### Added
- Added the `support` command to quickly create GitHub issues pre-populated with environment info.

### Changed
- The `info` command now displays additional info such as IP, PHP version, and links to common tasks.

### Fixed
- Fixed a bug confirmation prompts would take just about any input as a “yes”. ([#190](https://github.com/craftcms/nitro/issues/190))

## 1.0.0-RC1.1 - 2020-08-07

### Fixed
- Fixed an issue that occurred when creating new machines.

## 1.0.0-RC1 - 2020-08-07

### Added
- Added the `nginx start`, `nginx stop`, and `nginx restart` commands.
- Added the `php start`, `php stop`, and `php restart` commands. ([#57](https://github.com/craftcms/nitro/issues/57))
- Added the `php iniget` command, which can be used to see current php.ini values.
- Added the `php iniset` command, which can be used to modify the `max_execution_time`, `max_input_vars`, `max_input_time`, `upload_max_filesize`, `max_file_uploads`, and `memory_limit` php.ini settings. ([#138](https://github.com/craftcms/nitro/issues/138))
- Added the `xon` and `xoff` commands, which are shortcuts for `xdebug on` and `xdebug off`.
- Added the `nitrod` daemon that runs in the virtual machine and exposes a gRPC API on port 50051.
- Enabling and disabling Xdebug is now performed via the gRPC API, and is much faster now.

### Changed
Improved Composer performance when run from inside the virtual machine. ([#186](https://github.com/craftcms/nitro/issues/186))

### Fixed
- Fixed a bug where the default PHP version was not getting updated when running the `apply` command. ([#192](https://github.com/craftcms/nitro/issues/192))
- Fixed command completion when pressing the <kbd>Tab</kbd> key.

## 1.0.0-beta.10 - 2020-06-04

### Added
- Added support for a `NITRO_EDIT_HOSTS` environment variable so that when set to `false`, Nitro will never edit the host machine’s `hosts` file.

### Changed
- The `destroy` command now has a `--skip-hosts` option.

## 1.0.0-beta.9 - 2020-06-02

### Changed
- The `add` now has a `--skip-hosts` option. ([#163](https://github.com/craftcms/nitro/issues/163))
- The `db add` command now validates the database name. ([#160](https://github.com/craftcms/nitro/issues/160))

### Fixed
- Fixed a bug with the Nginx config template (run `nitro refresh` for the change to take effect).

## 1.0.0-beta.8 - 2020-06-01

### Changed
- Newly-created site configs are now based on the Nginx config provided by <https://github.com/nystudio107/nginx-craft> (run `nitro refresh` for the change to take effect). ([#35](https://github.com/craftcms/nitro/issues/161))

### Fixed
- Fixed an error that could remove mounted sites. ([#162](https://github.com/craftcms/nitro/issues/162))
- Fixed a bug where the `apply` command wasn’t removing deleted sites’ hostnames from the hosts file. ([#161](https://github.com/craftcms/nitro/issues/161))
- Fixed a bug where the `apply` command wasn’t removing deleted sites’ Nginx configurations from the virtual machine.

## 1.0.0-beta.7 - 2020-05-27

### Fixed
- Fixed a bug where keys transferred into the machine did not have the proper permissions. ([#154](https://github.com/craftcms/nitro/issues/154))
- Fixed a bug where the `init` command was not editing the hosts file. ([#155](https://github.com/craftcms/nitro/issues/155)) ([#156](https://github.com/craftcms/nitro/issues/156))
- Fixed a bug where the `db import` command was not working on PostgreSQL.

## 1.0.0-beta.6 - 2020-05-26

### Changed
- Removed an unnecessary debug command.

## 1.0.0-beta.5 - 2020-05-26

### Added
- Added Windows support.
- Added support for MySQL 8.0. ([#97](https://github.com/craftcms/nitro/issues/97))
- The PostgreSQL and MySQL client tools are now installed on new machines. ([#54](https://github.com/craftcms/nitro/issues/54), [#139](https://github.com/craftcms/nitro/issues/139))
- Added the `install postgres`, `install mysql`, `install composer`, and `install mailhog` commands.

### Changed
- New machines now use Ubuntu 20 LTS.
- Renamed the `--no-backups` option to `--skip-backup` for the `destroy` command.
- Composer is no longer installed on machines by default, but can be installed by running `nitro install composer`.
- The `init` command now sets the default CPU count based on the number of CPUs on the host machine.
- MySQL 5 and 8 now use version-specific configuration directories (`/home/ubuntu/.nitro/databases/mysql/conf.d/<version>/`).
- Removed the `xdebug configure` command, and moved its logic into the `xdebug on` command.

### Fixed
- Fixed a bug where Nitro wasn’t removing Nginx server configs when removing sites.
- Fixed a bug where the `apply` command wasn’t removing deleted mounts’ root directories within the machine. ([#96](https://github.com/craftcms/nitro/issues/96))
- Fixed a bug where the `init` command could return an exit code of 100. ([#96](https://github.com/craftcms/nitro/issues/96))
- The OPcache extension is no longer installed by default. ([#129](https://github.com/craftcms/nitro/issues/129))

## 1.0.0-beta.4 - 2020-05-21

### Added
- Added the `keys` command which prompts which keys should be imported into a machine. ([#141](https://github.com/craftcms/nitro/issues/141))
- Added the `--no-backups` flag to `destroy` which will skip database backups.
- Added `completion` commands for `bash` and `zshrc`.

### Changed
- The `destroy` command creates individual databases backups. ([#146](https://github.com/craftcms/nitro/issues/146))
- The `mysql` system database is no longer backed up using `db backup` or `destroy`. ([#147](https://github.com/craftcms/nitro/issues/147))

### Fixed
- Fixed a bug where the `refresh` command was failing silently.
- Fixed a permissions issue. ([#145](https://github.com/craftcms/nitro/issues/145))
- Fixed an issue when importing mysql databases using `db import`.
- Fixed an issue installing composer on new machines. ([#149](https://github.com/craftcms/nitro/issues/149))

## 1.0.0-beta.3 - 2020-05-19

### Added
- Added the `--skip-hosts` option to the `apply` command. ([#110](https://github.com/craftcms/nitro/issues/110))
- The `add` command will now prompt to create a machine if it does not exist. ([#103](https://github.com/craftcms/nitro/issues/103))
- The `init` command will edit the hosts file if sites are present in the config. ([#123](https://github.com/craftcms/nitro/issues/123))
- Added the `db restart`, `db stop`, `db add`, `db remove`, and `db backup` commands. The `import` command has also been renamed to `db import`.
- Added the `refresh` command, which helps keep scripts and configs updated between versions of Nitro.
- Databases now support custom configuration files. ([#133](https://github.com/craftcms/nitro/issues/133))

### Changed
- Nginx is now configured to allow file uploads up to 100MB. ([#126](https://github.com/craftcms/nitro/issues/126))
- Databases are now backed up automatically when a machine is destroyed. ([#136](https://github.com/craftcms/nitro/issues/136))
- When creating a new machine, the `DB_USER` and `DB_PASSWORD` are automatically set in the environment. ([#119](https://github.com/craftcms/nitro/issues/119))
- The default database is now called `nitro` for MySQL engines, to be consistent with PostgreSQL.
- The `destroy` command now prompts for confirmation. ([#116](https://github.com/craftcms/nitro/issues/116))
- The `init` command now prompts to initialize a new machine if there is no config file.
- The `add` command now mounts project files in `/home/ubuntu/sites/<name>` instead of `/nitro/sites/<name>`. ([#134](https://github.com/craftcms/nitro/issues/134))
- The `apply` command now provides more information. ([#95](https://github.com/craftcms/nitro/issues/95))
- The `init` command now checks if the machine already exists before prompting for input ([#102](https://github.com/craftcms/nitro/issues/102))

### Fixed
- Fixed a bug where the `self-update` command would non-interactively prompt to initialize the primary machine. ([#115](https://github.com/craftcms/nitro/issues/115))
- Fixed a bug where `import --help` didn’t show the required SQL file argument in the usage example. ([#124](https://github.com/craftcms/nitro/issues/124))
- Fixed a bug where the `apply` command wasn’t applying changes to sites’ webroots. ([#113](https://github.com/craftcms/nitro/issues/113))

## 1.0.0-beta.2 - 2020-05-06

### Fixed
- Fixed an error when using the select prompt. ([#104](https://github.com/craftcms/nitro/issues/104))

## 1.0.0-beta.1 - 2020-05-05

### Changed
- Improved the `init` command workflow.
- Changed the input package to use [pixelandtonic/prompt](https://github.com/pixelandtonic/prompt).

### Removed
- Removed MySQL 8.0 support for now.
- Removed PHP 7.0 and 7.1 support.

### Fixed
- Fixed a potential permission error when installing/updating.

## 0.11.4 - 2020-05-04

### Fixed
- Fixed a broken test which prevented a release.

## 0.11.3 - 2020-05-04

### Changed
- The [GMP](https://www.php.net/manual/en/book.gmp.php) and [BCMath](https://www.php.net/manual/en/book.bc.php) PHP extensions are now installed by default.
- Composer is now installed globally on machines. ([#92](https://github.com/craftcms/nitro/issues/92))

### Fixed
- Fixed a permission error when provisioning a PostgreSQL database.
- Fix a bug where the `import` command wasn’t importing.
- Fixed an issue where the machine DNS was not resolving in some environments. ([#91](https://github.com/craftcms/nitro/issues/91))
- Fixed an error when trying to create a database during PostgreSQL import. ([#94](https://github.com/craftcms/nitro/issues/94))

## 0.11.2 - 2020-04-09

### Changed
- The `init` command now prompts for how many CPU cores should be assigned to the machine.

## 0.11.0 - 2020-04-29

### Added
- Added the `rename` command to allow users to quickly rename sites.

### Changed
- The `destroy` command now has a `--clean` option which will delete a config file after destroying the machine.
- The `nitro` database user now has root privileges for `mysql` and `postgres` databases. ([#79](https://github.com/craftcms/nitro/issues/79))
- Added the `php` option back to the config file.
- All commands that perform config changes (e.g. `add`, `remove`, and `rename`) now use the same logic as the `apply` command.
- When importing a database using the `import` command, users will be prompted for the database name which will be created if it does not exist.
- The `apply` command will automatically update the machine's hosts file.
- The `destroy` command will now remove any sites in the machine config from the hosts file.
- The `init` command will use an existing config file and recreate the entire environment.
- Commands now output more _statuses_ where possible to provide the user more feedback.

### Fixed
- When using the `add` command, the config file checks for duplicate sites and mounts. ([#86](https://github.com/craftcms/nitro/issues/86))
- Fixed an issue when using some commands on Windows. ([#88](https://github.com/craftcms/nitro/issues/88))
- Fixed an issue in the `apply` command that would not detect new changes to the config file.

## 0.10.0 - 2020-04-23

> **Warning:** This release contains breaking changes. See the [upgrade notes](UPGRADE.md#upgrading-to-nitro-0100)
> for details.

### Added
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
