#!/bin/bash

VERSION=$(git describe --tags)

mkdir -p dist

# build lakitu-cli
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X 'main.version=$VERSION'" -o bin/lakitu-cli.exe cmd/cli/main.go
zip -j dist/lakitu-cli_$VERSION.zip bin/lakitu-cli.exe

# build front end
cd web
npm install
npm run build

# move build folder so go can use
if [ -d "../assets/build" ]; then rm -rf ../assets/build/; fi
mv build/ ../assets/

cd ..

# build executables
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X 'main.CookieStoreSecret=$COOKIE_SECRET'" -o bin/lakitu-windows_amd64.exe cmd/server/main.go
zip -j dist/lakitu_$VERSION-windows.zip bin/lakitu-windows_amd64.exe

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.CookieStoreSecret=$COOKIE_SECRET'" -o bin/lakitu-linux_amd64 cmd/server/main.go
zip -j dist/lakitu_$VERSION-linux.zip bin/lakitu-linux_amd64
