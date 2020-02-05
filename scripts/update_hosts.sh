#!/usr/bin/env sh

# Update /etc/hosts with hostname and ip

HOSTNAME=$1
IP=$2

if [ -z "$HOSTNAME" ]; then
  echo missing hostname
  exit
fi

if [ -z "$IP" ]; then
  echo missing ip
  exit
fi

HOST_LINE="${IP} ${HOSTNAME}"
grep "${HOST_LINE}" /etc/hosts
if [ $? == 0 ]; then
  exit
fi

echo "Adding mapping from $HOSTNAME to $IP in /etc/hosts"
# Remove old entries for HOSTNAME in /etc/hosts
sudo sed -i "/${HOSTNAME}/d" /etc/hosts

# Add entry for HOSTNAME in /etc/hosts pointing to IP
echo "${HOST_LINE}" | sudo tee -a /etc/hosts
