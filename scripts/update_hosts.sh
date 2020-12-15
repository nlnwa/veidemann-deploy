#!/usr/bin/env sh

# Update /etc/hosts with hostname and ip

IP=$1
HOSTNAME=$2
shift 2
ALIAS="$@"


if [ -z "$IP" ]; then
  echo missing ip
  exit
fi

if [ -z "$HOSTNAME" ]; then
  echo missing hostname
  exit
fi

echo "Adding mapping from $HOSTNAME to $IP in /etc/hosts"
# Remove old entries for HOSTNAME in /etc/hosts
sudo sed -i "/${HOSTNAME}/d" /etc/hosts

# Add entry for HOSTNAME in /etc/hosts pointing to IP
echo "${IP}" "${HOSTNAME}" "${ALIAS}" | sudo tee -a /etc/hosts
