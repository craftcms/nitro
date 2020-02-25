# Craft Nitro

A better and faster way to develop Craft CMS applications locally without Docker or Vagrant!

## Installation

```shell script
composer require --dev craftcms/nitro
```

## Usage

This package has a single executable named `nitro`. In order to 
use the CLI, run `./vendor/bin/nitro`. 

Creating a new machine for development is simple, run the following command:

```shell script
./vendor/bin/nitro pour
```

This will create a new machine and give it the default name of `nitro`. If you wish to assign another name to the machine, run the following command instead:

```shell script
./vendor/bin/nitro --machine my-custom-name init
``` 
