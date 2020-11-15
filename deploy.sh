#!/bin/bash

echo "Building the plugin..."
make build
echo "Deploying the plugin..."
make deploy
echo "Reloading telegraf to start latest plugin..."
make reload