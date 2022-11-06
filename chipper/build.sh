#!/bin/bash
  	export CGO_CXXFLAGS="-I/root/.coqui/"
	export CGO_ENABLED=1
	export CGO_CFLAGS="-I/root/.vosk/libvosk"
	export CGO_LDFLAGS="-L/root/.coqui/ -L /root/.vosk/libvosk -lvosk -ldl -lpthread"
	export LD_LIBRARY_PATH="/root/.coqui/:/root/.vosk/libvosk:$LD_LIBRARY_PATH"
UNAME=$(uname -a)
echo "Building chipper..."
/usr/local/go/bin/go build cmd/main.go
mv main chipper
echo "Built chipper!"
