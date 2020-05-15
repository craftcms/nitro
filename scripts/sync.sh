#!/usr/bin/env bash
version="$!"

if [ -z $version ]; then
  echo "you must specify a version of nitro"
  exit 1
fi

if [ $version == "1.0.0-beta.3" ]; then
  echo "running scripts for 1.0.0-beta.3"

  # create setup scripts
  mkdir -p /home/ubuntu/sites
  mkdir -p /home/ubuntu/.nitro/databases/imports
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/mysql/backups
  mkdir -p /home/ubuntu/.nitro/databases/postgres/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/mysql/conf.d
  mkdir -p /home/ubuntu/.nitro/databases/postgres/backups
  chown -R ubuntu:ubuntu /home/ubuntu/.nitro
  chown -R ubuntu:ubuntu /home/ubuntu/sites

  echo "setup script has completed!"
  exit 0
fi
