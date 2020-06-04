# Advanced Configuration

## `mysql.conf`

You can customize Nitro’s MySQL settings by editing `mysql.conf` and restarting the MySQL service:

1. Tunnel into the machine with `nitro ssh`.
2. Apply your changes to `~/.nitro/databases/mysql/conf.d/mysql.conf`.
3. Exit the machine tunnel by running `exit`.
4. Restart MySQL using `nitro db restart`.

## `NITRO_EDIT_HOSTS`

If you add a `NITRO_EDIT_HOSTS` environment variable to your system and set it to `false`, Nitro will never edit the `hosts` file on the host machine.

This is useful for people running some host file manager applications. 