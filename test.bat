@echo off

curl -X POST http://localhost:8080/replay -F "file=@F:\Games\Ubisoft\Tom Clancy's Rainbow Six Siege\MatchReplay\Match-2025-01-25_15-09-34-31184.zip" -H "Content-Type: multipart/form-data" -o "output.json"