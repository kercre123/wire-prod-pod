#!/bin/bash

sudo systemctl stop wire-pod wire-pod-updater
sudo systemctl disable wire-pod wire-pod-updater
sudo rm -rf /root/.coqui
sudo rm -f /lib/systemd/system/wire-pod.service /lib/systemd/system/wire-pod-updater.service
sudo systemctl daemon-reload
echo
echo "Uninstalled. Now you just need to remove this directory."
