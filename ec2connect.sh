#!/usr/bin/env bash

# Load config file for BMC_AUTO_START_STOPPED_INSTANCES setting
conffile="${HOME}/.config/bmc/config.env"
if [ -f ${conffile} ]; then
  source $conffile
fi

# Default to "prompt" if not set
BMC_AUTO_START_STOPPED_INSTANCES=${BMC_AUTO_START_STOPPED_INSTANCES:-"prompt"}

# Helper function to start an instance and wait for it to be running
startInstanceAndWait() {
  local instance_id=$1
  local bmcpath=$(dirname $0)

  # Source _bmclib.sh to access ec2CheckNewInstanceState function
  source ${bmcpath}/_bmclib.sh

  echo "-- Starting instance ${instance_id}..."
  aws ec2 start-instances --instance-ids ${instance_id} >/dev/null
  if [ $? -ne 0 ]; then
    echo "!! Error: Failed to start instance. Check error above."
    exit 1
  fi

  gum spin --spinner meter --title "Starting instance ${instance_id}" -- bash -c "source ${bmcpath}/_bmclib.sh && ec2CheckNewInstanceState ${instance_id} running"
  echo "-- Instance ${instance_id} is now running"
}

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

# Handle non-running instances
if [ "$instance_state" != "running" ]; then
  if [ "$instance_state" = "stopped" ]; then
    # Instance is stopped - handle based on config
    case "$BMC_AUTO_START_STOPPED_INSTANCES" in
      always)
        # Auto-start without prompting
        startInstanceAndWait "$instance_id"
        ;;
      never)
        # Exit with error message
        echo "!!! Instance chosen is not running. Current state is : ${instance_state}."
        exit 1
        ;;
      *)
        # Prompt user (default behavior)
        if gum confirm "Instance ${instance_id} is stopped. Start it?"; then
          startInstanceAndWait "$instance_id"
        else
          echo "Instance not started. Exiting."
          exit 0
        fi
        ;;
    esac
  else
    # Instance is in another non-running state (pending, stopping, etc.)
    echo "!!! Instance chosen is not running. Current state is : ${instance_state}."
    exit 1
  fi
fi

# Auto-select SSH if -u or -i flags were provided
if [[ -n ${user} || -n ${sshKey} ]]; then
  connectionMethod="ssh"
else
  connectionMethod=$(gum choose "ssh" "ssm")
fi

if [[ ${connectionMethod} == "ssh" ]]; then
  while [[ -z ${user} ]] ; do
    header=$(echo -e "Available Users")
    users_list=$(printf "%s\n" "${shell_users[@]}")
    user=$(echo -e "$header\n$users_list" | gum table -w 20 --height 20 | awk '{print $1}')
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
