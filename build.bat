set GO111MODULE=on
set GOARCH=386
rsrc -manifest main.manifest -ico serviceMonitor.ico -o rsrc.syso
go build -ldflags="-s -w -H=windowsgui"
upx -9 ./serviceMonitor.exe
