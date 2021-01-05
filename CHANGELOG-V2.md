# Release Notes for Craft Nitro

## Unreleased

### Added

- Nitro now creates HTTPS sites by default.
- You can now set a PHP version per site.
- `xdebug on` and `xdebug off` apply to only one site.
- Added the `nitro trust` command to import and trust the certificates used for sites.
- Added the `nitro composer` with the `--version` flag to install or update projects without composer installed locally.
- Added the `nitro npm` with the `--version` flag to install or update projects without node/npm installed locally.
- Added the `nitro ls` command to show all the running containers.
- Added the `nitro prep` command to help setup a project to deploy to a Docker based environment by adding a `Dockerfile` based on the sites PHP version, a `.dockerignore`, and multi-stage builds for handling dependency installs in Docker.
- You can now set the default ports for HTTP, HTTPS, and the API to avoid any port collisions using `NITRO_HTTP_PORT`, `NITRO_HTTPS_PORT`, and `NITRO
_API_PORT`.
- Added the `version` command to display the cli and gRPC API versions
- Added the `nitro craft` command to run console commands in the sites container.

### Changed

- Multipass is no longer a dependency. Docker is now the only dependency and setting up a development environment is faster.
- Nitro will now check for port collisions during `init` and when adding database engines.
- Terminal output now has colors to identify info and error output.
- The `nitro create` command now accepts custom GitHub repositories and installs composer and node dependencies automatically using `nitro composer` and `nitro npm`.
- Machine name is now the environment name, the `NITRO_DEFAULT_MACHINE` has been renamed to `NITRO_DEFAULT_ENVIRONMENT`.
- The `create` command now accepts Github URLs as an argument to allow you to build your own boilerplate.
- The `create` command will now install composer dependencies.
- PHP settings are now site specific.
- the `ssh` command allows you to ssh as the root user with the `--root` flag
