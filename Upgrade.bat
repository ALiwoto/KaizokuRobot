@echo off
git pull && powershell -command "Stop-service -Force -name "KaizokuRobot" -ErrorAction SilentlyContinue; go mod tidy; go build; Start-service -name "KaizokuRobot""
:: Hail Hydra
