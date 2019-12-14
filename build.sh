#!/bin/sh

env GOOS=windows GOARCH=amd64 go build -o ./dist/windows/amd64/ip2location.exe ./ip2location.go ./utils.go ./update.go


env GOOS=windows GOARCH=386 go build -o ./dist/windows/386/ip2location.exe ./ip2location.go ./utils.go ./update.go


env GOOS=linux GOARCH=amd64 go build -o ./dist/linux/amd64/ip2location ./ip2location.go ./utils.go ./update.go

env GOOS=linux GOARCH=386 go build -o ./dist/linux/386/ip2location ./ip2location.go ./utils.go ./update.go

env GOOS=darwin GOARCH=amd64 go build -o ./dist/darwin/amd64/ip2location ./ip2location.go ./utils.go ./update.go

env GOOS=darwin GOARCH=386 go build -o ./dist/darwin/386/ip2location ./ip2location.go ./utils.go ./update.go

env GOOS=linux GOARCH=amd64
