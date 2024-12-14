#!/bin/sh

# Replace environment variables in config file
envsubst < config.yaml.template > config.yaml

# Start the application
exec ./main
