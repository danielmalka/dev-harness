@echo off
set "binary=%~dp0windows_amd64\dh.exe"
if not exist "%binary%" (
  >&2 echo dh: no dh binary for windows_amd64; run dh doctor
  exit /b 0
)
"%binary%" %*
