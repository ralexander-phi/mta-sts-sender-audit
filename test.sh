#!/usr/bin/bash

set -e

docker compose --env-file dev.env build
docker compose --env-file dev.env down || true
docker compose --env-file dev.env up -d --wait

subdomains=(a b c d)
for i in {1..4}
do
    subdomain=${subdomains[$i-1]}
    echo "Checking server ($i) $subdomain"
    UUID=$(uuidgen)
    echo "Using USERID: ${UUID}"

    # Ensure the MTA-STS policy is available on both HTTP and HTTPS
    # Not on third server
    if (( $i != 3 ));
    then
    	curl -k -H "Host: $i.audit.alexsci.com" https://127.0.0.1:8443/.well-known/mta-sts.txt | grep "enforce"
    else
        echo "C won't have a policy hosted"
    fi

    echo "Checking that email hasn't been seen"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users=$UUID | grep "false"

    echo "Send the emails"
    ./test-send-email.exp 127.0.0.$i $UUID $subdomain.audit.alexsci.com

    # Email processing takes some time...
    sleep 1

    echo "Checking that email has been seen"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users=$UUID | grep "true"

    echo ""
    echo "Server $subdomain looks OK!"
    echo ""
done

docker compose --env-file dev.env down
echo ""
echo "SUCCESS!"
echo ""
