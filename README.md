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
- [Adding a Database](#adding-a-database)
- [Adding Mounts](#adding-mounts)
- [Running Multiple Machines](#running-multiple-machines)
- [Adding Multiple Database Engines](#adding-multiple-database-engines)
- [Using Blackfire](#using-blackfire)
- [Using Xdebug](#using-xdebug)
- [Using MailHog](#using-mailhog)
- [Advanced Configuration](#advanced-configuration)
- [Commands](#commands)
  - [`apply`](#apply)
  - [`add`](#add)
  - [`context`](#context)
  - [`db add`](#db-add)
  - [`db backup`](#db-backup)
  - [`db import`](#db-import)
  - [`db remove`](#db-remove)
  - [`db restart`](#db-restart)
  - [`db start`](#db-start)
  - [`db stop`](#db-stop)
  - [`destroy`](#destroy)
  - [`edit`](#edit)
  - [`info`](#info)
  - [`init`](#init)
  - [`install composer`](#install-composer)
  - [`install mysql`](#install-mysql)
  - [`install postgres`](#install-postgres)
  - [`keys`](#keys)
  - [`logs`](#logs)
  - [`remove`](#remove)
  - [`redis`](#redis)
  - [`rename`](#rename)
  - [`restart`](#restart)
  - [`self-update`](#self-update)
  - [`start`](#start)
  - [`stop`](#stop)
  - [`ssh`](#ssh)
  - [`update`](#update)
  - [`version`](#version)
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

----------

macOS and Linux:

1. Install [Multipass](https://multipass.run) (requires 1.2.0+).
2. Run this terminal command:

    ```sh
    bash <(curl -sLS http://installer.getnitro.sh)
    ```

3. Follow the prompts to create your machine.

----------

Windows 10 Pro (with Hyper-V enabled):

> 💡 **Tip:** Windows doesn't currently have an automated install script, so installation and updating must be done manually.

1. Install [Multipass](https://multipass.run) (requires 1.2.0+).
2. Download `nitro_windows_x86_64.zip` from the latest [release](https://github.com/craftcms/nitro/releases)
3. Create a `Nitro` folder in your home folder, if it does not exist. i.e. `C:\Users\<username>\Nitro`
4. Extract the zip file and copy `nitro.exe` into the `Nitro` folder you just created in your home folder.
5. If this is your first installation, run this from the command line to add `nitro` to your global path: `setx path "%PATH%;%USERPROFILE%\Nitro"`
6. Start the Windows terminal (cmd.exe) with Administrator permissions and run `nitro init` to create your first machine.

----------

Once complete, you will have a Multipass machine called `nitro-dev`, and a new configuration file for the machine
 stored at `~/.nitro/nitro-dev.yaml`.

### Uninstalling Nitro

To completely remove Nitro, first [destroy](#destroy) your machine:

```shell script
nitro destroy
```

> 💡 **Tip:** If you have multiple machines, you can destroy them all using the `multipass delete` command:
>
> ```sh
> multipass delete --all -p
> ```

Then remove your `nitro` command:

macOS and Linux:

```shell script
sudo rm /usr/local/bin/nitro
```

Windows:

```shell script
rm -rf $HOME/Nitro
```

You can optionally remove your machine configs as well:

macOS and Linux

```shell script
rm -rf ~/.nitro
```

Windows:

```shell script
rm -rf $HOME/.nitro
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

```shell script
$ cd /path/to/project
$ nitro add
Enter the hostname [plugins-dev] example.test
Enter the webroot [web]
Added plugins-dev to config file.
Apply changes from config? [yes]
Mounting /path/to/project to nitro-dev
Adding site example.test to nitro-dev
Applied changes from /Users/vin/.nitro/nitro-dev.yaml
Editing your hosts file
Password: ******
```

> 💡 **Tip:** Multipass requires Full Disk Access on macOS. If you’re seeing mount “not readable” issues, ensure `multipassd` is checked under System Preferences → Security & Privacy → Privacy → Full Disk Access.

### Mounting your entire dev folder at once

If you manage all of your projects within a single dev folder, you can mount that entire folder once within Nitro,
and point your sites’ webroots to the appropriate folders within it.

To do that, open your `~/.nitro/nitro-dev.yaml` file in a text editor (or run the [`edit`](#edit) command), and add
a new mount for the folder that contains all of your projects, plus list out all of your sites you wish to add to
Nitro within that folder:

```yaml
mounts:
 - source: ~/dev
   dest: /home/ubuntu/sites
sites:
 - hostname: example1.test
   webroot: /home/ubuntu/sites/example1.test/web
 - hostname: example2.test
   webroot: /home/ubuntu/sites/example2.test/web
```

Then run `nitro apply` to apply your `nitro.yaml` changes to the machine.

> 💡 **Tip:** To avoid permission issues, we recommend you always mount folders into `/home/ubuntu/*` within the
  machine.

> ⚠️ **Warning:** If your projects contain any symlinks, such as `path` Composer repositories, those symlinks
  **must** be relative (`../`), not absolute (`/` or `~/`).

## Connecting to the Database

To connect to the machine from a Craft install, set the following environment variables in your `.env` file:

```
DB_USER="nitro"
DB_PASSWORD="nitro"
```

To connect to the database from your host operating system, you’ll first need to get the IP address of your Nitro machine. You can find that by running the [info](#info) command.

```shell script
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

## Adding a Database

Nitro creates its initial database for you. You can add as many as you’d like running the following command, which will prompt for your desired database engine and name:

```shell script
$ nitro db add
```

## Adding Mounts

Nitro can mount various system directories into your Nitro machine. You can either mount each of your projects’
root directories into Nitro individually (as you’d get when [adding a site with `nitro
add`](#add-a-site-with-nitro-add)), or you can mount your entire development folder, or some combination of the two.

To add a new mount, follow these steps:

1. Open your `~/.nitro/nitro.yaml` file in a text editor, and add the new mount:

   ```yaml
   mounts:
     - source: /Users/vin/dev
       dest: /home/ubuntu/sites
   ```

2. Run `nitro apply` to apply the `nitro.yaml` change to the machine.

Once that’s done, you should be able to tunnel into your machine using the [`ssh`](#ssh) command, and see the
newly-mounted directory in there.

## Running Multiple Machines

You can have Nitro manage more than just your primary machine (`nitro-dev`) if you want. For example, you could
create a machine for a specific dev project.

To create a new machine, run the following command:

```shell script
$ nitro init -m <machine>
```

Replace `<machine>` with the name you want to give your new machine. Machine names can only include letters,
numbers, underscores, and hyphens.

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

## Using MailHog

See [Using MailHog with Craft](MAILHOG.md) for instructions on configuring Craft to send email to MailHog for local development and troubleshooting.

## Advanced Configuration

See [Advanced Configuration](ADVANCED.md) for instructions on customizing Nitro’s default settings.

## Commands

### `apply`

Ensures that the machine exists, and applies any changes in its config file to it.

```shell script
nitro apply [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--skip-hosts</code></dt>
<dd>Skips updating the <code>hosts</code> file.</dd>
</dl>

Example:

```shell script
$ nitro apply
There are 2 mounted directories and 1 new mount(s) in the config file.
Applied changes from nitro.yaml.
```

### `add`

Adds a new site to the machine.

```shell script
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
<dt><code>--skip-hosts</code></dt>
<dd>Skips updating the <code>hosts</code> file.</dd>
</dl>

Example:

```shell script
$ cd /path/to/project
$ nitro add
Enter the hostname [plugins-dev] example.test
Enter the webroot [web]
Added plugins-dev to config file.
Apply changes from config? [yes]
Mounting /path/to/project to nitro-dev
Adding site example.test to nitro-dev
Applied changes from /Users/vin/.nitro/nitro-dev.yaml
Editing your hosts file
Password: ******
```

### `context`

Shows the machine’s configuration.

```shell script
nitro context [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro context
Machine: nitro-dev
------
php: "7.4"
cpus: "1"
disk: 40G
memory: 4G
mounts:
- source: ~/sites/demo-site
  dest: /home/ubuntu/sites/demo-site
databases:
- engine: mysql
  version: "5.7"
  port: "3306"
- engine: postgres
  version: "12"
  port: "5432"
sites:
- hostname: demo-site
  webroot: /home/ubuntu/sites/demo-site/web
------
```

### `db add`

Create a new database on a database engine in a machine.

```shell script
nitro db add [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db add
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine [1] 2
Enter the name of the database: mynewproject
Added database "mynewproject" to "mysql_5.7_3306".
```

### `db backup`

Backup one or all databases from a database engine in a machine. 

```shell script
nitro db backup [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db backup
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine [1] 
  1 - all-dbs
  2 - postgres
  3 - nitro
  4 - project-one
Select database to backup? [1] 
Created backup "all-dbs-200519_100730.sql", downloading...
Backup completed and stored in "/Users/vin/.nitro/backups/nitro-dev/postgres_11_5432/all-dbs-200519_100730.sql"
```

### `db import`

Import a SQL file into a database engine in a machine. You will be prompted with a list of running database engines
(MySQL and PostgreSQL) to import the file into.

```shell script
nitro db import <file> [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db import mybackup.sql
  1 - mysql_5.7_3306
  2 - postgres_11_5432
Select database engine [1] 
Enter the database name to create for the import: new-project
Uploading "mybackup.sql" into "nitro-dev" (large files may take a while)...
Successfully import the database backup into new-project
```

### `db remove`

Will remove a database from a database engine in a machine, but not from the config file.

```shell script
nitro db remove [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db remove
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine: [1] 
  1 - nitro
  2 - project-one
 
Are you sure you want to permanently remove the database nitro? [no] 
Removed database nitro
```

### `db restart`

Will restart a database engine in a machine.

```shell script
nitro db restart [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db restart
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine to restart: [1]  
Restarted database engine postgres_11_5432
```

### `db start`

Will start a stopped database engine in a machine.

```shell script
nitro db start [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db start
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine to start: [1]  
Started database engine postgres_11_5432
```

### `db stop`

Will stop a database engine in a machine.

```shell script
nitro db stop [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro db stop
  1 - postgres_11_5432
  2 - mysql_5.7_3306
Select database engine to stop: [1]
Stopped database engine postgres_11_5432
```

### `destroy`

Destroys a machine.

```shell script
nitro destroy [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--clean</code></dt>
<dd>Remove the configuration file after destroying the machine. Defaults to `false`</dd>
<dt><code>--skip-hosts</code></dt>
<dd>Skips updating the <code>hosts</code> file.</dd>
</dl>

### `edit`

Edit allows you to quickly open your machine configuration to make changes.

```shell script
nitro edit [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
nitro edit
```

> 💡 **Tip:** If you’re running macOS or Linux, you can set an `EDITOR` environment variable in `~/.bash_profile` to your preferred text editor of choice.
>
> ```shell script
> export EDITOR="/Applications/Sublime Text.app/Contents/MacOS/Sublime Text"
> ```
>
> After adding that line, restart your terminal or run `source ~/.bash_profile` for the change to take effect.
>
> Alternatively, you can open the configuration file using your operating system’s default text editor for `.yaml` files by running this command:
>
> ```shell script
> open ~/.nitro/nitro-dev.yaml
> ```
> 
> (Replace `nitro-dev` with the appropriate machine name if it’s not that.)

### `info`

Shows the _running_ information for a machine like the IP address, memory, disk usage, and mounts.

```shell script
nitro info [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro info
Name:           nitro-dev
State:          Running
IPv4:           192.168.64.48
Release:        Ubuntu 18.04.4 LTS
Image hash:     2f6bc5e7d9ac (Ubuntu 18.04 LTS)
Load:           0.09 0.15 0.22
Disk usage:     2.7G out of 38.6G
Memory usage:   379.8M out of 3.9G
Mounts:         /Users/vin/sites/demo-site => /home/ubuntu/sites/demo-site
                    UID map: 501:default
                    GID map: 20:default
```

### `init`

Initializes a machine.

```shell script
nitro init [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to use. If not passed, the command will prompt for it.</dd>
<dt><code>--cpus</code></dt>
<dd>The max number of CPUs that the machine can use. If not passed, Nitro will try to determine the best number based on the host computer.</dd>
<dt><code>--memory</code></dt>
<dd>The max amount of system RAM that the machine can use. If not passed, the command will prompt for it.</dd>
<dt><code>--disk</code></dt>
<dd>The max amount of disk space that the machine can use. If not passed, the command will prompt for it.</dd>
</dl>

If the machine already exists, it will be reconfigured.

### `install composer`

Install composer inside of a virtual machine.

```shell script
nitro install composer
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro install composer
Composer is now installed on "nitro-dev".
```

### `install mysql`

This will add a new MySQL database engine. 

```shell script
nitro install mysql
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro install mysql
Enter the MySQL version to install: 5.6
Enter the MySQL port number: 3306
Adding MySQL version "5.6" on port "3306"
Apply changes from config now? [yes]
```

### `install postgres`

This will add a new PostgreSQL database engine. 

```shell script
nitro install postgres
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro install postgres
Enter the PostgreSQL version to install: 11
Enter the MySQL port number: 5432
Adding MySQL version "11" on port "5432"
Apply changes from config now? [yes]
```

### `keys`

Import SSH keys intro a virtual machine for use with Composer, git, etc.

```shell script
nitro keys [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

Example:

```shell script
$ nitro keys
  1 - id_rsa
  2 - personal_rsa
Select the key to add to "nitro-dev"? [1]
Transferred the key "id_rsa" into "nitro-dev".
```

### `logs`

Views the machine’s logs. This command will prompt you for a type of logs to view, including e.g. `nginx`,
`database`, or `docker` (for a specific container).

```shell script
nitro logs [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `redis`

Starts a Redis shell.

```shell script
nitro redis [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `remove`

Removes a site from the machine.

```shell script
nitro remove [<options>]
```

You will be prompted to select the site that should be removed. If the site has a corresponding
[mount](#adding-mounts) at `/home/ubuntu/sites/<hostname>`, the mount will be removed as well.

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `rename`

Rename a site in a configuration file. Will prompt for which site to rename.

```shell script
nitro rename [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `restart`

Restarts a machine.

```shell script
nitro restart [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `self-update`

Perform updates to the Nitro CLI.

```shell script
nitro self-update
```

> ⚠️ **Warning:** This command does not work on Windows. You will need to perform a [manual installation.](https://github.com/craftcms/nitro/blob/master/README.md#installation)

### `ssh`

Tunnels into the machine as the default `ubuntu` user over SSH.

```shell script
nitro ssh [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `start`

Starts the machine.

```shell script
nitro start [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `stop`

Stops the machine.

```shell script
nitro stop [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `update`

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

```shell script
nitro update [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
</dl>

### `version`

Checks the currently version of nitro against the releases and shows any updated versions.

```shell script
nitro version
```

### `xdebug on`

Enables Xdebug, which is installed and disabled by default on each machine.

```shell script
nitro xdebug on [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to enable Xdebug for</dd>
</dl>

This ensures Xdebug is installed for PHP and enables it:

### `xdebug off`

Disables Xdebug on a machine.

```shell script
nitro xdebug off [<options>]
```

Options:

<dl>
<dt><code>-m</code>, <code>--machine</code></dt>
<dd>The name of the machine to use. Defaults to <code>nitro-dev</code>.</dd>
<dt><code>--php-version</code></dt>
<dd>The PHP version to disable Xdebug for</dd>
</dl>
