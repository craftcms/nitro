<p><img src="resources/craft-nitro.svg" width="60" height="60" alt="Craft Nitro icon" /></p>

# Craft Nitro

A better, faster way to develop Craft CMS apps locally without Docker or Vagrant! Nitro’s one dependency is Multipass, which allows you to create Ubuntu virtual machines.

---

## Table of Contents

- [Requirements](#requirements)
- [What’s Included](#whats-included)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)

---

## Requirements

- [Multipass](https://multipass.run)

## What’s Included

Nitro installs the following on every machine:

- PHP 7.4 (+ option to use 7.3, 7.2, 7.1, 7.0)
- MySQL
- PostgreSQL (optional)
- Redis
- Xdebug
- Blackfire

> note: for a more detailed writeup on how to configure xdebug and nitro with PHPStorm, view this document on [how to configure xdebug and PHPStorm for both web and console debugging](XDEBUG.md).

## Installation

```
curl https://raw.githubusercontent.com/pixelandtonic/nitro/develop/get.sh | sudo sh
```

## Getting Started

In order to create a new development server, you must create a new Nitro machine. By default, this will not attach any directories and is equivalent to getting a brand new Virtual Private Server (VPS).

Nitro defaults to a global `nitro.yaml`. The default location is `~/.nitro/nitro.yaml`. However, you can override the configuration for each command by providing a `--config` (or the shorthand `-f`) with the path to the file (e.g. `nitro -f /path/to/nitro.yaml <command>`). Here is an example config:

```yaml
name: diesel
php: "7.4"
cpus: "2"
disk: 40G
memory: 4G
databases:
- engine: mysql
  version: "5.7"
  port: "3306"
- engine: postgres
  version: "12"
  port: "5432"
```

This works like you might expect, it will create a new machine named `diesel` with 2 CPUs and 4GB of RAM. In addition, it will create to database "servers" inside the virtual machine for MySQL and PostgreSQL and make the assigned ports available on the machines IP address (not your localhost). 

> Note: Nitro can run multiple versions of the same database engine (e.g. PostgreSQL 11 and 12) because it utilizes Docker underneath. See [this file](examples/nitro-multiple-versions.yaml) for an example.

```bash
nitro machine create
```

> Note: `nitro machine create` has options you can pass when creating a new server. However, we can set some "“sane” defaults for most scenarios. To view the options, run `nitro machine create --help`.

After running `machine create`. The bootstrap process will install the latest PHP version, MySQL, Postgres, and Redis from the `nitro.yaml` file.

The next step is to add a new virtual host to the server:

```bash
nitro site add /Users/jason/Sites/craftcms myclientsite.test
```

> Note: you can use any top level domain you wish, but we recomend using .test

This process will perform the following tasks:

1. Set up a new nginx virtual server for `myclientsite.test`.
2. Attach the directory `/Users/jason/Sites/craftcms` to that virtual server.
3. Edit your `/etc/hosts` file to point `myclientsite.test` to the virtual server for use locally.

You can now visit `http://myclientsite.test` in your browser!

## Commands

The following commands will help you manage your virtual server.

- [`apply`](#apply)
- [`add`](#add)
- [`context`](#context)
- [`edit`](#edit)
- [`info`](#info)
- [`logs`](#logs)
- [`mount`](#mount)
- [`machine create`](#machine-create)
- [`machine destroy`](#machine-destroy)
- [`redis`](#redis)
- [`self-update`](#self-update)
- [`ssh`](#ssh)
- [`stop`](#stop)
- [`update`](#update)
- [`xdebug configure`](#xdebug-configure)
- [`xdebug on`](#xdebug-on)
- [`xdebug off`](#xdebug-off)
- [`version`](#version)

> Note: these examples use a custom config file `nitro-example.yaml`. If you’d like to use Nitro’s default server name (`nitro-dev`), you can skip adding the `--machine` argument.

### `apply`

Apply will look at a config file and make changes from the mounts, and sites in the config file by adding or removing. The config file is the source of truth for your nitro machine.

```bash
nitro apply
```

```bash
$ nitro apply
ok, there are 2 mounted directories and 1 new mount(s) in the config file
applied changes from nitro.yaml
```

### `add`

Add will create an interactive prompt to add a site (and mount it) into your nitro machine. By default, it will look at your current working directory and assume that it is a Craft project.

```bash
cd /Users/brandon/Sites/example.test
$ nitro add
→ what should the hostname be? [example.test] $ ex.test
→ what is the webroot? [web] $
ex.test has been added to nitro.yaml
→ apply nitro.yaml changes now? [y] $ n
ok, you can apply new nitro.yaml changes later by running `nitro apply`.
```

You can optionally pass a path to the directory as the first argument to use that directory:

```bash
cd /Users/brandon/Sites/
$ nitro -f nitro.yaml add demo-site
✔ what should the hostname be? [demo-site]: $
what is the webroot? [web]: $
✔ apply nitro.yaml changes now? [y]: $
ok, we applied the changes and added demo-site to nitro  
````

| Argument     | Default                                        | Options | Description                                 |
|--------------|------------------------------------------------|---------|---------------------------------------------|
| `--hostname` | (the current working directory name)           |         | The hostname to use for accessing the site. |
| `--webroot`  | (looks for web, public, public_html, and www)) |         | The webroot to configure nginx.             |

### `context`

Shows the currently used configuration file for quick reference.

```shell
$ nitro -f nitro-example.yaml context
Using config file: nitro.yaml
------
name: nitro
php: "7.4"
cpus: "1"
disk: 40G
memory: 4G
mounts:
- source: ~/sites/demo-site
  dest: /nitro/sites/demo-site
databases:
- engine: mysql
  version: "5.7"
  port: "3306"
- engine: postgres
  version: "12"
  port: "5432"
sites:
- hostname: demo-site
  webroot: /nitro/sites/demo-site/web
------
```

### `edit`

Edit allows you to quickly open your nitro.yaml to make changes. However, it is recommended to use `nitro` commands to edit your config.

```shell
nitro edit
```

### `info`

Shows the _running_ information for a machine like the IP address, memory, disk usage, and mounts.

```shell
$ nitro info
Name:           nitro
State:          Running
IPv4:           192.168.64.48
Release:        Ubuntu 18.04.4 LTS
Image hash:     2f6bc5e7d9ac (Ubuntu 18.04 LTS)
Load:           0.09 0.15 0.22
Disk usage:     2.7G out of 38.6G
Memory usage:   379.8M out of 3.9G
Mounts:         /Users/jasonmccallister/sites/demo-site => /nitro/sites/demo-site
                    UID map: 501:default
                    GID map: 20:default
```

### `logs`

Views virtual machines logs. This command will prompt you for a type of logs to view (e.g. `nginx`, `database`, or `docker` (for a specific container)). 

```bash
nitro logs
```

### `machine create`

Creates a new server. The following options are available:

| Argument        | Default | Options                           | Description                                       |
|-----------------|---------|-----------------------------------|---------------------------------------------------|
| `--php-version` | `7.4`   | `7.4`, `7.3`, `7.2`, `7.1`, `7.0` | Specifies PHP version used for bootstrap command. |
| `--cpus`        | `2`     | max host CPUs\*                   | Number of CPUs to allocate to the server.         |
| `--memory`      | `2G`    | max host memory\*                 | Gigabytes of memory to allocate to the server.    |
| `--disk`        | `20G`   | max host disk\*                   | Disk space to allocate to the server.             |

<small>\*: CPU, memory, and disk are shared with the host—not reserved—and represent maximum resources to be made available.</small>


### `machine destroy`

Destroys a machine.

> Note: by default, Multipass does not permanently delete a machine and can cause name conflicts (e.g. `instance "nitro-dev" already exists`). This will not affect any local files or directories attached to the machine.

Options:

- `--permanent` permanently deletes a machine **(this is non-recoverable!)**

This soft-destroys the `diesel` machine:

```bash
nitro machine destroy
```

This **permanently** destroys the `diesel` machine:

```bash
nitro machine destroy --permanent
```

### `mount`

Mounts a local directory to a path on the machine.

```bash
nitro mount ~/sites/project-folder /home/ubuntu/project-folder
```

### `redis`

Access a Redis shell.

This launches a Redis console shell for the `diesel` machine:

```bash
nitro redis
```

### `start`

Starts, or turns on, a machine.

```bash
nitro start
```

### `stop`

Stops, or turns off, a machine.

```bash
nitro stop
```

### `self-update`

Perform updates to the nitro CLI.

```bash
nitro self-update
```

### `ssh`

Nitro gives you full root access to your virtual server. The default user is `ubuntu` and has `sudo` permissions without a password. Once you’re in the virtual server, you can run `sudo` commands as usual (e.g. `sudo apt install golang`).

```bash
nitro ssh
```

### `xdebug configure`

Configures Xdebug for remote access and debugging with PHPStorm or other IDE.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

```bash
nitro xdebug configure --php-version 7.3
```

### `xdebug on`

Enables Xdebug, which is installed and disabled by default on each machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.3 and enables it:

```bash
nitro xdebug on --php-version 7.3
```

### `xdebug off`

Disables Xdebug on a machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.2 but disables it:

```bash
nitro xdebug off --php-version 7.2
```

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

This upgrades the `diesel` machine’s software packages to their newest versions:

```bash
nitro update
```
