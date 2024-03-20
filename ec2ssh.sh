#!/usr/bin/env bash

# Check AWS credentials via AWS_PROFILE and AWS_DEFAULT_REGION
if [[ ! -z ${AWS_PROFILE} || ! -z ${AWS_DEFAULT_REGION} ]]; then


  # execute AWS CLI-command to retrieve info about EC2-instances. Exit when command fails
  aws_output=$(aws ec2 describe-instances 2>/dev/null)
  if [[ $? -ne 0 ]]; then
    echo "!! Error listing instances. Check AWS credentials"
    exit 1
  fi


  # Extract fields needed from JSON_output with the use of jq-command
  instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | "\(.InstanceId) - \(.PrivateIpAddress) - \(.PublicIpAddress) - \(.Tags[] | select(.Key=="Name") | .Value)"')

  # Check if array constains instance prefix i-. Exit when no running instances are found
  if [[ ! ${instances[@]} == *"i-"* ]]; then
    echo "!! No active instances found"
    exit 1
  fi

  # Print the proper fields
  user=$(gum choose "root" "ubuntu" "other")

  if [[ $user == "other" ]]; then
    user=$(gum input)
  fi
  instance=$(echo "$instances" | gum filter)
  instance_name=$(echo $instance | cut -f 1 -d " ")
  ssh $user@$instance_name
else
  echo "!! No AWS_PROFILE or AWS_DEFAULT_REGION set"
fi
