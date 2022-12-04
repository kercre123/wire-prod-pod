#!/bin/bash

if [[ ! -d ./chipper ]]; then
    echo "This must be run in the wire-prod-pod/ directory."
    exit 0
fi

if [[ $1 == "-d" ]]; then
    sleep 20
    while true; do
        sudo git pull --force > /tmp/gitTest 2> /tmp/gitTest
        gitTestOut=$(cat /tmp/gitTest)
        if [[ ${gitTestOut} != *"up to date"* ]]; then
            git fetch --all
            git reset --hard origin/main
            systemctl stop wire-pod
            cd chipper
            sudo ./build.sh
            cd ..
            echo
            systemctl start wire-pod
            echo "Updated!"
            rm -f /tmp/gitTest
        else
            echo "No update needed."
        fi
        echo
        sleep 86400
    done
else
    systemctl stop wire-pod
    git fetch --all
    git reset --hard origin/main
    cd chipper
    sudo ./build.sh
    cd ..
    echo
    systemctl start wire-pod
    echo "Updated!"
fi
