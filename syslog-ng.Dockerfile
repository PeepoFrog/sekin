# Use an existing syslog-ng image as a base
FROM linuxserver/syslog-ng

# Install necessary packages
RUN apk add --no-cache logrotate inotify-tools

# Copy custom logrotate config and monitoring script from the host to the container
COPY config/logrotate.conf /etc/logrotate.conf
COPY config/monitoring.sh /usr/local/bin/monitoring.sh

# Ensure the monitoring script is executable
RUN chmod +x /usr/local/bin/monitoring.sh

# Set up environment to run the monitoring script in the background
CMD ["/usr/local/bin/monitoring.sh"]

