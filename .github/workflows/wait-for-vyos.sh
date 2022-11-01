#!/bin/bash

HOST="https://localhost/retrieve"
KEY="vyos"

STATUS=000
while [[ "$STATUS" != "200" ]]; do
    echo "Waiting for VyOS... status=$STATUS"

    STATUS=$(curl -sk -o /dev/null -w "%{http_code}" -X POST $HOST -F data='{"op": "showConfig", "path": []}' -F key="$KEY")
    curl -k -D - -X POST $HOST -F data='{"op": "showConfig", "path": ["system", "host-name"]}' -F key="$KEY"

    sleep 5
done