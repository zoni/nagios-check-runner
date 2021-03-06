#!/bin/sh

getent group nagios > /dev/null || groupadd --system nagios
useradd --system ncr --home /opt/ncr --groups nagios

if [ -x /bin/systemctl ]; then
	/bin/systemctl start ncr.service
elif [ -x /usr/sbin/service ]; then
	/usr/sbin/service ncr start
fi
