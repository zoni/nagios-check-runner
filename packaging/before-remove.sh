#!/bin/sh

if [ -x /bin/systemctl ]; then
	/bin/systemctl stop ncr.service
elif [ -x /usr/sbin/service ]; then
	/usr/sbin/service ncr stop
fi
