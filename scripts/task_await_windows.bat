@echo off
REM
if "%~1"=="" (
    echo Please provide a wait time in seconds.
    exit /b 1
)

REM
set /a waitTime=%~1 / 1000

REM
echo Awaiting %waitTime% seconds...
ping -n %waitTime% 127.0.0.1 > nul

echo Finished waiting for %waitTime% seconds.