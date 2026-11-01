@echo off
title Instagram Platform Launcher
echo ===================================================
echo     Launching Instagram Full-Stack Servers
echo ===================================================
echo.

echo [1/3] MongoDB Service is running in the background.
echo.

echo [2/3] Starting Go Concurrency Engine (Port 8081)...
start "Instagram - Go Concurrency Service (Port 8081)" cmd /k "cd /d "%~dp0backend-concurrency" && concurrency-service.exe"

echo [3/3] Starting Java 21 Spring Boot Core (Port 8080)...
start "Instagram - Java 21 Spring Boot Core (Port 8080)" cmd /k "cd /d "%~dp0backend-core" && mvn spring-boot:run"

echo.
echo ===================================================
echo All services launched!
echo - Core Security & Auth API : http://localhost:8080
echo - Go Concurrency WebSockets: ws://localhost:8081/ws
echo.
echo Opening Instagram in your default browser...
start "" "%~dp0index.html"
echo ===================================================
pause
