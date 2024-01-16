curl -X POST localhost:8080/update \
   -H 'Content-Type: application/json' \
   -d '{"id":"test_gauge","type":"gauge", "value":7.77}' \
