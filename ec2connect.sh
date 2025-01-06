#!/usr/bin/env bash

shell_users=("root" "ubuntu" "ec2_user" "other")

while getopts 'i:h:u:' opt; do
	case "$opt" in
		i)
			sshKey="$OPTARG"
			if [[ -f ${sshKey} ]]; then
				echo "Using ssh-key: '${OPTARG}' "
				sshKey="-i ${OPTARG}"
			else
				echo "ssh-key not found: '${OPTARG}'. Not using it"
				unset sshKey
			fi
			;;
    h)
      instance_id=${OPTARG}
      ;;
		u)
			user="$OPTARG"
			echo "Using user: '${OPTARG}'"
			;;
		?)
			echo -e "Invalid command option.\nUsage: $(basename $0) [-i path_to/ssh_key]i [-u username]"
			exit 1
			;;
	esac
done
shift "$(($OPTIND - 1))"



aws_output=$(aws ec2 describe-instances)
if [ $? -ne 0 ]; then
	echo "!! Error: Can't build list of instances. Check error above."
	exit 1
fi

header=$(echo "$aws_output" | jq -r '["InstanceId", "PrivateIpAddress", "PublicIpAddress", "State", "Name"] | @csv')
instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | select(.State.Code != 48) | [
.InstanceId,
.PrivateIpAddress,
.PublicIpAddress // "null",
.State.Name,
(.Tags[] | select(.Key=="Name") | .Value)
] | @csv')

formatted_instances=$(echo -e "$header\n$instances")

if [[ -z ${instance_id} ]]; then
  instance_id=$(echo -e "$formatted_instances" | gum table -w 20,16,16,8,50 --height 20 | awk -F, '{print $1}')
fi

instance_state=$(echo "$aws_output" | jq -r --arg INSTANCE_ID "$instance_id" '.Reservations[].Instances[] | select(.InstanceId == $INSTANCE_ID) | .State.Name')
if [ "$instance_state" != "running" ]; then
  echo "!!! Instance chosen is not running. Current state is : ${instance_state}. Not executing the SSH-command"
  exit 1
fi

connectionMethod=$(gum choose "ssh" "ssm")

if [[ ${connectionMethod} == "ssh" ]]; then
  while [[ -z ${user} ]] ; do
    header=$(echo -e "Available Users")
    users_list=$(printf "%s\n" "${shell_users[@]}")
    user=$(echo -e "$header\n$users_list" | gum table -w 20 | awk '{print $1}')
    while [[ ${user} == "other" ]]; do
      user=$(gum input --prompt="Enter username >")
    done
  done

  echo "-- Executing: ssh ${sshKey} ${user}@${instance_id}"
  ssh ${sshKey} ${user}@${instance_id}
fi

if [[ ${connectionMethod} == "ssm" ]]; then
  echo "-- Executing: aws ssm start-session --target ${instance_id}"
  aws ssm start-session --target ${instance_id}
fi


echo "## END"
