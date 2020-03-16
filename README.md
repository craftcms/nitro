# Craft Nitro

A better and faster way to develop Craft CMS applications locally without Docker or Vagrant! Nitro has one dependency, Multipass, which will allow create virtual machines based on Ubuntu.

# Requirements

- [Multipass](https://multipass.run)

## What's Included

Nitro will install the following default on the server, the version of PHP and database engine are configurable.

- PHP 7.4
- MariaDB
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

### init

TODO

### bootstrap

TODO

### add

TODO

### remove

TODO

## attach

TODO

### ssh

TODO

### xon

TODO

### xoff

TODO

### start

TODO

### stop

TODO

### destroy

TODO

### sql

TODO

### redis

TODO

### update

TODO

### logs

TODO

### ip

TODO

### services

TODO

### refresh

TODO
