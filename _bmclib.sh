MISSING_DEPS=()

function show_version(){
  version=`cat $thisdir/VERSION-bmc`
  echo
  echo "    bmc v${version}"
  echo "    Bill McCloud's Toolbox"
  echo
  echo "    http://github.com/wearetechnative/bmc"
  echo
  echo "    by Wouter, Pim, et al."
  echo "    © Technative 2024"
  echo
}

function selectEC2Instance() {
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
  echo $instance_id
}



function loadConfig(){
  conffile="${HOME}/.config/bmc/config.env"
  if [ -f ${conffile} ]; then
    source $conffile
  fi
}

getEc2State() {
  local instanceId=$1
  instance_state=$(aws ec2 describe-instances --instance-ids ${instanceID} --query 'Reservations[].Instances[].State.Name' --output text)
  echo "${instance_state}"
}

isEc2HibernationEnabled() {
  local instanceId=$1
  result=$(aws ec2 describe-instances --instance-ids ${instanceID} --query 'Reservations[].Instances[].HibernationOptions.Configured' --output json | jq '.[0]')
	if [ $result == "true" ]; then 
		return 0
	else
		return 1
	fi 
}

userConfirm() {
	local prompt=$1
  gum confirm "${prompt}"
   echo $?
}

ec2Action() {
    local instanceId=$1
    local action=$2
    local targetState=$3
    local timeout=120 # Timeout in seconds
    local interval=5  # Interval between checks

    echo "Performing action '$action' on EC2 instance: $instanceId. Target state: $targetState"
    local startTime=$(date +%s)

    # Placeholder: Perform the action (e.g., AWS CLI)
    case $action in
        start)
            echo "Sending start command for $instanceId"
            ;;
        stop)
            echo "Sending stop command for $instanceId"
            ;;
        hibernate)
            echo "Sending hibernate command for $instanceId"
            ;;
        *)
            echo "Unknown action: $action"
            return 1
            ;;
    esac

    # Wait until the desired state is reached
    while true; do
        local currentState
        currentState=$(getEc2State "$instanceId")
        echo "Current state of $instanceId: $currentState"

        if [[ $currentState == "$targetState" ]]; then
            echo "Target state '$targetState' reached for EC2 instance: $instanceId"
            return 0
        fi

        # Check if the timeout has been exceeded
        local currentTime=$(date +%s)
        if (( currentTime - startTime >= timeout )); then
            echo "Error: Target state '$targetState' not reached within $timeout seconds for EC2 instance: $instanceId"
            return 1
        fi

        sleep $interval
    done
}

# Main logic for managing an EC2 instance
manageEc2Instance() {
    local instanceId=$1

    echo "Managing EC2 instance: $instanceId"

    # Get the current state of the EC2 instance
    local originalState
    originalState=$(getEc2State "$instanceId")
    echo "Original state of EC2 instance $instanceId: $originalState"

    if [[ $originalState == "stopped" ]]; then
        # If the instance is stopped, start it
        ec2Action "$instanceId" "start" "running"
    elif [[ $originalState == "running" ]]; then
        # If the instance is running, check for hibernation
        if isEc2HibernationEnabled "$instanceId"; then
            if userConfirm "Do you want to hibernate the instance?"; then
                ec2Action "$instanceId" "hibernate" "stopped"
            elif userConfirm "Do you want to stop the instance?"; then
                ec2Action "$instanceId" "stop" "stopped"
            fi
        else
            if userConfirm "Do you want to stop the instance?"; then
                ec2Action "$instanceId" "stop" "stopped"
            fi
        fi
    else
        echo "Unknown state: $originalState"
        exit 1
    fi
}

# Run the script with an instance-id as an argument
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <instance-id>"
    exit 1
fi

function checkdeps(){
  if ! command -v $1 &> /dev/null
  then
    MISSING_DEPS+=("$1")
    echo "<$1> could not be found"
    echo "  install this program first"
  fi
}

