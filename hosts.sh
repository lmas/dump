#!/bin/sh
# Create a new hosts file and fill it with blocked adservers.

FILE="/etc/hosts"
REDIRECT_IP="127.0.0.1"

echo "# HOSTS file based on http://www.mvps.org/winhelp2002/hosts.txt" > $FILE
echo "" >> $FILE
echo "#<ip-address>   <hostname.domain.org>   <hostname>" >> $FILE
echo "127.0.0.1       localhost" >> $FILE
echo "::1             localhost" >> $FILE
echo "192.168.1.1     gateway $(hostname)" >> $FILE
echo "" >> $FILE

wget -q -O - http://www.mvps.org/winhelp2002/hosts.txt | \
grep -v localhost | \
grep 0.0.0.0 | \
awk -v ip=$REDIRECT_IP '{print ip " " $2}' >> $FILE

echo "$REDIRECT_IP www.facebook.com" >> $FILE
