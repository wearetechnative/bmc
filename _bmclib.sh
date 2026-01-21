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

function selectProfileGroup(){
	local awsProfileGroup=$(jsonify-aws-dotfiles | jq -r '[.config[].group] | unique | sort | .[]' | grep -v null | gum choose --height 25)
  local aws_profiles=($(jsonify-aws-dotfiles | jq -r ".config | to_entries[] | select(.value.group == \"${awsProfileGroup}\") | .key"))


	if [ ${#aws_profiles[@]} -eq 0 ]; then
		echo "No AWS-profiles found for group '$awsProfileGroup'."
		# exit 1
	fi

	echo "Starting EC2 listing for profiles..."
	for profile in "${aws_profiles[@]}"; do
  	echo -e "\n------  $profile  ------"

    AWS_PROFILE="$profile" ec2ListInstances

		if [[ $? -ne 0 ]]; then
			echo "[ERROR] Command failed for profile: $profile"
		fi

	done

	echo -e "All profiles processed."

}


function ec2FindInstance(){
	local searchString=$1
	if [[ -z ${searchString} ]]; then
		echo "!!! No search string"
    echo "Usage: $(basename $0) ec2find <search-string>"
		exit 0
	fi

	outputFile=$(mktemp)
	echo "Searching for: ${searchString}"
	selectProfileGroup > $outputFile
	output="InstanceId,PrivateIpAddress,PublicIpAddress,State,Hibernate,Name,Profile\n"
	while IFS= read -r line
	do
		if [[ ${line} == **---** ]]; then
			profileHit=$(echo "${line}" | sed 's/^-*//; s/-*$//; s/  */ /g')
		fi
		if [[ ${line} == **${searchString}** ]] ; then
			searchHit="true"
			hitstring=$(echo ${line}|sed 's/│/,/g')
			while read -r items; do
				# Zet de waarden in de juiste CSV-indeling
				instance_id=$(echo $items | awk -F "," '{print $2}')
				private_ip=$(echo $items | awk -F "," '{print $3}')
				public_ip=$(echo $items | awk -F "," '{print $4}')
				state=$(echo $items | awk -F "," '{print $5}')
				name=$(echo $items | awk -F "," '{print $7}')
				hibernation_status=$(echo $items | awk -F "," '{print $6}')
				profile=${profileHit}
				output+="$instance_id,$private_ip,$public_ip,$state,$hibernation_status,$name,$profile\n"
			done <<< ${hitstring}
		fi
	done < $outputFile

	if [[ -z ${searchHit} ]]; then
		echo "!! String not found in AWS PROFILE Group: ${awsProfileGroup}"
	fi

	rm $outputFile

	echo -e "$output" | gum table -p

}

function ec2ListInstances(){

instances=$(aws ec2 describe-instances \
    --query 'Reservations[*].Instances[*].[InstanceId, PrivateIpAddress, PublicIpAddress, State.Name, Tags[?Key==`Name`].Value | [0], HibernationOptions.Configured]' \
    --output text)

# Zet de tekst output om naar CSV-formaat
output="InstanceId,PrivateIpAddress,PublicIpAddress,State,Hibernate,Name\n"
while read -r line; do
    # Zet de waarden in de juiste CSV-indeling
    instance_id=$(echo $line | awk '{print $1}')
    private_ip=$(echo $line | awk '{print $2}')
    public_ip=$(echo $line | awk '{print $3}')
    state=$(echo $line | awk '{print $4}')
    name=$(echo $line | awk '{print $5}')
    hibernation_status=$(echo $line | awk '{print $6}')

    # Voeg de hibernation status vóór de Name toe aan de output
    output+="$instance_id,$private_ip,$public_ip,$state,$hibernation_status,$name\n"
done <<< "$instances"

# Gebruik gum table om de CSV in een mooie tabel weer te geven
echo -e "$output" | gum table -p
}



function ec2SelectInstance(){
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
    instance_id=$(echo -e "$formatted_instances" | gum table -w 20,16,16,8,50 --height 20 | awk -F, '{print $1}')
    echo $instance_id
}

function ec2CheckInstanceStatus(){
  local instance_id=$1
  ec2instancestate=$(aws ec2 describe-instances --instance-ids ${instance_id} --query 'Reservations[].Instances[].State.Name' --output text)
  echo ${ec2instancestate}
}

function ec2CheckHibernationState() {
	local instance_id=$1
	ec2hibernationenabled=$(aws ec2 describe-instances --instance-ids ${instance_id} --query 'Reservations[].Instances[].HibernationOptions.Configured' --output json | jq '.[0]')
	case ${ec2hibernationenabled} in
		true)
			return 0
			;;
		*)
			return 1
			;;
	esac
}

