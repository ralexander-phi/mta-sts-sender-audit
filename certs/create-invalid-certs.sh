#!/bin/bash

set -e

export CAROOT=`pwd`

# Remove all old leaf certs
rm *alexsci*.pem || true

create_cert() {
    DOMAIN=$1
    DN_PREFIX=$2
    CERT_PREFIX=$3
    ~/code/mkcert/mkcert "${DN_PREFIX}${DOMAIN}"
    rm -Rf $DOMAIN
    mkdir -p $DOMAIN
    mv "${CERT_PREFIX}${DOMAIN}-key.pem" $DOMAIN/key.pem
    mv "${CERT_PREFIX}${DOMAIN}.pem"     $DOMAIN/fullchain.pem
}

for sub in {a,b,c,d}
do
    create_cert "mail-${sub}.audit.alexsci.com" "" ""
done

# This is also used as mta-sts.*.alexsci.com although it isn't valid for those host names
# The whole CA is untrusted so ¯\_(ツ)_/¯
create_cert "audit.alexsci.com" "*." "_wildcard."

