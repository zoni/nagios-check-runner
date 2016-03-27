#!/bin/sh

if [ -x /bin/systemctl ]; then
	/bin/systemctl restart ncr.service
elif [ -x /usr/sbin/service ]; then
	/usr/sbin/service ncr restart
fi
