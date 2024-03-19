#!/usr/bin/env bash

# Uitvoeren van het AWS CLI-commando om informatie over EC2-instances op te halen
aws_output=$(aws ec2 describe-instances)

# Extract de gewenste velden uit de JSON-output met behulp van jq
instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | "\(.InstanceId) - \(.PrivateIpAddress) - \(.PublicIpAddress) - \(.Tags[] | select(.Key=="Name") | .Value)"')

# Print de gewenste velden
user=$(gum choose "root" "ubuntu" "other")

if [[ $user == "other" ]]; then
  user=$(gum input)
fi
instance=$(echo "$instances" | gum filter)
instance_name=$(echo $instance | cut -f 1 -d " ")
ssh $user@$instance_name
