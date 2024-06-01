#!/bin/bash
#
# Define log file paths
SEKAI_LOG="/syslog-data/syslog-ng/logs/shidai.log"
INTERX_LOG="/syslog-data/syslog-ng/logs/interx.log"
SHIDAI_LOG="/syslog-data/syslog-ng/logs/sekai.log"

# Monitoring script to check log file sizes and rotate if necessary
while true; do
  # Wait for modifications on any of the specified log files
  inotifywait -e modify $SEKAI_LOG $INTERX_LOG $SHIDAI_LOG

  # Check and rotate $SEKAI_LOG if it exceeds 12MB
  if [ -e "$SEKAI_LOG" ]; then
    if [ $(stat -c %s $SEKAI_LOG) -ge 12582912 ]; then
      logrotate /etc/logrotate.d/mylogs
    fi
  else
    echo "Warning: $SEKAI_LOG does not exist."
  fi

  # Check and rotate $INTERX_LOG if it exceeds 12MB
  if [ -e "$INTERX_LOG" ]; then
    if [ $(stat -c %s $INTERX_LOG) -ge 12582912 ]; then
      logrotate /etc/logrotate.d/mylogs
    fi
  else
    echo "Warning: $INTERX_LOG does not exist."
  fi

  # Check and rotate $SHIDAI_LOG if it exceeds 12MB
  if [ -e "$SHIDAI_LOG" ]; then
    if [ $(stat -c %s $SHIDAI_LOG) -ge 12582912 ]; then
      logrotate /etc/logrotate.d/mylogs
    fi
  else
    echo "Warning: $SHIDAI_LOG does not exist."
  fi

  # Wait for 60 minutes before next check
  sleep 3600
done
