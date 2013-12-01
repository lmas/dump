#!/bin/sh
# Create a new hosts file and fill it with blocked adservers.

FILE='/etc/hosts'

echo "# HOSTS file based on http://www.mvps.org/winhelp2002/hosts.txt" > $FILE
echo "" >> $FILE
echo "#<ip-address>   <hostname.domain.org>   <hostname>" >> $FILE
echo "127.0.0.1       localhost" >> $FILE
echo "::1             localhost" >> $FILE
echo "192.168.1.1     gateway $(hostname)" >> $FILE
echo "" >> $FILE

wget -q -O - http://www.mvps.org/winhelp2002/hosts.txt | \
grep -v localhost | \
grep 127.0.0.1 | \
awk '{print "0.0.0.0 " $2}' >> $FILE

echo "0.0.0.0 www.facebook.com" >> $FILE
