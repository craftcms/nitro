<p><img src="resources/craft-nitro.svg" width="60" height="60" alt="Craft Nitro icon" /></p>

# Craft Nitro

Nitro is a command-line tool focused on making local Craft development quick and easy. Nitro’s one dependency is [Multipass](https://multipass.run/), which allows you to create Ubuntu virtual machines.

---

## Table of Contents

- [Requirements](#requirements)
- [What’s Included](#whats-included)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Commands](#commands)

---

## What’s Included

Nitro installs the following on every machine:

- PHP 7.4 (+ option to use 7.3, 7.2, 7.1, 7.0)
- MySQL
- PostgreSQL
- Redis
- Xdebug
- Blackfire

> Note: For a more detailed writeup on how to configure Xdebug and Nitro with PhpStorm, view this document on [how to configure Xdebug and PhpStorm for both web and console debugging](XDEBUG.md).

## Installation

1. Install [Multipass](https://multipass.run).
2. Run this terminal command:

    ```bash
    curl https://raw.githubusercontent.com/craftcms/nitro/master/get.sh | bash
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

This works like you might expect, it will create a new machine named `diesel` with 2 CPUs and 4GB of RAM. In addition, it will create two database "servers" inside the virtual machine for MySQL and PostgreSQL and make the assigned ports available on the machine's IP address (i.e. not localhost). 

> Note: Nitro can run multiple versions of the same database engine (e.g. PostgreSQL 11 and 12) because it utilizes Docker underneath. See [this file](examples/nitro-multiple-versions.yaml) for an example.

```bash
nitro machine create
```

> Note: If you run `nitro machine create` and it cannot locate the `nitro.yaml` it will walk you through setting up the machine.

After running `machine create`. The bootstrap process will install the latest PHP version, MySQL, Postgres, and Redis from the `nitro.yaml` file.

The next step is to add a new site to the machine:

```bash
cd /Users/jasonmccallister/dev
$ nitro add my-project
→ What should the hostname be? [my-project] $ myproject.test
→ Where is the webroot? [web] $
myproject.test has been added to nitro.yaml.
→ apply nitro.yaml changes now? [y] $ n
You can apply new nitro.yaml changes later by running `nitro apply`.
```

> Note: You can use any top-level domain you wish, but we recommend using .test

This process will perform the following tasks:

1. Set up a new nginx virtual server for `myproject.test`.
2. Attach the directory `/Users/jasonmccallister/dev/my-project` to that virtual server.
3. Edit your `/etc/hosts` file to point `myproject.test` to the virtual server for use locally.

You can now visit `http://myproject.test` in your browser!

## Commands

The following commands will help you manage your virtual server.

- [`apply`](#apply)
- [`add`](#add)
- [`context`](#context)
- [`edit`](#edit)
- [`info`](#info)
- [`import`](#import)
- [`logs`](#logs)
- [`machine create`](#machine-create)
- [`machine destroy`](#machine-destroy)
- [`redis`](#redis)
- [`self-update`](#self-update)
- [`ssh`](#ssh)
- [`stop`](#stop)
- [`update`](#update)
- [`version`](#version)
- [`xdebug configure`](#xdebug-configure)
- [`xdebug on`](#xdebug-on)
- [`xdebug off`](#xdebug-off)

> Note: These examples use a custom config file `nitro-example.yaml`. If you’d like to use Nitro’s default server name (`nitro-local`), you can skip adding the `--machine` argument.

### `apply`

`apply` will look at a config file and make changes from the mounts and sites in the config file by adding or removing them. The config file is the source of truth for your Nitro machine.

```bash
nitro apply
```

```bash
$ nitro apply
There are 2 mounted directories and 1 new mount(s) in the config file.
Applied changes from nitro.yaml.
```

### `add`

Add will create an interactive prompt to add a site (and mount it) into your Nitro machine. By default, it will look at your current working directory and assume that it is a Craft project.

```bash
cd /Users/brandon/Sites/example.test
$ nitro add
→ What should the hostname be? [example.test] $ ex.test
→ Where is the webroot? [web] $
ex.test has been added to nitro.yaml.
→ apply nitro.yaml changes now? [y] $ n
You can apply new nitro.yaml changes later by running `nitro apply`.
```

You can optionally pass a path to the directory as the first argument to use that directory:

```bash
cd /Users/brandon/Sites/
$ nitro -f nitro.yaml add demo-site
✔ What should the hostname be? [demo-site]: $
Where is the webroot? [web]: $
✔ apply nitro.yaml changes now? [y]: $
Applied the changes and added demo-site to nitro.  
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

Edit allows you to quickly open your nitro.yaml file to manually make changes. However, it is recommended to use `nitro` commands to edit your config.

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

### `import`

Import allows you to import a SQL file into a database. You will be prompted with a list of running database engines (mysql and postgres) to import the file into.

```shell
$ nitro import mybackup.sql
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select database:
  ▸ mysql_5.7_3306
```

### `logs`

Views virtual machines logs. This command will prompt you for a type of logs to view (e.g. `nginx`, `database`, or `docker` (for a specific container)). 

```​shell
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

> Note: By default, Multipass does not permanently delete a machine and can cause name conflicts (e.g. `instance "nitro-local" already exists`). This will not affect any local files or directories attached to the machine.

Options:

- `--permanent` permanently deletes a machine **(this is non-recoverable!)**

To soft-destroy the `diesel` machine, so it is recoverable later:

```bash
nitro machine destroy
```

To **permanently** destroy the `diesel` machine:

```bash
nitro machine destroy --permanent
```

### `redis`

Access a Redis shell.

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

Perform updates to the Nitro CLI.

```bash
nitro self-update
```

### `ssh`

Nitro gives you full root access to your virtual server. The default user is `ubuntu` and has `sudo` permissions without a password. Once you’re in the virtual server, you can run `sudo` commands as usual (e.g. `sudo apt install golang`).

```bash
nitro ssh
```

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

This upgrades the `diesel` machine’s software packages to their newest versions:

```bash
nitro update
```

### `version`

Checks the currently version of nitro against the releases and shows any updated versions.  

```bash
nitro version
```

### `xdebug configure`

Configures Xdebug for remote access and debugging with PhpStorm or other IDEs.

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
