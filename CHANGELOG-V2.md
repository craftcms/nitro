# Release Notes for Craft Nitro

## Unreleased

### Added
- Added `containers configure` command. [#170](https://github.com/craftcms/nitro/issues/170)

### Changed
- Commands that rely on existing config files now prompt to run the `init` command if there is no config file.

## 2.0.0-beta.3 - 2021-02-05

### Added
- Nitro will now remove entries from your hosts file when running `destroy`. (#235)[https://github.com/craftcms/nitro/issues/235]

### Changed
- Nitro will use MariaDB instead of MySQL on M1 Macs. [#234](https://github.com/craftcms/nitro/issues/234)
- The `add` command will prompt to run `init` if no configuration file is found. [#237](https://github.com/craftcms/nitro/issues/237)

## 2.0.0-beta.2 - 2021-02-04

### Added
- Nitro now supports Arm CPUs, including M1 Macs.

### Changed
- The `share` command now outputs a more helpful message if Ngrok isn’t installed.
- Xdebug is no longer supported for PHP 7.0.

### Fixed
- Fixed a bug where Nitro wasn’t respecting custom service ports.
- Fixed a bug where the `update` command wasn’t immediately updating containers.
- Fixed a bug where the `add` command would overwrite an existing `.env` file. ([#232](https://github.com/craftcms/nitro/issues/232)

## 2.0.0-beta.1 - 2021-02-03

### Added
- Nitro now runs on Docker rather than Multipass. ([#224](https://github.com/craftcms/nitro/issues/224), [#222](https://github.com/craftcms/nitro/issues/222), [#215](https://github.com/craftcms/nitro/issues/215), [#205](https://github.com/craftcms/nitro/issues/205), [#182](https://github.com/craftcms/nitro/issues/182), [#181](https://github.com/craftcms/nitro/issues/181), [#180](https://github.com/craftcms/nitro/issues/180), [#152](https://github.com/craftcms/nitro/issues/152), [#22](https://github.com/craftcms/nitro/issues/22), [#18](https://github.com/craftcms/nitro/issues/18), [#216](https://github.com/craftcms/nitro/issues/216))
- PHP versions and settings are now applied on a per-site basis. ([#200](https://github.com/craftcms/nitro/issues/200), [#105](https://github.com/craftcms/nitro/issues/105), [#225](https://github.com/craftcms/nitro/issues/225))
- Xdebug is now applied on a per-site basis.
- Added support for SSL. ([#10](https://github.com/craftcms/nitro/issues/10))
- Added support for PHP 8.
- Added support for Xdebug 3 when using PHP 7.2 or later.
- Added the `alias` command.
- Added the `clean` command.
- Added the `composer` command.
- Added the `craft` command, which will run `craft` commands within a site’s container. ([#189](https://github.com/craftcms/nitro/issues/189))
- Added the `db ssh` command.
- Added the `enable` and `disable` commands.
- Added the `extensions` command.
- Added the `iniset` command.
- Added the `npm` command.
- Added the `portcheck` command.
- Added the `queue` command. ([#189](https://github.com/craftcms/nitro/issues/189))
- Added the `share` command, which shares a site via Ngrok. ([#2](https://github.com/craftcms/nitro/issues/189))
- Added the `trust` command.
- Added the `validate` command.
- Added the `version` command.

### Changed
- Nitro now has a single `~/.nitro/nitro.yaml` file to manage everything, instead of a YAML file per machine.
- Most Nitro commands are now context aware of the directory they are executed in. ([#167](https://github.com/craftcms/nitro/issues/167))
- The `apply` command will only update the `hosts` file if it has changed. ([#117](https://github.com/craftcms/nitro/issues/117))
- The `create` command now accepts custom GitHub repositories and installs Composer and Node dependencies automatically. ([#101](https://github.com/craftcms/nitro/issues/101))
- The `db import` command can now import database backups that live outside the project directory.
- The `ssh` command now has a `--root` flag, which will SSH into the container as the root user.
- It’s now possible to set the default ports for HTTP, HTTPS, and the API to avoid any port collisions using `NITRO_HTTP_PORT`, `NITRO_HTTPS_PORT`, and `NITRO_API_PORT`.
- Nitro will now check for port collisions during `init` and when adding database engines.
- Sites’ containers’ `hosts` files now list other Nitro site host names. ([#150](https://github.com/craftcms/nitro/issues/150))

### Fixed
- Fixed and error that could occur when downloading database images.
