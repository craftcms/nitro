# Release Notes for Craft Nitro

## Unreleased

### Added
- Added the `nitro alias` command to quickly setup alias domains for a site.

## 2.0.0-alpha - 2021-02-02

### Added
- Added support for HTTPS on a per-site basis. [#10](https://github.com/craftcms/nitro/issues/10)
- Added support for PHP 8.
- Added support for Xdebug 3 for PHP 7.2 and higher.
- Added the `nitro trust` command to import and trust the certificates used for sites.
- Added the `nitro composer` command with the `--version` flag to install or update projects without Composer installed locally.
- Added the `nitro npm` command with the `--version` flag to install or update projects without node/npm installed locally.
- Added the `nitro version` command to display the cli and gRPC API versions
- Added the `nitro clean` command to remove unused containers for composer/npm
- Added the `nitro craft` command to run `craft` console commands in a site’s container. [#189](https://github.com/craftcms/nitro/issues/189)
- Added the `nitro db ssh` command, which allows you to SSH into a database container.
- Added `nitro enable` and `nitro disable` commands which let you enable and disable services.
- Added the `nitro queue` command, which allows you to listen to Craft’s queue. [#189](https://github.com/craftcms/nitro/issues/189)
- Added the `nitro validate` command, which validates your `~/.nitro/nitro.yml` file.
- You can now set the default ports for HTTP, HTTPS, and the API to avoid any port collisions using `NITRO_HTTP_PORT`, `NITRO_HTTPS_PORT`, and `NITRO
  _API_PORT`.
  - Added the `nitro portcheck` command to quickly check if a port is available.
  - Added `nitro share` to quickly share sites with ngrok. [#2](https://github.com/craftcms/nitro/issues/189)
  - Added the `nitro iniset` command to walk user through setting PHP settings for a site.
  - Added the `nitro extensions` command to walk users through setting custom PHP extensions to install during `nitro apply`.

### Changed
- Nitro now requires Docker, instead of Multipass. [#224](https://github.com/craftcms/nitro/issues/224) [#222](https://github.com/craftcms/nitro/issues/222) [#215](https://github.com/craftcms/nitro/issues/215) [#205](https://github.com/craftcms/nitro/issues/205) [#182](https://github.com/craftcms/nitro/issues/182) [#181](https://github.com/craftcms/nitro/issues/181) [#180](https://github.com/craftcms/nitro/issues/180) [#152](https://github.com/craftcms/nitro/issues/152) [#22](https://github.com/craftcms/nitro/issues/22) [#18](https://github.com/craftcms/nitro/issues/18) [#216](https://github.com/craftcms/nitro/issues/216)
- Greatly improved Nitro support on Windows.
- Nitro now has a single `~/.nitro/nitro.yaml` file to manage everything, instead of a YAML file per machine like in v1.
- Nitro will now check for port collisions during `init` and when adding database engines.
- The `nitro create` command now accepts custom GitHub repositories and installs Composer and node dependencies automatically using `nitro composer` and `nitro npm`. [#101](https://github.com/craftcms/nitro/issues/101)
- The `create` command now accepts Github URLs as an argument to allow you to build your own boilerplate.
- The `create` command will now install composer dependencies.
- PHP versions and settings are now applied on a per-site basis. [#200](https://github.com/craftcms/nitro/issues/200) [#105](https://github.com/craftcms/nitro/issues/105) [#225](https://github.com/craftcms/nitro/issues/225)
- Xdebug is now applied on a per-site basis.
- The `ssh` command allows you to ssh as the root user with the `--root` flag.
- Most Nitro commands are now context aware of the directory they are executed in. [#167](https://github.com/craftcms/nitro/issues/167)
- The `apply` command will only update the hosts file if it has changed. [#117](https://github.com/craftcms/nitro/issues/117)
- Site hostnames and aliases are set on the containers hosts file. [#150](https://github.com/craftcms/nitro/issues/150)
