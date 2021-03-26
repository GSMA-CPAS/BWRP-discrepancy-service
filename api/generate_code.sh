#!/bin/bash
echo "Generating types..."
oapi-codegen -config=types.cfg.yaml openapi.yaml
sleep 1
echo "Generating server..."
oapi-codegen -config=server.cfg.yaml openapi.yaml
