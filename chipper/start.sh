#!/bin/bash

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root. sudo ./start.sh"
  exit 1
fi

if [[ -d ./chipper ]]; then
   cd chipper
fi

#if [[ ! -f ./chipper ]]; then
#   if [[ -f ./go.mod ]]; then
#     echo "You need to build chipper first. This can be done with the setup.sh script."
#   else
#     echo "You must be in the chipper directory."
#   fi
#   exit 0
#fi

if [[ -f ./source.sh ]]; then
  source source.sh
fi

if [[ -f ./chipper ]]; then
      	export CGO_CXXFLAGS="-I$HOME/.coqui/"
	export CGO_ENABLED=1
	export CGO_CFLAGS="-I/root/.vosk/libvosk"
	export CGO_LDFLAGS="-L/root/.coqui/ -L /root/.vosk/libvosk -lvosk -ldl -lpthread"
	export LD_LIBRARY_PATH="/root/.coqui/:/root/.vosk/libvosk:$LD_LIBRARY_PATH"
  ./chipper
else
  	export CGO_CXXFLAGS="-I/root/.coqui/"
	export CGO_ENABLED=1
	export CGO_CFLAGS="-I/root/.vosk/libvosk"
	export CGO_LDFLAGS="-L/root/.coqui/ -L /root/.vosk/libvosk -lvosk -ldl -lpthread"
	export LD_LIBRARY_PATH="/root/.coqui/:/root/.vosk/libvosk:$LD_LIBRARY_PATH"
   echo "ERROR: chipper binary does not exist. Please open a Github issue if you see this."
  /usr/local/go/bin/go run cmd/main.go
fi
