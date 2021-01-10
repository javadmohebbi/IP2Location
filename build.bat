
SET GOOS=windows
SET GOARCH=amd64
go build -o ./dist/windows/amd64/ip2location.exe ./ip2location.go ./utils.go ./update.go


SET GOOS=windows
SET GOARCH=386
go build -o ./dist/windows/386/ip2location.exe ./ip2location.go ./utils.go ./update.go


SET GOOS=linux
SET GOARCH=amd64
go build -o ./dist/linux/amd64/ip2location ./ip2location.go ./utils.go ./update.go


SET GOOS=linux
SET GOARCH=386
go build -o ./dist/linux/386/ip2location ./ip2location.go ./utils.go ./update.go


SET GOOS=darwin
SET GOARCH=amd64
go build -o ./dist/darwin/amd64/ip2location ./ip2location.go ./utils.go ./update.go


SET GOOS=darwin
SET GOARCH=386
go build -o ./dist/darwin/386/ip2location ./ip2location.go ./utils.go ./update.go



SET GOOS=linux
SET GOARCH=amd64