function ec2CheckNewInstanceState(){
	local  instance_id=$1
	local  desiredinstancestate=$2
	local  elapsed=0
	local  interval=5
	local  timeout=300

	while [[ ${elapsed} -lt ${timeout}  && ${currentinstancestate} != ${desiredinstancestate} ]] ; do
		((elapsed += interval))
    currentinstancestate=$(aws ec2 describe-instances --instance-ids ${instance_id} --query 'Reservations[].Instances[].State.Name' --output text)

		sleep ${interval}
	done

	if [[ ${currentinstancestate} == ${desiredinstancestate} ]]; then
		echo "Instance ${instance_id} has reached new state ${currentinstancestate} in ${elapsed} seconds.";
	else
		echo "Instance ${instance_id} has not reached desired state ${currentinstancestate} within ${timeout} seconds.";
	fi
	exit 0
}



function ec2StopStartInstance(){
	local instance_id=$(ec2SelectInstance)
	local instance_state=$(ec2CheckInstanceStatus ${instance_id})
  local bmcpath=$(dirname $0)
	case ${instance_state} in
		stopped)
			aws ec2 start-instances --instance-ids ${instance_id} >/dev/null
      answer=$(gum spin --spinner meter --title "Starting instance ${instance_id}" -- bash -c "source ${bmcpath}/_bmclib.sh && ec2CheckNewInstanceState ${instance_id} running")
			;;
		running)
			local stoppingoptions="stop\nexit menu"
			if $(ec2CheckHibernationState ${instance_id}); then
				local stoppingoptions="hibernate\n${stoppingoptions}"
			fi
			stoppingmethod=$(echo -e $stoppingoptions | gum choose --header "Choose stop-method for instance: ${instance_id}")

			case ${stoppingmethod} in
				hibernate)
					aws ec2 stop-instances --instance-ids ${instance_id} --hibernate  >/dev/null
					answer=$(gum spin --spinner meter --title "Hibernating instance ${instance_id}" -- bash -c "ec2CheckNewInstanceState ${instance_id} stopped")
					;;
				stop)
					aws ec2 stop-instances --instance-ids ${instance_id}  >/dev/null
					answer=$(gum spin --spinner meter --title "Stopping instance ${instance_id}" -- bash -c "source ${bmcpath}/_bmclib.sh && ec2CheckNewInstanceState ${instance_id} stopped")
					;;
							esac
			;;
	*)
					echo -e "Instance ${instance_id} not in running/stopped state.\nCurrent state: ${instance_state}\n"
					exit 0
					;;
esac

 if [[ ! -z ${answer} ]]; then
	echo "ANSWER: $answer"
fi

}


function loadConfig(){
  conffile="${HOME}/.config/bmc/config.env"
  if [ -f ${conffile} ]; then
    source $conffile
  fi
}

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
    #dateCmd="date -j -f "
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

function useOrSelectAWSProfile {
  if [[ -z $AWS_PROFILE ]]; then
    selectAWSProfile "$@"
    setMFA
    export AWS_PROFILE=$selectedProfileName
  fi
}

