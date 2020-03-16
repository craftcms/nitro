<p><img src="resources/craft-nitro.svg" width="60" height="60" alt="Craft Nitro icon" /></p>

# Craft Nitro

A better, faster way to develop Craft CMS apps locally without Docker or Vagrant! Nitro’s one dependency is Multipass, which allows you to create Ubuntu virtual machines.

## Requirements

- [Multipass](https://multipass.run)

## What's Included

Nitro installs the following on every machine:

- PHP 7.4 (+ option to use 7.3 or 7.2)
- MySQL (MariaDB)
- PostgreSQL (optional)
- Redis

## Installation

```shell script
composer require --dev craftcms/nitro
```

This package has a single executable named `nitro`. In order to use the CLI, run `./vendor/bin/nitro`.

To use `composer run nitro`, add the following to your `composer.json`’s scripts section:

```json
"scripts": {
  "nitro": "./vendor/bin/nitro"
},
```

Create a new development machine by running the following:

```bash
composer run nitro init
```

This will create a new machine with the default name `nitro-dev`. If you wish to assign another name to the machine, provide one with the `--machine` argument instead:

```bash
./vendor/bin/nitro --machine my-custom-name init
```

## Usage

In order to create a new development server, you must “initialize” Nitro. By default, this will not attach any directories and is equivalent to getting a brand new Virtual Private Server (VPS).

```bash
nitro init
```

> Note: `nitro init` has options you can pass when creating a new server. However, we can set some "“sane” defaults for most scenarios. To view the options, run `nitro init -h`.

After running `init`, the system will default to automatically bootstrap the server. The bootstrap process will install the latest PHP version, MySQL, and Redis.

> Note: if you wish to avoid bootstrapping, pass `--bootstrap false` when running init (e.g. `nitro init --bootstrap false`)

The next step is to add a new virtual host to the server:

```bash
nitro add mysite.test /Users/jason/Sites/craftcms
```

This process will perform the following tasks:

1. Set up a new nginx virtual server for `mysite.test`.
2. Attach the directory `/Users/jason/Sites/craftcms` to that virtual server.
3. Edit your `/etc/hosts` file to point `mysite.test` to the virtual server for use locally.

You can now visit `http://mysite.test` in your browser!

## Commands

The following commands will help you manage your virtual server.

> Note: these examples use a custom server name of `diesel`. If you’d like to use Nitro’s default server name (`nitro-dev`), you can skip adding the `--machine` argument.

### `init`

Creates a new server. The following options are available:

| Argument        | Default | Options             | Description                                       |
| --------------- | ------- | ------------------- | ------------------------------------------------- |
| `--bootstrap`   | `true`  |                     | Bootstraps installation of PHP, MySQL, and Redis. |
| `--php-version` | `7.4`   | `7.4`, `7.3`, `7.2` | Specifies PHP version used for bootstrap command. |
| `--database`    | `mysql` | `mysql`, `postgres` | Specifies database used for bootstrap command.    |
| `--cpus`        | `2`     | max host CPUs\*     | Number of CPUs to allocate to the server.         |
| `--memory`      | `2G`    | max host memory\*   | Gigabytes of memory to allocate to the server.    |
| `--disk`        | `5G`    | max host disk\*     | Disk space to allocate to the server.             |

<small>\*: CPU, memory, and disk are shared with the host—not reserved—and represent maximum resources to be made available.</small>

### `bootstrap`

Installs the specified version of PHP, the database engine, and Redis server onto a server. This should only be run once per virtual server.

Options:

- `--php-version [argument]` install a specific version of PHP. Options are `7.4`, `7.3`, and `7.2`.
- `--database [argument]` install a database engine. Options are `mysql` or `postgres`.

This boostraps a machine with the custom name `diesel`, using PHP 7.2 and PostgreSQL:

```bash
nitro --machine diesel bootstrap --php-version 7.2 --database postgres
```

### `add`

Adds a new virtual host to nginx and mounts a local directory to the server.

> Note: if you pass a version of PHP that was _not_ used when provisioning the server, Nitro will install that version of PHP for you.

This adds a host using `mysite.test` to the `diesel` machine, using PHP 7.4 and a document root of `/Users/jason/Sites/craftcms`.

```bash
nitro --machine diesel add --php-version 7.4 mysite.test /Users/jason/Sites/craftcms
```

### `remove`

Removes the specified virtual server from nginx and detaches the attached directory from the virtual server.

This removes the `mysite.test` host from the `diesel` machine:

```bash
nitro --machine diesel remove mysite.test
```

### `attach`

Attaches, or mounts, a local directory to an nginx server’s web root.

This mounts the local directory `/Users/jason/Sites/craftcms` as the web root for the `diesel` machine’s `mysite.test` host:

```bash
nitro --machine diesel attach mysite.test /Users/jason/Sites/craftcms
```

### `ssh`

Nitro gives you full root access to your virtual server. The default user is `ubuntu` and has `sudo` permissions without a password. Once you’re in the virtual server, you can run `sudo` commands as usual (e.g. `sudo apt install golang`).

This launches a new shell within the `diesel` machine:

```bash
nitro --machine diesel ssh
```

### `xon`

Enables Xdebug, which is installed and disabled by default on each machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.3 and enables it for the `diesel` machine:

```bash
nitro --machine diesel xon --php-version 7.3
```

### `xoff`

Disables Xdebug on a machine.

Options:

- `--php-version [argument]` install a specific version of PHP to enable for Xdebug

This ensures Xdebug is installed for PHP 7.2 but disables it for the `diesel` machine:

```bash
nitro --machine diesel xoff --php-version 7.2
```

### `start`

Starts, or turns on, a machine.

This turns on the `diesel` machine:

```bash
nitro --machine diesel start
```

### `stop`

Stops, or turns off, a machine.

This turns off the `diesel` machine:

```bash
nitro --machine diesel stop
```

### `destroy`

Destroys a machine. By default, Multipass does not permanently delete a machine and can cause name conflicts (e.g. `instance "nitro-dev" already exists`). This will not affect any files or directories attached to the machine.

Options:

- `--permanent` permanently deletes a machine **(this is non-recoverable!)**

This soft-destroys the `diesel` machine:

```bash
nitro --machine diesel destroy --permanent
```

This **permanently** destroys the `diesel` machine:

```bash
nitro --machine diesel destroy --permanent true
```

### `sql`

Launches a database shell as the root user.

- `--postgres` access the PostgreSQL shell rather than MySQL (default)

This launches a PostgreSQL console shell for the `diesel` machine:

```bash
nitro --machine diesel sql --postgres
```

### `redis`

Access a Redis shell.

This launches a Redis console shell for the `diesel` machine:

```bash
nitro --machine diesel redis
```

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

This upgrades the `diesel` machine’s software packages to their newest versions:

```bash
nitro --machine diesel update
```

### `logs`

Views the virtual machines logs.

Options:

- `nginx`

This displays nginx logs for the `diesel` machine:

```bash
nitro --machine diesel logs nginx
```

### `ip`

Prints the machine’s locally-accessible IP address.

This prints the IP address of the `diesel` machine:

```bash
nitro --machine diesel ip
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
nitro --machine diesel services restart --nginx --mysql
```

This stops PostgreSQL on the `diesel` machine:

```bash
nitro --machine diesel services stop --postgres
```

This starts Redis on the `diesel` machine:

```bash
nitro --machine diesel services start --redis
```

### `refresh`

Updates the scripts used to create virtual servers for nginx and other utilities. This is only needed when updating the Nitro CLI.

This updates the `diesel` machine to use Nitro’s most current virtual server scripts for the CLI:

```bash
nitro --machine diesel refresh
```
