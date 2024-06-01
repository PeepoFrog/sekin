#!/bin/bash
#
# Define log file paths
SEKAI_LOG="/syslog-data/syslog-ng/logs/shidai.log"
INTERX_LOG="/syslog-data/syslog-ng/logs/interx.log"
SHIDAI_LOG="/syslog-data/syslog-ng/logs/sekai.log"

# Monitoring script to check log file sizes and rotate if necessary
while true; do
      logrotate -v /etc/logrotate.d/logrotate.conf
  sleep 100
done
