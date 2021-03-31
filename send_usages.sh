#!/bin/bash

curl -X 'PUT' \
  'http://127.0.0.1:8080/usages/6056f30bc257e800281964f75f1b?partnerUsageId=6056f30bc257e800281964f75f1b' \
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
          "yearMonth": "202001",
          "homeTadig": "FRAF1",
          "visitorTadig": "DEUD1",
          "service": "MOC Local",
          "usage": 23573.98,
    	  "unit": "min"
        },
        {
          "yearMonth": "202002",
          "homeTadig": "BELMO",
          "visitorTadig": "DEUD1",
          "service": "GPRS",
          "usage": 48740.74,
    	  "unit": "MB"
        }
      ],
      "outbound": [
        {
          "yearMonth": "202002",
          "homeTadig": "DEUD1",
          "visitorTadig": "BELMO",
          "service": "GPRS",
          "usage": 11787.03,
          "unit": "MB"
        },
        {
          "yearMonth": "202001",
          "homeTadig": "DEUD1",
          "visitorTadig": "FRAF1",
          "service": "MOC Local",
          "usage": 45085.07,
          "unit": "min"
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
          "yearMonth": "202002",
          "homeTadig": "DEUD1",
          "visitorTadig": "BELMO",
          "service": "GPRS",
          "usage": 11728.39,
          "unit": "MB"  
        },
        {
          "yearMonth": "202001",
          "homeTadig": "DEUD1",
          "visitorTadig": "FRAF1",
          "service": "MOC Local",
          "usage": 45678.9,
          "unit": "min"  
        }
      ],
      "outbound": [
        {
          "yearMonth": "202001",
          "homeTadig": "FRAF1",
          "visitorTadig": "DEUD1",
          "service": "MOC Local",
          "usage": 23456.7,
          "unit": "min"
        },
        {
          "yearMonth": "202002",
          "homeTadig": "BELMO",
          "visitorTadig": "DEUD1",
          "service": "GPRS",
          "usage": 49382.72,
          "unit": "MB"
        }
      ]
    }
  }
]'
