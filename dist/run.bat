@echo off
setlocal

echo Enter Instance Hostname (Domain only):
set /p HOSTNAME=

echo Enter API Token:
set /p API_TOKEN=

echo Enable Debug? (true or false):
set /p DEBUG=

set HOSTNAME=%HOSTNAME%
set API_TOKEN=%API_TOKEN%
set DEBUG=%DEBUG%

echo Environment variables are set. Running MMSB.exe.
MMSB.exe

endlocal
