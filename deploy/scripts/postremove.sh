#!/bin/sh
if [ "$1" = "remove" ]; then
  systemctl disable byteload.service || true
  systemctl stop byteload.service || true
fi
systemctl daemon-reload 