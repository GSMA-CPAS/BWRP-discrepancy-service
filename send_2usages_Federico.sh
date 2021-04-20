#!/bin/bash

curl -X 'PUT' \
  'http://127.0.0.1:8080/usages/6056f30bc257e800281964f75f1b?partnerUsageId=6056f30bc257e800281964f75f1b' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d @usages_Federico_payload.json
