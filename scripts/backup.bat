@echo on
set local

pg_dump players > %USERPROFILE%\players-api\backup\players.db

