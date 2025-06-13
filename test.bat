@echo off

curl -X POST http://localhost:8080/replay -F "file=@F:\Games\Ubisoft\Tom Clancy's Rainbow Six Siege\MatchReplay\Match-2025-06-11_16-46-21-25200.zip" -H "Content-Type: multipart/form-data" -o "output.json"