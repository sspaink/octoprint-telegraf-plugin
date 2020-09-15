#!/bin/bash

echo "Building the plugin..."
env GOOS=linux GOARCH=arm GOARM=5 go build -o bin/octoprint cmd/main.go
echo "Deploying the plugin..."
ssh pi "rm ~/plugins/octoprint"
scp bin/octoprint pi:~/plugins/octoprint
ssh pi "chmod +x ~/plugins/octoprint"