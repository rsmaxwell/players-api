@echo off

set mydate=%date:~6,4%-%date:~3,2%-%date:~0,2%T%time:~0,8%Z
set mydate=%mydate: =0%
echo %mydate%
