#!/usr/bin/env bash

shell_users=("root" "ubuntu" "ec2_user" "other")

while getopts 'i:u:' opt; do
  case "$opt" in
  k)
    sshKey="$OPTARG"
    if [[ -f ${sshKey} ]]; then
      echo "Using ssh-key: '${OPTARG}' "
      sshKey="-i ${OPTARG}"
    else
      echo "ssh-key not found: '${OPTARG}'. Not using it"
      unset sshKey
    fi
    ;;
  u)
    user="$OPTARG"
    echo "Using user: '${OPTARG}'"
    ;;
  ?)
    echo -e "Invalid command option.\nUsage: $(basename $0) [-i path_to/ssh_key]"
    exit 1
    ;;
  esac
done
shift "$(($OPTIND - 1))"

# Check AWS credentials via AWS_PROFILE
if [[ -z ${AWS_PROFILE} ]]; then
  echo "!! AWS_PROFILE  not set"
  exit 1
fi

# Check AWS credentials via AWS_PROFILE and AWS_DEFAULT_REGION

# execute AWS CLI-command to retrieve info about EC2-instances. Exit when command fails
aws_output=$(aws ec2 describe-instances 2>/dev/null)
if [[ $? -ne 0 ]]; then
  echo "!! Error listing instances. Check AWS credentials"
  exit 1
fi

# Extract fields needed from JSON_output with the use of jq-command
instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | select(.State.Code != 48) | "\(.InstanceId) - \(.PrivateIpAddress) - \(.PublicIpAddress) - \(.Tags[] | select(.Key=="Name") | .Value)"')

# Check if array constains instance prefix i-. Exit when no running instances are found
if [[ ! ${instances[@]} == *"i-"* ]]; then
  echo "!! No active instances found"
  exit 1
fi

# Print the proper fields
if [[ -z ${shell_users} ]]; then
  shell_users=("root" "ubuntu" "other")
fi

if [[ ${shell_users[@]} != *"other"* ]]; then
  shell_users+=('other')
fi

user=$(gum choose ${shell_users[@]})

if [[ $user == "other" ]]; then
  user=$(gum input)
fi
instance=$(echo "$instances" | gum filter)
instance_name=$(echo $instance | cut -f 1 -d " ")
ssh $sshKey $user@$instance_name
else
echo "!! No AWS_PROFILE or AWS_DEFAULT_REGION set"

# execute AWS CLI-command to retrieve info about EC2-instances. Exit when command fails
aws_output=$(aws ec2 describe-instances 2>/dev/null)
if [[ $? -ne 0 ]]; then
  echo "!! Error listing instances. Check AWS credentials"
  exit 1
fi

# Extract fields needed from JSON_output with the use of jq-command
instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | "\(.InstanceId) - \(.PrivateIpAddress) - \(.PublicIpAddress) - \(.Tags[] | select(.Key=="Name") | .Value)"' | grep -v null)

# Check if array constains instance prefix i-. Exit when no running instances are found
if [[ ! ${instances[@]} == *"i-"* ]]; then
  echo "!! No active instances found"
  exit 1
fi

instance=$(echo "$instances" | gum filter --reverse --prompt="Select instance:> ")
instance_name=$(echo $instance | cut -f 1 -d " ")

# Print the proper fields
if [[ -z ${shell_users} ]]; then
  shell_users=("root" "ubuntu" "other")

fi

if [[ ${shell_users[@]} != *"other"* ]]; then
  shell_users+=('other')
fi

if [[ -z ${user} ]]; then
  user=$(gum choose ${shell_users[@]})
fi

if [[ $user == "other" ]]; then
  user=$(gum input --prompt="Enter username >")
fi

ssh $sshKey $user@$instance_name
