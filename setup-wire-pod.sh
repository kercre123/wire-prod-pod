#!/bin/bash

if [[ -f /usr/bin/apt ]]; then
        TARGET="debian"
        echo "Debian-based Linux confirmed."
elif [[ -f /usr/bin/pacman ]]; then
        TARGET="arch"
        echo "Arch Linux confirmed."
elif [[ -f /usr/bin/dnf ]]; then
        TARGET="fedora"
        echo "Fedora/openSUSE detected."
else
        echo "This OS is not supported. This script currently supports Linux with either apt, pacman, or dnf."
        exit 1
fi

if [[ ${TARGET} == "debian" ]]; then
                apt update -y
                apt install -y git
elif [[ ${TARGET} == "arch" ]]; then
                pacman -Sy --noconfirm
                sudo pacman -S --noconfirm git
elif [[ ${TARGET} == "fedora" ]]; then
                dnf update
                dnf install -y git
fi

cd ~

git clone https://github.com/kercre123/wire-prod-pod
cd wire-prod-pod
sudo ./setup.sh
