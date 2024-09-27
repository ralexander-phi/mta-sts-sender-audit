#!/bin/bash

echo "Send an email to these addresses:"

USERS=""

for s in a b c d e f
do
  U=$(uuidgen)
  echo "${U}@${s}.audit.alexsci.com, "
  if [[ -z ${USERS} ]];
  then
    USERS="${U}"
  else
    USERS="${USERS},${U}"
  fi
done

while true
do
  echo "Press any key to poll for messages, q to quit"
  echo ""

  read -n 1 key
  if [[ $key == "q" ]]; then
    break
  fi

  curl https://api.audit.alexsci.com/poll -F users="" | jq .
  curl https://api.audit.alexsci.com/poll -F users="${USERS}" | jq .
done

echo ""

