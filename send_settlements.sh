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
              "local": 34,
              "backHome": 12,
              "international": 20,
              "premium": 23
            },
            "MTC": 2343
          },
          "SMS": {
            "MO": 5707,
            "MT": 5654
          },
          "Data": [
            {
              "name": "M2M",
              "value": 526263
            }
          ]
        }
      },
      "outbound": {
        "currency": "string",
        "services": {
          "voice": {
            "MOC": {
              "local": 123.98,
              "backHome": 231.7,
              "international": 232,
              "premium": 342,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 223,
            "MT":12334
          },
          "Data": [
            {
              "name": "M2M",
              "value": 735265.8
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
              "local": 123,
              "backHome": 342,
              "international": 2342,
              "premium": 332,
              "ROW": 0
            },
            "MTC": 0
          },
          "SMS": {
            "MO": 234,
            "MT": 123
          },
          "Data": [
            {
              "name": "M2M",
              "value": 432321
            }
          ]
        }
      },
      "outbound": {
        "currency": "string",
        "services": {
          "voice": {
            "MOC": {
              "local": 2332,
              "backHome": 3432,
              "international": 34322,
              "premium": 223,
              "ROW": 0
            },
            "MTC": 233
          },
          "SMS": {
            "MO": 2345,
            "MT": 2313
          },
          "Data": [
            {
              "name": "M2M",
              "value": 123234
            }
          ]
        }
      }
    }
  }
]'
