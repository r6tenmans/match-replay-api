@echo off

curl -X POST http://localhost:8080/replay -F "file=@F:\Games\Ubisoft\Tom Clancy's Rainbow Six Siege\MatchReplay\Match-2021-07-22_15-06-35-99.zip" -H "Content-Type: multipart/form-data" -o "output.json"