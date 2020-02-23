# Dev

A better way to develop PHP application locally without Docker or Vagrant.

## Installation

```shell script
composer require --dev craftcms\dev
```

## Usage

This package has a single executable named `dev`. In order to 
use the CLI, run `./vendor/bin/dev`. 

### Starting a new Machine

Creating a new machine is simple, run the following command:

```shell script
./vendor/bin/dev init
```

This will create a new machine and give it a default name of `dev`. If you wish to assign another name to the machine, run the following command:

```shell script
./vendor/bin/dev --machine my-custom-name init
``` 

### Install PHP

After you have created a machine, you can install PHP on that machine with the following command:

```shell script
./vendor/bin/dev install php --version 7.4
```

> Note: the --version flag is optional and will default to the latest PHP version.

### Install Nginx

```shell script
./vendor/bin/dev install nginx
```

### Install MariaDB

```shell script
./vendor/bin/dev install mariadb
```
