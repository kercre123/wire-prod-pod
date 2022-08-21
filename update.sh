#!/bin/bash

if [[ ! -d ./chipper ]]; then
  echo "This must be run in the jank-escape-pod/ directory."
  exit 0
fi

if [[ $1 == "-d" ]]; then
while true; do
systemctl stop wire-pod
git pull
cd chipper
sudo ./build.sh
cd ..
echo
systemctl start wire-pod
echo "Updated!"
echo
sleep 600
done
else
systemctl stop wire-pod
git pull
cd chipper
sudo ./build.sh
cd ..
echo
systemctl start wire-pod
echo "Updated!"
fi
