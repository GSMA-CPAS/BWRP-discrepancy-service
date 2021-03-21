#!/bin/bash

curl -X 'PUT' \
  'http://127.0.0.1:8080/usages/1?partnerUsageId=2' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '[
  {
    "header": {
      "version": "1.0",
      "type": "usage",
      "mspOwner": "DTAG",
      "context": "home"
    },
    "body": {
      "inbound": [
        {
          "yearMonth": "string",
          "homeTadig": "string",
          "visitorTadig": "string",
          "service": "string",
          "usage": 0
        }
      ],
      "outbound": [
        {
          "yearMonth": "string",
          "homeTadig": "string",
          "visitorTadig": "string",
          "service": "string",
          "usage": 0
        }
      ]
    }
  },
  {
    "header": {
      "version": "1.0",
      "type": "usage",
      "mspOwner": "ORANGE",
      "context": "partner"
    },
    "body": {
      "inbound": [
        {
          "yearMonth": "string",
          "homeTadig": "string",
          "visitorTadig": "string",
          "service": "string",
          "usage": 0
        }
      ],
      "outbound": [
        {
          "yearMonth": "string",
          "homeTadig": "string",
          "visitorTadig": "string",
          "service": "string",
          "usage": 0
        }
      ]
    }
  }
]'