function selectAWSProfile {

  if [[ -z $preferedProfile ]]; then
    awsProfileGroups=$(jsonify-aws-dotfiles | jq -r '[.config[].group] | unique | sort | .[]' | grep -v null | gum filter --height 25)

    # Check if profile group selection was cancelled
    if [[ -z $awsProfileGroups ]]; then
      unset selectedProfileName
      return
    fi

    #OLD TABLE SELECTOR
    #selectedProfile=$(jsonify-aws-dotfiles | jq -r --arg group "$awsProfileGroups" '.config | to_entries | map(select(.value.group == $group)) | (["AWS ACCOUNT", "ROLE"] | @csv), (.[] | [.key, .value.role_arn] | @csv)' | gum table -w 40,120 --height 30)

    selectedProfileTable=$(jsonify-aws-dotfiles | jq -r --arg group "$awsProfileGroups" '.config | to_entries | map(select(.value.group == $group)) | (["AWS ACCOUNT", "ROLE"] | @csv), (.[] | [.key, .value.role_arn] | @csv)' | sed -e s/\"//g | sed -e s/ROLE/ACCOUNT\ ID,ROLE/ | sed -e s/arn:aws:iam::// | sed -e s/:role\\//,/ | column -s, -t)
    header=$(echo "  $selectedProfileTable" | head -n1)
    selectedProfile=$(echo "$selectedProfileTable" |tail -n +2 | gum filter --header="$header")

    # Check if profile selection was cancelled
    if [[ -z $selectedProfile ]]; then
      unset selectedProfileName
      return
    fi

    # Split the selectedProfile string in a way that works in both bash and zsh
    local profile_name=$(echo "$selectedProfile" | awk '{print $1}')
    local account_id=$(echo "$selectedProfile" | awk '{print $2}')
    local role_name=$(echo "$selectedProfile" | awk '{print $3}')

    selectedProfile="${profile_name},arn:aws:iam::${account_id}:role/${role_name}"
    selectedProfileARN=$(echo "${selectedProfile}" | awk -F "," '{print $2}')
  else
    selectedProfileARN=$(jsonify-aws-dotfiles| jq -r ".config.\"${preferedProfile}\".role_arn")
    selectedProfile="$preferedProfile,$selectedProfileARN"
  fi

  selectedProfileName=$(echo "${selectedProfile}" | awk -F "," '{print $1}')

  sourceProfile=$(jsonify-aws-dotfiles | jq -r --arg arn "$selectedProfileARN" ' .config | to_entries | map(select(.value.role_arn == $arn)) | .[0].value.source_profile // "Error" ')

  if [[ ${sourceProfile} == "Error" ]]; then
    inCredentials=$(jsonify-aws-dotfiles | jq -r ".credentials.\"$selectedProfileName\"")
    if [ "$inCredentials" = 'null' ]; then
      unset sourceProfile
    else
      sourceProfile=${selectedProfileName}
    fi
  fi

  unset preferedProfile
}

function setMFA {
  if [[ -z $sourceProfile ]]; then
    echo "Error could not set MFA without valid sourceProfile"
    return 1
  else
    echo "-- Using AWS source-profile: $sourceProfile"
  fi

  checkOS
  setDates
  echo
  if [[  ${mfa} == "true" ]]; then
    awsMFADevice=$(awk -v profile="${sourceProfile}-long-term" ' $0 == "[" profile "]" {found=1; next} /^\[.*\]/ {found=0} found && /^aws_mfa_device/ {print $3; exit} ' ~/.aws/credentials)
    if [[ -z ${currentMFASessionExpirationDate} ]]; then expiration="1" ;fi
    if [[ ${currentMFASessionExpirationDate} -lt ${date_now} ]]; then
      if [[ ! -z ${awsMFADevice} ]]; then
        echo "-- Refreshing MFA session for ${sourceProfile}..."
        if [[ ! -z $totpScript ]]; then
          echo "-- Executing TOTP script..."
          totpCode=$("${totpScript[@]}")
          echo "${totpCode}"
          if [[ ! -z $clipboardCopyCommand ]]; then
            if echo ${totpCode} | "${clipboardCopyCommand[@]}" 2>/dev/null; then
              echo "-- Copied to clipboard"
            else
              echo "-- Note: Clipboard copy failed (command not found or error)"
            fi
          fi
        else
          echo "-- No TOTP script configured. Please enter MFA code manually."
        fi
        aws-mfa --profile ${sourceProfile} --force --device ${awsMFADevice}
        if [[ $? -ne 0 ]]; then echo "!!  Error with AWS MFA code for device. Wrong TOPT?"; return;fi
      else
        echo "!! AWS MFA Device not found. Can't renew session"
        echo
      fi
    else
      echo "Current MFA Session Valid, until: $(convertTime ${currentMFASessionExpirationDate})"
      echo
    fi
  fi
}
