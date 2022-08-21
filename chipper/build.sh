#!/bin/bash
	export CGO_LDFLAGS="-L/root/.coqui/"
	export CGO_CXXFLAGS="-I/root/.coqui/"
	export LD_LIBRARY_PATH="/root/.coqui/:$LD_LIBRARY_PATH"
UNAME=$(uname -a)
echo "Building chipper..."
/usr/local/go/bin/go build cmd/main.go
mv main chipper
echo "Built chipper!"
