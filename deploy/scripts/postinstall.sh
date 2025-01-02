#!/bin/sh
if ! getent group byteload >/dev/null; then
  groupadd -r byteload
fi
if ! getent passwd byteload >/dev/null; then
  useradd -r -g byteload -s /sbin/nologin -d /var/lib/byteload byteload
fi
mkdir -p /var/lib/byteload
chown -R byteload:byteload /var/lib/byteload /etc/byteload
systemctl daemon-reload
systemctl enable byteload.service
systemctl start byteload.service || true 