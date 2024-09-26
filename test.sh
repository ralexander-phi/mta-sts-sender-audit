#!/usr/bin/bash

set -e

docker compose --env-file dev.env build
docker compose --env-file dev.env down || true

# Purge old DB
docker volume rm postfix-tls-audit_postfix-audit-db || true

docker compose --env-file dev.env up -d --wait --remove-orphans


subdomains=(a b c d e f g)
ips=(127.0.0.1 127.0.0.2 "::1" 127.0.0.10 127.0.0.11 127.0.0.12 127.0.0.13)
for i in {0..6}
do
    subdomain=${subdomains[$i]}
    ip=${ips[$i]}
    echo "Checking server ($i / $ip) $subdomain"
    UUID=$(uuidgen)
    echo "Using USERID: ${UUID}"

    # Ensure the MTA-STS policy is available on both HTTP and HTTPS
    # Not on 5th, or 6th
    if (( $i != 5 && $i != 6 ));
    then
        curl -k -H "Host: mta-sts.$subdomain.audit.alexsci.com" https://127.0.0.1:8443/.well-known/mta-sts.txt | grep "enforce"
	# Make sure it was logged
        curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users= | grep "mta-sts.${subdomain}.audit.alexsci.com"
    else
        echo "$subdomain won't have a policy hosted"
    fi

    echo "Checking that email hasn't been seen"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/health | grep "pong"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users=$UUID | grep "{}"

    echo "Send the emails"
    if (( $i == 2 ));
    then
        # This one uses IPv6
        ./test-send-email.exp "[${ip}]" "${UUID}" "${subdomain}.audit.alexsci.com"
    elif (( $i == 3 || $i == 5));
    then
        # These ones doesn't support TLS
        ./test-send-email-no-tls.exp "${ip}" "${UUID}" "${subdomain}.audit.alexsci.com"
    else
        ./test-send-email.exp "${ip}" "${UUID}" "${subdomain}.audit.alexsci.com"
    fi

    if (( $i != 4 && $i != 6));
    then
      # All but the 4th and 6th support unencrypted emails
      # Try to send an email to an unrelated domain (should fail)
      ./test-open-relay.exp "${ip}" "${UUID}" "${subdomain}.audit.alexsci.com"
    fi

    # Email processing takes some time...
    sleep 1

    echo "Checking that email has been seen"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users=$UUID | grep "$UUID" | grep "Message Received"
    curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users=$UUID | grep "$UUID" | grep "MSG: This Is The Message"

    echo ""
    echo "Server $subdomain looks OK!"
    echo ""
done

# Check TLS reporting
curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/tlsrpt -d "TLS REPORT"
curl -k -H "Host: api.audit.alexsci.com" https://127.0.0.1:8443/poll -F users= | grep "TLS REPORT"

docker compose --env-file dev.env down
echo ""
echo "SUCCESS!"
echo ""
