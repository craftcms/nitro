#!/bin/bash

# download the
curl -s https://api.github.com/repos/craftcms/nitro/releases/latest \
  | grep "browser_download_url" \
  | grep "nitrod_linux_x86_64" \
  | cut -d : -f 2,3 | tr -d \" \
  | wget --directory-prefix=/tmp -qi -

# move nitrod into place
cd /tmp && tar xfz /tmp/nitrod_linux_x86_64.tar.gz
mv /tmp/nitrod /usr/sbin/
mv /tmp/nitrod.service /etc/systemd/system/

# setup the service
systemctl daemon-reload
systemctl start nitrod
systemctl enable nitrod
