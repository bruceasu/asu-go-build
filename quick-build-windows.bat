@echo off
setLocal
set CWD=%CD%
cd %~dp0
:: windows, linux, darwin
set GOOS=windows
:: amd64,  386
set GOARCH=amd64
go build -o go-build.exe
call cpfile go-build.exe c:\green\
call cpfile go-build.exe c:\green\go\bin
cd %CWD%
endLocal