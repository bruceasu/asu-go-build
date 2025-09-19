@echo off
setLocal
set CWD=%CD%
cd %~dp0
:: windows, linux, darwin
set GOOS=linux
:: amd64,  386
set GOARCH=amd64
go build -o go-build
cd %CWD%
endLocal