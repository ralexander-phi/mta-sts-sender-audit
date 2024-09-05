#!/bin/bash
export CAROOT=`pwd`

for sub in {a,b,c,d}
do
	~/code/mkcert/mkcert $sub.audit.alexsci.com
	cat $sub.audit.alexsci.com-key.pem >  $sub.pem
	cat $sub.audit.alexsci.com.pem     >> $sub.pem
	#cat rootCA.pem                     >> $sub.pem
done

~/code/mkcert/mkcert "*.audit.alexsci.com"
cat _wildcard.audit.alexsci.com-key.pem >  _wildcard.pem
cat _wildcard.audit.alexsci.com.pem     >> _wildcard.pem
