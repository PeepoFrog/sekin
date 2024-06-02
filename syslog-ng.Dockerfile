# Use an existing syslog-ng image as a base
FROM linuxserver/syslog-ng

# Install necessary packages
RUN apk add --no-cache logrotate 

# Copy custom logrotate config and monitoring script from the host to the container
COPY config/logrotate.conf /etc/logrotate.d/logrotate.conf
COPY config/cron_logrotate /etc/crontabs/
COPY config/syslog-ng.conf /etc/syslog-ng/syslog-ng.conf

# Ensure the monitoring script is executable
RUN chown root:root /etc/logrotate.d/logrotate.conf && \
	chmod 0640 /etc/logrotate.d/logrotate.conf && \
	chown root:root /etc/syslog-ng/syslog-ng.conf && \
  chmod 0640 /etc/syslog-ng/syslog-ng.conf

RUN crontab /etc/crontabs/cron_logrotate

CMD ["syslog-ng","-F","-c","/run/syslog-ng.ctl"]
