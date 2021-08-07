#!/bin/bash

## Script to easily copy SSH keys to multiple devices

echo Enter SSH user:
read user
echo Enter SSH password:
read -s password

for ip in `cat inventory`; do
    if [[ "$ip" =~ ^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$ ]]; then
      sshpass -p "$password" ssh-copy-id -i ~/.ssh/id_rsa.pub $user@$ip
    fi
done
