#!/bin/bash

# Wildcard cert used for:
#      mail-X.audit.alexsci.com
#         api.audit.alexsci.com
# Also covers:
#   mta-sts.X.audit.alexsci.com
sudo certbot certonly \
      --dns-digitalocean \
      --dns-digitalocean-credentials ~/certbot-creds.ini \
      -d '*.audit.alexsci.com'
      -d 'mta-sts.a.audit.alexsci.com'
      -d 'mta-sts.b.audit.alexsci.com'
      -d 'mta-sts.c.audit.alexsci.com'
      -d 'mta-sts.d.audit.alexsci.com'

