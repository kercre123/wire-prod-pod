#!/bin/bash

if [[ -f /usr/bin/apt ]]; then
	apt install -y libasound2-dev
fi

if [[ ! -d /root/.vosk ]]; then
    origDir=$(pwd)
    echo "Getting VOSK assets"
    rm -fr /root/.vosk
    mkdir /root/.vosk
    cd /root/.vosk
    if [[ ${ARCH} == "x86_64" ]]; then
        VOSK_DIR="vosk-linux-x86_64-0.3.43"
    elif [[ ${ARCH} == "aarch64" ]]; then
        VOSK_DIR="vosk-linux-aarch64-0.3.43"
    elif [[ ${ARCH} == "armv7l" ]]; then
        VOSK_DIR="vosk-linux-armv7l-0.3.43"
    fi
    VOSK_ARCHIVE="$VOSK_DIR.zip"
    wget -q --show-progress "https://github.com/alphacep/vosk-api/releases/download/v0.3.43/$VOSK_ARCHIVE"
    unzip "$VOSK_ARCHIVE"
    mv "$VOSK_DIR" libvosk
    rm -fr "$VOSK_ARCHIVE"
    cd ${origDir}/chipper
    export CGO_ENABLED=1
    export CGO_CFLAGS="-I/root/.vosk/libvosk"
    export CGO_LDFLAGS="-L /root/.vosk/libvosk -lvosk -ldl -lpthread"
    export LD_LIBRARY_PATH="/root/.vosk/libvosk:$LD_LIBRARY_PATH"
    /usr/local/go/bin/go get -u github.com/alphacep/vosk-api/go/...
    /usr/local/go/bin/go get github.com/alphacep/vosk-api
    /usr/local/go/bin/go install github.com/alphacep/vosk-api/go
    cd ${origDir}
    rm -fr vosk
    mkdir -p vosk
    mkdir -p vosk/models
    echo "Downloading English (US) model"
    mkdir -p vosk/models/en-US
    cd vosk/models/en-US
    wget https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip
    unzip vosk-model-small-en-us-0.15.zip
    mv vosk-model-small-en-us-0.15 model
    rm vosk-model-small-en-us-0.15.zip
    cd ${origDir}
    cd ${origDir}/vosk
    touch completed
    echo
    cd ..
else
    echo "postupdate - Nothing to be done!"
fi
