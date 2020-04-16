<p><img src="resources/craft-nitro.svg" width="60" height="60" alt="Craft Nitro icon" /></p>

# Craft Nitro

A better, faster way to develop Craft CMS apps locally without Docker or Vagrant! Nitro’s one dependency is Multipass, which allows you to create Ubuntu virtual machines.

![Go test](https://github.com/pixelandtonic/nitro/workflows/Go%20test/badge.svg)

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

## Installation

TODO

## Usage

In order to create a new development server, you must create a new Nitro machine. By default, this will not attach any directories and is equivalent to getting a brand new Virtual Private Server (VPS).

Nitro defaults to a global `nitro.yaml`. The default location is `~/.nitro/nitro.yaml`. However, you can override the configuration for each command by providing a `--config` (or the shorthand `-f`) with the path to the file (e.g. `nitro -f /path/to/nitro.yaml <command>`).


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
- [`info`](#info)
- [`logs`](#logs)
- [`ssh`](#ssh)
- [`stop`](#stop)
- [`machine create`](#destroy)
- [`machine destroy`](#destroy)
- [`redis`](#redis)
- [`restart`](#restart)
- [`self-update`](#self-update)
- [`update`](#update)

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

### `machine create`

Creates a new server. The following options are available:

| Argument        | Default | Options                           | Description                                       |
|-----------------|---------|-----------------------------------|---------------------------------------------------|
| `--php-version` | `7.4`   | `7.4`, `7.3`, `7.2`, `7.1`, `7.0` | Specifies PHP version used for bootstrap command. |
| `--cpus`        | `2`     | max host CPUs\*                   | Number of CPUs to allocate to the server.         |
| `--memory`      | `2G`    | max host memory\*                 | Gigabytes of memory to allocate to the server.    |
| `--disk`        | `20G`   | max host disk\*                   | Disk space to allocate to the server.             |

<small>\*: CPU, memory, and disk are shared with the host—not reserved—and represent maximum resources to be made available.</small>

### `site`

Adds a new virtual host to nginx and mounts a local directory to the server.

> Note: if you pass a version of PHP that was _not_ used when provisioning the server, Nitro will install that version of PHP for you.

This adds a host using `mysite.test` to the `diesel` machine, using PHP 7.4 and a document root of `/Users/jason/Sites/craftcms`.

```bash
nitro site --php-version 7.4 --path /Users/jason/Sites/craftcms mysite.test
```

> Note: The `--path` argument is only required if you are not mounting the current working directory. If you are in a current Craft CMS project, you can omit the `--path`.

```bash
nitro site mysite.test
```

To remove a site from the virtual machine, run the following command:

```bash
nitro site --remove mysite.test
```

### `attach`

Attaches, or mounts, a local directory to an nginx server’s web root.

This mounts the local directory `/Users/jason/Sites/craftcms` as the web root for the `diesel` machine’s `mysite.test` host:

```bash
nitro attach mysite.test /Users/jason/Sites/craftcms
```

### `ssh`

Nitro gives you full root access to your virtual server. The default user is `ubuntu` and has `sudo` permissions without a password. Once you’re in the virtual server, you can run `sudo` commands as usual (e.g. `sudo apt install golang`).

This launches a new shell within the `diesel` machine:

```bash
nitro ssh
```

### `xon`

Enables Xdebug, which is installed and disabled by default on each machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.3 and enables it for the `diesel` machine:

```bash
nitro xon --php-version 7.3
```

### `xoff`

Disables Xdebug on a machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.2 but disables it for the `diesel` machine:

```bash
nitro xoff --php-version 7.2
```

### `start`

Starts, or turns on, a machine.

This turns on the `diesel` machine:

```bash
nitro start
```

### `stop`

Stops, or turns off, a machine.

This turns off the `diesel` machine:

```bash
nitro stop
```

### `destroy`

Destroys a machine. By default, Multipass does not permanently delete a machine and can cause name conflicts (e.g. `instance "nitro-dev" already exists`). This will not affect any files or directories attached to the machine.

Options:

- `--permanent` permanently deletes a machine **(this is non-recoverable!)**

This soft-destroys the `diesel` machine:

```bash
nitro destroy --permanent
```

This **permanently** destroys the `diesel` machine:

```bash
nitro destroy --permanent true
```

### `sql`

Launches a database shell as the root user.

- `--postgres` access the PostgreSQL shell rather than MySQL (default)

This launches a PostgreSQL console shell for the `diesel` machine:

```bash
nitro sql --postgres
```

### `redis`

Access a Redis shell.

This launches a Redis console shell for the `diesel` machine:

```bash
nitro redis
```

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

This upgrades the `diesel` machine’s software packages to their newest versions:

```bash
nitro update
```

### `logs`

Views the virtual machines logs.

Options:

- `nginx`
- `xdebug`

This displays nginx logs for the `diesel` machine:

```bash
nitro logs nginx
```

This displays xdebug logs for the `diesel` machine:

```bash
nitro logs xdebug
```

### `services`

Stops, starts, or restarts services on a virtual machine. The commands are nested under the `services` command, so when calling `nitro services` the sub commands will be listed.

Options:

- `--nginx`
- `--mysql`
- `--postgres`
- `--redis`

This restarts nginx and MySQL on the `diesel` machine:

```bash
nitro services restart --nginx --mysql
```

This stops PostgreSQL on the `diesel` machine:

```bash
nitro services stop --postgres
```

This starts Redis on the `diesel` machine:

```bash
nitro services start --redis
```

### Xdebug

#### Debugging web requests

Once you have a machine created, you can run `nitro xdebug on` to enable Xdebug
support in your machine.

Install the Xdebug browser helper in your favorite browser.

* [Chrome](https://chrome.google.com/extensions/detail/eadndfjplgieldjbigjakmdgkmoaaaoc)

* [Firefox](https://addons.mozilla.org/en-US/firefox/addon/xdebug-helper-for-firefox/)

* [Internet Explorer](https://www.jetbrains.com/phpstorm/marklets/)

* [Safari](https://github.com/benmatselby/xdebug-toggler)

* [Opera](https://addons.opera.com/addons/extensions/details/xdebug-launcher/)

Go to the Xdebug browser helper options, choose "PhpStorm" and save.

![Xdebug Browser Helper Chrome](resources/xdebug_chrome_settings.png)

Create a new server in PhpStorm using your machine's domain name.

![PhpStorm Server Settings](resources/phpstorm_server.png)

Setup path mappings to that `/app/sites/example.test` in your Nitro machine is
mapped to your project's root on your host machine.

Create a new "PHP Remote Debug" configuration and select the server you just created.

Check "Filter debug connection by IDE key" and enter "PHPSTORM" for the IDE key.

![PhpStorm Remote Debug Settings](resources/phpstorm_remote_debug.png)

Click the "Start Listening for PHP Debug Connections" button in PhpStorm.

![PhpStorm Remote Debug Settings](resources/start_listening.png)

Click the "Debug" button on your browser's Xdebug helper.

![PhpStorm Remote Debug Settings](resources/xdebug_chrome.png)

Then load the site in your browser and whatever breakpoints you've set will be hit.

#### Debugging PHP console requests

Do everything above except Xdebug browser helper. SSH into your Nitro machine using
`nitro ssh`, then run your PHP script from the console and any breakpoints you've
set will be hit.
