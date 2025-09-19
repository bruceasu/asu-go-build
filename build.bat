@echo off
setLocal
set CWD=%CD%
cd %~dp0
:: build self
md bin
go-build
cd %CWD%
endLocal