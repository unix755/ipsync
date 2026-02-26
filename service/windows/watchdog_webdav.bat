@echo off
chcp 65001

set base=%~dp0

if not "%1"=="am_admin" (powershell start -verb runas '%0' am_admin & exit /b)
cd /d %base%

set KEY=""
set WG_REMOTE_INTERFACE=""
set WG_INTERFACE=""
set WEBDAV_ENDPOINT=""
set WEBDAV_USERNAME=""
set WEBDAV_PASSWORD=""
set WEBDAV_PATH=""

ipsync watchdog webdav
