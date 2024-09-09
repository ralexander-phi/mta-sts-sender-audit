#!/bin/bash

set -e

if [ ! -f /postgres-password.txt ];
then
    echo $POSTGRES_PASSWORD > /postgres-password.txt
fi

/usr/sbin/postfix -c /etc/postfix -vv start-fg

