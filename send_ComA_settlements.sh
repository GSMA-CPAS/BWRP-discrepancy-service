curl -X 'PUT' \
  'http://127.0.0.1:8080/settlements/1?partnerSettlementId=2' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d @sett_ComA_payload.json
