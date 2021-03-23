curl -X 'PUT' \
  'http://127.0.0.1:8080/settlements/1?partnerSettlementId=2' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '[
  {
    "header": {
      "version": "string",
      "type": "string",
      "mspOwner": "DTAG",
      "context": "home"
    },
    "body": {
      "fromdate": "string",
      "todate": "string",
      "inbound": {
        "currency": "string",
        "services": {
          "voice": {
            "MOC": {
              "local": 0,
              "backHome": 0,
              "international": 0,
              "premium": 0,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 0,
            "MT": 0
          },
          "Data": [
            {
              "name": "string",
              "value": 0
            }
          ]
        }
      },
      "outbound": {
        "currency": "string",
        "services": {
          "voice": {
            "MOC": {
              "local": 0,
              "backHome": 0,
              "international": 0,
              "premium": 0,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 0,
            "MT": 0
          },
          "Data": [
            {
              "name": "string",
              "value": 0
            }
          ]
        }
      }
    }
  },
  {
    "header": {
      "version": "string",
      "type": "string",
      "mspOwner": "ORANGE",
      "context": "partner"
    },
    "body": {
      "fromdate": "string",
      "todate": "string",
      "inbound": {
        "currency": "EURO",
        "services": {
          "voice": {
            "MOC": {
              "local": 0,
              "backHome": 0,
              "international": 0,
              "premium": 0,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 0,
            "MT": 0
          },
          "Data": [
            {
              "name": "string",
              "value": 0
            }
          ]
        }
      },
      "outbound": {
        "currency": "string",
        "services": {
          "voice": {
            "MOC": {
              "local": 0,
              "backHome": 0,
              "international": 0,
              "premium": 0,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 0,
            "MT": 0
          },
          "Data": [
            {
              "name": "string",
              "value": 0
            }
          ]
        }
      }
    }
  }
]'
