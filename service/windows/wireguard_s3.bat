@echo off
chcp 65001

set base=%~dp0

if not "%1"=="am_admin" (powershell start -verb runas '%0' am_admin & exit /b)
cd /d %base%

set KEY=""
set WG_REMOTE_INTERFACE=""
set WG_INTERFACE=""
set S3_ENDPOINT=""
set S3_ACCESS_KEY_ID=""
set S3_SECRET_ACCESS_KEY=""
set S3_PATH_STYLE=true
set S3_BUCKET=""
set S3_OBJECT_PATH=""

ipsync wireguard s3
