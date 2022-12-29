#!/bin/bash

    botAddress=$1
    keyPath=$2
	if [[ ! -n ${botAddress} ]]; then
		echo "not enough args"
		exit 0
	fi
	if [[ ! -f ../certs/server_config.json ]]; then
		echo "server_config.json file missing. You need to generate this file with ./setup.sh's 6th option."
		exit 0
	fi
	if [[ ! -n ${keyPath} ]]; then
		echo
		if [[ ! -f ./ssh_root_key ]]; then
			echo "Key not provided, downloading ssh_root_key..."
			wget http://wire.my.to:81/ssh_root_key
		else
			echo "Key not provided, using ./ssh_root_key (already there)..."
		fi
		chmod 600 ./ssh_root_key
		keyPath="./ssh_root_key"
	fi
	if [[ ! -f ${keyPath} ]]; then
		echo "The key that was provided was not found. Exiting."
		exit 0
	fi
	ssh -oStrictHostKeyChecking=no -o ConnectTimeout=5  -i ${keyPath} root@${botAddress} "cat /build.prop" >/tmp/sshTest 2>>/tmp/sshTest
	botBuildProp=$(cat /tmp/sshTest)
    if [[ "${botBuildProp}" == *"no mutual signature"* ]]; then
	    echo "PubkeyAcceptedKeyTypes +ssh-rsa" >>/etc/ssh/ssh_config
	    botBuildProp=$(ssh -oStrictHostKeyChecking=no -i ${keyPath} root@${botAddress} "cat /build.prop")
    fi
	if [[ ! "${botBuildProp}" == *"ro.build"* ]]; then
		echo "Unable to communicate with robot. The key may be invalid, the bot may not be unlocked, or this device and the robot are not on the same network."
		exit 0
	fi
	scp -oStrictHostKeyChecking=no -o ConnectTimeout=5 -v -i ${keyPath} root@${botAddress}:/build.prop /tmp/scpTest >/tmp/scpTest 2>>/tmp/scpTest
	scpTest=$(cat /tmp/scpTest)
	if [[ "${scpTest}" == *"sftp"* ]]; then
		oldVar="-O"
	else
		oldVar=""
	fi
	if [[ ! "${botBuildProp}" == *"ro.build"* ]]; then
		echo "Unable to communicate with robot. The key may be invalid, the bot may not be unlocked, or this device and the robot are not on the same network."
		exit 0
	fi
    ssh  -oStrictHostKeyChecking=no -i ${keyPath} root@${botAddress} "mount -o rw,remount / && mount -o rw,remount,exec /data && systemctl stop anki-robot.target && mv /anki/data/assets/cozmo_resources/config/server_config.json /anki/data/assets/cozmo_resources/config/server_config.json.bak"
    scp  -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ../vector-cloud/build/vic-cloud root@${botAddress}:/anki/bin/
    scp  -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ../certs/server_config.json root@${botAddress}:/anki/data/assets/cozmo_resources/config/
    scp  -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ../vector-cloud/pod-bot-install.sh root@${botAddress}:/data/
    if [[ -f ./useepod ]]; then
        scp -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ./epod/ep.crt root@${botAddress}:/anki/etc/wirepod-cert.crt
        scp -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ./epod/ep.crt root@${botAddress}:/data/data/wirepod-cert.crt
    else
        scp -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ../certs/cert.crt root@${botAddress}:/anki/etc/wirepod-cert.crt
        scp -oStrictHostKeyChecking=no ${oldVar} -i ${keyPath} ../certs/cert.crt root@${botAddress}:/data/data/wirepod-cert.crt
    fi
    ssh -oStrictHostKeyChecking=no -i ${keyPath} root@${botAddress} "chmod +rwx /anki/data/assets/cozmo_resources/config/server_config.json /anki/bin/vic-cloud /data/data/wirepod-cert.crt /anki/etc/wirepod-cert.crt /data/pod-bot-install.sh && /data/pod-bot-install.sh"
    rm -f /tmp/sshTest
    rm -f /tmp/scpTest
	echo
	echo "Everything has been copied to the bot! Voice commands should work now without needing to reboot Vector."
	echo