function deps_missing(){
  if [ ${#MISSING_DEPS[@]} -gt 0 ]
  then
    exit 1
  fi
}

function checkOS {
  if [ -f /etc/lsb-release ]; then
    osType="linux"
  elif [ -f /System/Library/CoreServices/SystemVersion.plist ]; then
    osType="macos"
  else
    osType="other"
  fi
}

function setDates {
  unset currentMFASessionExpirationDate
  expiration=$(sed -n -e "/\[$sourceProfile\]/,/^$/ s/^[[:space:]]*expiration[[:space:]]*=[[:space:]]*\(.*\)/\1/p" "$HOME/.aws/credentials")
  if [[ -z ${expiration} ]]; then expiration="1970-01-01 01:00:00"; fi


  if [[ ${osType} == "macos" ]]; then
    currentMFASessionExpirationDate=$(date -j -f "%Y-%m-%d %H:%M:%S" "${expiration}" "+%s" 2>/dev/null)
    dateCmd="date -j -f "
  elif [[ ${osType} == "linux" ]]; then
    currentMFASessionExpirationDate=$(date -d "$expiration" +%s 2>/dev/null)
  else
    currentMFASessionExpirationDate="0"
  fi
  date_now=$(date +%s)
}

function convertTime() {
  local input_time=$1

  if [[ $input_time =~ ^[0-9]+$ ]]; then
    if [[ $(uname) == "Darwin" ]]; then
      date -j -f "%s" $input_time +"%Y-%m-%d %H:%M:%S"
    else
      date -d @$input_time +"%Y-%m-%d %H:%M:%S"
    fi
  else
    if [[ $(uname) == "Darwin" ]]; then
      date -j -f "%Y-%m-%d %H:%M:%S" "$input_time" +"%s"
    else
      date -d "$input_time" +"%s"
    fi
  fi
}

function printAWSProfiles {
  jsonify-aws-dotfiles | jq -r '
  .config | to_entries |
  map({profile: .key, group: .value.group, arn_number: (.value.role_arn // "" | capture("arn:aws:iam::(?<number>\\d+):").number // "")}) |
  group_by(.group) |
  map({group: .[0].group, profiles: map({profile: .profile, arn_number: .arn_number, group: .group})}) |
  .[] |
  .profiles | map("\(.group)\t\(.profile)\t\(.arn_number)") |
  join("\n")
' | awk 'BEGIN {print "Group\tName\tARN number"} {print}' | column -t -s $'\t'
}

function selectAWSProfile {
  awsProfileGroups=$(jsonify-aws-dotfiles | jq -r '[.config[].group] | unique | sort | .[]' | grep -v null | gum choose --height 25)
  selectedProfile=$(jsonify-aws-dotfiles | jq -r --arg group "$awsProfileGroups" '.config | to_entries | map(select(.value.group == $group)) | (["AWS ACCOUNT", "ROLE"] | @csv), (.[] | [.key, .value.role_arn] | @csv)' | gum table -w 40,120 --height 30)
  selectedProfileName=$(echo "${selectedProfile}" | awk -F "," '{print $1}')
  selectedProfileARN=$(echo "${selectedProfile}" | awk -F "," '{print $2}')
  selectedProfileAccountID=$(echo "${selectedProfileARN}" | awk -F ":" '{print $5}')

  #  if ! expr "${selectedProfileAccountID}" + 0 &>/dev/null; then echo "Error determing AccountID from ARN" ; fi

  sourceProfile=$(jsonify-aws-dotfiles | jq -r --arg arn "$selectedProfileARN" ' .config | to_entries | map(select(.value.role_arn == $arn)) | .[0].value.source_profile // "Error" ')

  if [[ ${sourceProfile} == "Error" ]]; then sourceProfile=${selectedProfileName}; fi
}

function setMFA {
  checkOS
  setDates
  echo
  echo "MFA: ${mfa}"
  if [[  ${mfa} == "true" ]]; then
    awsMFADevice=$(awk -v profile="${sourceProfile}-long-term" ' $0 == "[" profile "]" {found=1; next} /^\[.*\]/ {found=0} found && /^aws_mfa_device/ {print $3; exit} ' ~/.aws/credentials)
    if [[ -z ${currentMFASessionExpirationDate} ]]; then expiration="1" ;fi
    if [[ ${currentMFASessionExpirationDate} -lt ${date_now} ]]; then
      if [[ ! -z ${awsMFADevice} ]]; then
        echo aws-mfa --profile ${sourceProfile} --force --device ${awsMFADevice}
        if [[ ! -z $totpScript ]]; then
          totpCode=$(${totpScript})
          echo ${totpCode} |  ${clipboardCommand}
          echo "-- Copied to clipboard";
          echo "${totpCode}"
        else
          echo "Code: ${totpCode}"
        fi
        aws-mfa --profile ${sourceProfile} --force --device ${awsMFADevice}
        if [[ $? -ne 0 ]]; then echo "!!  Error with AWS MFA code for device. Wrong TOPT?"; return;fi
      else
        echo "!! awsMFADevice not found. Can't renew session"
        echo
      fi
    else
      echo "Current MFA Session Valid, until: $(convertTime ${currentMFASessionExpirationDate})"
      echo
    fi
  fi
}

