env GOOS=windows GOARCH=amd64 go build -o nethelp
mkdir builds/nethelp-windows
mv nethelp builds/nethelp-windows
cp README.md builds/nethelp-windows

env GOOS=linux GOARCH=amd64 go build -o nethelp
mkdir builds/nethelp-linux
mv nethelp builds/nethelp-linux
cp README.md builds/nethelp-linux

go build -o nethelp
mkdir builds/nethelp-mac
mv nethelp builds/nethelp-mac
cp README.md builds/nethelp-mac
