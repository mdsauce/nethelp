env GOOS=windows GOARCH=amd64 go build -o nethelp
mkdir builds/nethelp-windows
mv nethelp builds/nethelp-windows
cp README.md builds/nethelp-windows
tar -czf nethelp-windows.tar.gz builds/nethelp-windows

env GOOS=linux GOARCH=amd64 go build -o nethelp
mkdir builds/nethelp-linux
mv nethelp builds/nethelp-linux
cp README.md builds/nethelp-linux
tar -czf nethelp-linux.tar.gz builds/nethelp-linux


go build -o nethelp
mkdir builds/nethelp-mac
mv nethelp builds/nethelp-mac
cp README.md builds/nethelp-mac
tar -czf nethelp-mac.tar.gz builds/nethelp-mac

mv *.tar.gz builds