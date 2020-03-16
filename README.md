# Craft Nitro

A better and faster way to develop Craft CMS applications locally without Docker or Vagrant! Nitro has one dependency, Multipass, which will allow create virtual machines based on Ubuntu.

## Requirements

- [Multipass](https://multipass.run)

## What's Included

Nitro will install the following default on the machine, both the version of PHP and database engine are configurable.

- PHP 7.4
- MySQL (mariadb flavor)
- PostgreSQL (optional)
- Redis

## Installation

```shell script
composer require --dev craftcms/nitro
```

Add the following script to your `composer.json` scripts section:

```
"scripts": {
    // other scripts
    "nitro": "./vendor/bin/nitro"
  },
```

This package has a single executable named `nitro`. In order to 
use the CLI, run `./vendor/bin/nitro`. 

Creating a new machine for development is simple, run the following command:

```shell script
composer run nitro init
```

This will create a new machine and give it the default name of `nitro-dev`. If you wish to assign another name to the machine, run the following command instead:

```shell script
./vendor/bin/nitro --machine my-custom-name init
``` 

## Usage

In order to create a new development server, you must "initialize" nitro. By default, this will not mount any directories and can be viewed as getting a brand new Virtual Private Server (VPS).

```bash
nitro init
``` 

> Note: `nitro init` has options that you can pass when creating a new server. However, we can set some "sane" defaults that should work for most scenarios. To view the options, run `nitro init -h`
 
After running `init` the system will default to automatically `bootstrap` the server. The bootstrap process will install the latest PHP version, MySQL, Redis.

> Note: if you wish to avoid bootstrapping, pass the --bootstrap false when running init (e.g. `nitro init --bootstrap false`)

The next step is to add a new virtual host to the server. To make this simple, you can run:

```bash
nitro add mysite.test path/to/site
```

This process will perform the following tasks:

1. Set up a new virtual server in Nginx for `mysite.test`
2. Attach the directory `./path/to/site` to your virtual server
3. Edit your /etc/hosts files to point `mysite.test` to your virtual server

You can now visit `http://mysite.test` in your browser!

## Commands

The following commands will help you manage your virtual server.

> Note: `nitro` will default to the virtual server name `nitro-dev`, these commands are assuming you are connecting to a virtual server named `mysite`. If you are using the default servername, you can skip adding the `--machine` argument. 

### init

This will create a new server. The following options are available when creating a new virtual server:

- `--boostrap (default: true)` will bootstrap the installation of PHP, MySQL, and Redis
- `--php-version (default 7.4)` will be passed to the bootstrap command to install a specific version of PHP
- `--database (default mysql)` passed to bootstrap command for installation, valid options are `mysql` or `postgres` 
- `--cpus (default 2)` the number of CPUs to allocate to the server
- `--memory (default 2G)` how much memory to allocate to the server
- `--disk (default 5G)` how much disk space to allocate

### bootstrap

Will install the specified version of PHP, the database engine, and Redis server onto a server. This should only be run once per virtual server. 

- `--php-version (default 7.4)` install a specific version of PHP
- `--database (default mysql)` install a database engine, valid options are `mysql` or `postgres`

#### Full Example

```bash
nitro --machine diesel bootstrap --php-version 7.2 --database postgres 
```

### add

Adds a new virtual host to nginx and attaches a directory to a server.

> Note: if you pass a version of PHP that was _not_ used when provisioning the server, nitro will install that version of PHP for you.

#### Full Example

```bash
nitro --machine diesel add --php-version 7.4 mysite.test ./path/to/site
```

### remove

this will remove the specified virtual server from nginx and then detach the directory from the virtual server.

#### Full Example

```bash
nitro --machine diesel remove mysite.test
```

### attach

To attach, or mount, a directory to a virtual server in nginx, use this command. The first argument is the path to the virtual server root followed by a path to the sites directory.
 
#### Full Example

```bash
nitro --machine diesel attach mysite.test ./path/to/mysite
```

### ssh

Nitro gives you full root access to your virtual server. By default, the user is `ubuntu`. This user has `sudo` permissions without a password. Once you are in the virtual server, you can run sudo commands as usual (e.g. `sudo apt install golang`)

#### Full Example

```bash
nitro --machine diesel ssh
```

### xon

xDebug is installed on each machine, but it not enabled by default. Quickly enable xDebug on a machine.

- `--php-version (default 7.4)` install a specific version of PHP to enable for xDebug

#### Full Example

```bash
nitro --machine diesel xon --php-version 7.3
``` 

### xoff

Quickly disable xDebug on a machine. The PHP version can be specified when 

- `--php-version (default 7.4)` install a specific version of PHP to enable for xDebug

#### Full Example

```bash
nitro --machine diesel xoff --php-version 7.2
``` 

### start

Start, or turn on, a machine.

#### Full Example

```bash
nitro --machine diesel start
``` 

### stop

Stop, or turn off, a machine.

#### Full Example

```bash
nitro --machine diesel stop
``` 

### destroy

Destroy a machine. My default, Multipass does not permanently delete a machine and can cause name conflicts (e.g. `instance "nitro-dev" already exists`).

- `--permanent (default false)` to permanently delete a machine (this is non-recoverable)

#### Full Example

```bash
nitro --machine diesel destroy --permanent
``` 

### sql

Access a SQL shell, as the root user, without using a GUI.

- `--postgres (default false)` to access the postgres SQL shell

#### Full Example

```bash
nitro --machine diesel sql --postgres
``` 

### redis

Access a Redis shell.

#### Full Example

```bash
nitro --machine diesel redis
``` 

### update

Performs system updates (e.g. `sudo apt get update && sudo apt upgrade -y`).

#### Full Example

```bash
nitro --machine diesel update
``` 

### logs

View the virtual machines logs.

#### Full Example

```bash
nitro --machine diesel logs nginx
```   

### ip

View the virtual machines IP address.

#### Full Example

```bash
nitro --machine diesel ip
```   

### services

Stop, start, or restart services on a virtual machine. The commands are nested under the `services` command so calling `nitro services` the sub commands will be listed.

#### Full Example

```bash
nitro --machine diesel services restart --nginx|--mysql|--postgres|--redis
```   

```bash
nitro --machine diesel services stop --nginx|--mysql|--postgres|--redis
```   

```bash
nitro --machine diesel services start --nginx|--mysql|--postgres|--redis
```   

### refresh

This command will update the scripts used to create virtual servers for nginx and other utilities. This is only needed when updating the nitro cli.

#### Full Example
  
```bash
nitro --machine diesel refresh
```   
