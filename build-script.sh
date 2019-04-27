go install

env GOOS=windows GOARCH=amd64 go build -o nethelp
mkdir -p builds/nethelp-windows
mv nethelp builds/nethelp-windows
cp README.md builds/nethelp-windows
tar -C /Users/maxdobeck/go/src/github.com/mdsauce/nethelp/builds -czf nethelp-windows.tar.gz nethelp-windows

env GOOS=linux GOARCH=amd64 go build -o nethelp
mkdir -p builds/nethelp-linux
mv nethelp builds/nethelp-linux
cp README.md builds/nethelp-linux
tar -C  /Users/maxdobeck/go/src/github.com/mdsauce/nethelp/builds -czf nethelp-linux.tar.gz nethelp-linux


go build -o nethelp
mkdir -p builds/nethelp-mac
mv nethelp builds/nethelp-mac
cp README.md builds/nethelp-mac
tar -C  /Users/maxdobeck/go/src/github.com/mdsauce/nethelp/builds  -czf nethelp-mac.tar.gz nethelp-mac

mv *.tar.gz builds
