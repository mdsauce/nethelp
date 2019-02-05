env GOOS=windows GOARCH=amd64 go build -o nethelp-windows
env GOOS=linux GOARCH=amd64 go build -o nethelp-linux
go build -o nethelp-mac

mv nethelp* builds
