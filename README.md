<p align="center"><img src="resources/craft-nitro.svg" width="100" height="100" alt="Craft Nitro icon"></p>

<h1 align="center">Craft Nitro</h1>

Nitro is a speedy local development environment that’s tuned for [Craft CMS](https://craftcms.com/), powered by [Multipass](https://multipass.run/).

- [What’s Included](#whats-included)
- [Installation](#installation)
  - [Uninstalling Nitro](#uninstalling-nitro)
- [Adding Sites](#adding-sites)
  - [Adding a site with `nitro add`](#add-a-site-with-nitro-add)
  - [Mounting your entire dev folder at once](#mounting-your-entire-dev-folder-at-once)
- [Connecting to the Database](#connecting-to-the-database)
- [Adding Mounts](#adding-mounts)
- [Running Multiple Machines](#running-multiple-machines)
- [Adding Multiple Database Engines](#adding-multiple-database-engines)
- [Using Blackfire](#using-blackfire)
- [Using Xdebug](#using-xdebug)
- [Commands](#commands)
  - [`apply`](#apply)
  - [`add`](#add)
  - [`context`](#context)
  - [`destroy`](#destroy)
  - [`edit`](#edit)
  - [`info`](#info)
  - [`init`](#init)
  - [`import`](#import)
  - [`logs`](#logs)
  - [`remove`](#remove)
  - [`redis`](#redis)
  - [`start`](#start)
  - [`stop`](#stop)
  - [`rename`](#rename)
  - [`restart`](#restart)
  - [`self-update`](#self-update)
  - [`ssh`](#ssh)
  - [`update`](#update)
  - [`version`](#version)
  - [`xdebug configure`](#xdebug-configure)
  - [`xdebug on`](#xdebug-on)
  - [`xdebug off`](#xdebug-off)

---

## What’s Included

Nitro installs the following on every machine:

- PHP 7.4 (+ option to use 7.3 or 7.2)
- MySQL
- PostgreSQL
- Redis
- Xdebug
- Blackfire

## Installation

> ⚠️ **Note:** Windows support is a [work-in-progress](https://github.com/craftcms/nitro/issues/88).

1. Install [Multipass](https://multipass.run) (requires 1.2.0+).
2. Run this terminal command:

    ```sh
    bash <(curl -sLS http://installer.getnitro.sh)
    ```

3. Follow the prompts to create your machine.

Once complete, you will have a Multipass machine called `nitro-dev`, and a new configuration file for the machine
 stored at `~/.nitro/nitro-dev.yaml`.

### Uninstalling Nitro

To completely remove Nitro, first [destroy](#destroy) your machine:

```sh
nitro destroy
```

> 💡 **Tip:** If you have multiple machines, you can destroy them all using the `multipass delete` command:
>
> ```sh
> multipass delete --all -p
> ```

Then remove your `nitro` command:

```sh
sudo rm /usr/local/bin/nitro
```

You can optionally remove your machine configs as well:

```sh
rm -rf ~/.nitro
```

If you wish to uninstall Multipass as well, uninstall instructions can be found on the installation guide for your platform:

- [macOS](https://multipass.run/docs/installing-on-macos)
- [Windows](https://multipass.run/docs/installing-on-windows)
- [Linux](https://multipass.run/docs/installing-on-linux)

## Adding Sites

To add a site to Nitro, three things need to happen:

- Your project files need to be [mounted](#adding-mounts) into the Nitro machine.
- The web server within your Nitro machine needs to be configured to serve your site.
- Your system’s `hosts` file needs to be updated to associate your site’s hostname with Nitro.

### Add a site with `nitro add`

If your project files are completely contained within a single folder, then you can quickly accomplish these using
the [`add`](#add) command:

```sh
$ cd /path/to/project
$ nitro add
What should the hostname be? [plugins-dev] example.test
Where is the webroot? [web]
plugins-dev has been added to config file.
Apply changes from config? [yes]
Applied changes from /Users/jasonmccallister/.nitro/nitro-dev.yaml
Editing your hosts file
Password: ******
example.test added successfully!
```

### Mounting your entire dev folder at once

If you manage all of your projects within a single dev folder, you can mount that entire folder once within Nitro,
and point your sites’ webroots to the appropriate folders within it.

To do that, open your `~/.nitro/nitro-dev.yaml` file in a text editor (or run the [`edit`](#edit) command), and add
a new mount for the folder that contains all of your projects, plus list out all of your sites you wish to add to
Nitro within that folder:

```yaml
mounts:
 - source: ~/dev
   dest: /nitro/sites
sites:
 - hostname: example1.test
   webroot: /nitro/sites/example1.test/web
 - hostname: example2.test
   webroot: /nitro/sites/example2.test/web
```

Then run `nitro apply` to apply your `nitro.yaml` changes to the machine.

> 💡 **Tip:** To avoid permission issues, we recommend you always mount folders into `/nitro/*` within the
  machine.

> ⚠️ **Warning:** If your projects contain any symlnks, such as `path` Composer repositories, those symlinks
  **must** be relative (`../`), not absolute (`/` or `~/`).

## Connecting to the Database

To connect to the machine from a Craft install, set the following environment variables in your `.env` file:

```
DB_USER="nitro"
DB_PASSWORD="nitro"
```

To connect to the database from your host operating system, you’ll first need to get the IP address of your Nitro machine. You can find that by running the [info](#info) command.

```sh
$ nitro info
Name:           nitro-dev
State:          Running
IPv4:           192.168.64.2
Release:        Ubuntu 18.04.4 LTS
Image hash:     2f6bc5e7d9ac (Ubuntu 18.04 LTS)
Load:           0.71 0.74 0.60
Disk usage:     2.7G out of 38.6G
Memory usage:   526.4M out of 3.9G
```

Then from your SQL client of choice, create a new database connection with the following settings:

- **Host**: _The `IPv4` value from `nitro info`_
- **Port**: _The port you configured your database with (3306 for MySQL or 5432 for PostgreSQL by default)._
- **Username**: `nitro`
- **Password**: `nitro`

## Adding Mounts

Nitro can mount various system directories into your Nitro machine. You can either mount each of your projects’
root directories into Nitro individually (as you’d get when [adding a site with `nitro
add`](#add-a-site-with-nitro-add)), or you can mount your entire development folder, or some combination of the two.

To add a new mount, follow these steps:

1. Open your `~/.nitro/nitro.yaml` file in a text editor, and add the new mount:

   ```yaml
   mounts:
     - source: /Users/cathy/dev
       dest: /nitro/sites
   ```

2. Run `nitro apply` to apply the `nitro.yaml` change to the machine.

Once that’s done, yous should be able to tunnel into your machine using the [`ssh`](#ssh) command, and see the
newly-mounted directory in there.

## Running Multiple Machines

You can have Nitro manage more than just your primary machine (`nitro-dev`) if you want. For example, you could
create a machine for a specific dev project.

To create a new machine, run the following command:

```sh
$ nitro init -m <machine>
```

Replace `<machine>` with the name you want to give your new machine. Machine names can only include letters,
numbers, underscores, and hyphen.

This command will run through the same prompts you saw when creating your primary machine after you first installed
Nitro. Once it’s done, you’ll have a new Multipass machine, as well as a new configuration file for it at
`~/.nitro/<machine>.yaml`.

All of Nitro’s [commands](#commands) accept an `-m` option, which you can use to specify which machine the command
should be run against. (`nitro-dev` will always be used by default.)

## Adding Multiple Database Engines

To run multiple database engines on the same machine, open your `~/.nitro/nitro-dev.yaml` file in a text editor (or
run the [`edit`](#edit) command), and list additional databases under the `databases` key:

```yaml
databases:
 - engine: mysql
   version: "5.7"
   port: "3306"
 - engine: mysql
   version: "5.6"
   port: "33061"
 - engine: postgres
   version: "11"
   port: "5432"
```

> ⚠️ **Warning:** Each database engine needs its own unique port.

Then run `nitro apply` to apply your `nitro.yaml` changes to the machine.

## Using Blackfire

See [Using Blackfire with Nitro](BLACKFIRE.md) for instructions on how to configure and run Blackfire.

## Using Xdebug

See [Using Xdebug with Nitro and PhpStorm](XDEBUG.md) for instructions on how to configure Xdebug and PhpStorm for web/console debugging.

## Commands

The following commands will help you manage your virtual server.

- [`apply`](#apply)
- [`add`](#add)
- [`context`](#context)
- [`destroy`](#destroy)
- [`edit`](#edit)
- [`info`](#info)
- [`init`](#init)
- [`import`](#import)
- [`logs`](#logs)
- [`remove`](#remove)
- [`redis`](#redis)
- [`rename`](#rename)
- [`restart`](#restart)
- [`self-update`](#self-update)
- [`ssh`](#ssh)
- [`start`](#start)
- [`stop`](#stop)
- [`update`](#update)
- [`version`](#version)
- [`xdebug configure`](#xdebug-configure)
- [`xdebug on`](#xdebug-on)
- [`xdebug off`](#xdebug-off)

### `apply`

Ensures that the machine exists, and applies any changes in its config file to it.

```sh
nitro apply [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```sh
$ nitro apply
There are 2 mounted directories and 1 new mount(s) in the config file.
Applied changes from nitro.yaml.
```

### `add`

Adds a new site to the machine.

```sh
nitro add [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--hostname</code></dt>
<dd>The hostname to use for accessing the site. If not passed, the command will prompt for it.</dd>
<dt><code>--webroot</code></dt>
<dd>The relative path to the site’s webroot. If not passed, the command will prompt for it.</dd>
</dl>

Example:

```sh
$ cd /path/to/project
$ nitro add
What should the hostname be? [plugins-dev]
Where is the webroot? [web]
plugins-dev has been added to config file.
Apply changes from config? [yes]
Applied changes from /Users/jasonmccallister/.nitro/nitro-dev.yaml
Editing your hosts file
Password: ******
plugins-dev added successfully!
```

### `context`

Shows the machine’s configuration.

```sh
nitro context [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```sh
$ nitro context
Machine: nitro-dev
------
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

### `destroy`

Destroys a machine.

```sh
nitro destroy [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--clean</code></dt>
<dd>Remove the configuration file after destroying the machine. Defaults to `false`</dd>
</dl>

### `edit`

Edit allows you to quickly open your machine configuration to make changes.

```sh
nitro edit [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```sh
nitro edit
```

> 💡 **Tip:** If you’re running macOS or Linux, you can set an `EDITOR` environment variable in `~/.bash_profile` to your preferred text editor of choice.
>
> ```sh
> export EDITOR="/Applications/Sublime Text.app/Contents/MacOS/Sublime Text"
> ```
>
> After adding that line, restart your terminal or run `source ~/.bash_profile` for the change to take effect.
>
> Alternatively, you can open the configuration file using your operating system’s default text editor for `.yaml` files by running this command:
>
> ```sh
> open ~/.nitro/nitro-dev.yaml
> ```
> 
> (Replace `nitro-dev` with the appropriate machine name if it’s not that.)

### `info`

Shows the _running_ information for a machine like the IP address, memory, disk usage, and mounts.

```sh
nitro info [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```sh
$ nitro info
Name:           nitro-dev
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

### `init`

Initializes a machine.

```sh
nitro init [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to use. If not passed, the command will prompt for it.</dd>
<dt><code>--cpus</code></dt>
<dd>The max number of CPUs that the machine can use. If not passed, one CPU will be used by default.</dd>
<dt><code>--memory</code></dt>
<dd>The max amount of system RAM that the machine can use. If not passed, the command will prompt for it.</dd>
<dt><code>--disk</code></dt>
<dd>The max amount of disk space that the machine can use. If not passed, the command will prompt for it.</dd>
</dl>

If the machine already exists, it will be reconfigured.

### `import`

Import a SQL file into a database in the machine. You will be prompted with a list of running database engines
(MySQL and PostgreSQL) to import the file into.

```sh
nitro import <file> [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```sh
$ nitro import mybackup.sql
Use the arrow keys to navigate: ↓ ↑ → ←
? Select database:
  ▸ mysql_5.7_3306
```

### `logs`

Views the machine’s logs. This command will prompt you for a type of logs to view, including e.g. `nginx`,
`database`, or `docker` (for a specific container).

```sh
nitro logs [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `remove`

Removes a site from the machine.

```sh
nitro remove [<options>]
```

You will be prompted to select the site that should be removed. If the site has a corresponding
[mount](#adding-mounts) at `/nitro/sites/<hostname>`, the mount will be removed as well.

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `redis`

Starts a Redis shell.

```sh
nitro redis [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `start`

Starts the machine.

```sh
nitro start [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `stop`

Stops the machine.

```sh
nitro stop [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `rename`

Rename a site in a configuration file. Will prompt for which site to rename.

```sh
nitro rename [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `restart`

Restarts a machine.

```sh
nitro restart [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `self-update`

Perform updates to the Nitro CLI.

```sh
nitro self-update
```

### `ssh`

Tunnels into the machine as the default `ubuntu` user over SSH.

```sh
nitro ssh [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

```sh
nitro update [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `version`

Checks the currently version of nitro against the releases and shows any updated versions.

```sh
nitro version
```

### `xdebug configure`

Configures Xdebug for remote access and debugging with PhpStorm or other IDEs.

```sh
nitro xdebug configure [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to configure Xdebug for</dd>
</dl>

### `xdebug on`

Enables Xdebug, which is installed and disabled by default on each machine.

```sh
nitro xdebug on [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to enable Xdebug for</dd>
</dl>

This ensures Xdebug is installed for PHP 7.3 and enables it:

### `xdebug off`

Disables Xdebug on a machine.

```sh
nitro xdebug off [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to disable Xdebug for</dd>
</dl>
