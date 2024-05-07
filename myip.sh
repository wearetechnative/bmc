#!/usr/bin/env bash

public_ip=$(curl -s https://api.ipify.org)
echo "Your public IP is: $public_ip"

