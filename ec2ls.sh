#!/usr/bin/env bash



function selectProfileGroup(){
	awsProfileGroup=$(jsonify-aws-dotfiles | jq -r '[.config[].group] | unique | sort | .[]' | grep -v null | gum choose --height 25)
  aws_profiles=($(jsonify-aws-dotfiles | jq -r ".config | to_entries[] | select(.value.group == \"${awsProfileGroup}\") | .key"))
  

	if [ ${#aws_profiles[@]} -eq 0 ]; then
		echo "Geen AWS-profielen gevonden voor groep '$awsProfileGroup'."
		exit 1
	fi

	echo "Starting EC2 listing for profiles..." 
	for profile in "${aws_profiles[@]}"; do
	echo -e "\n------  $profile  ------" 

		#AWS_PROFILE="$profile" bmc ec2ls
    AWS_PROFILE="$profile" bmc ec2ls 

		if [[ $? -ne 0 ]]; then
			echo "[ERROR] Command failed for profile: $profile"
		fi

	done

	echo -e "All profiles processed." 

}

function findSearchPattern() {

  outputFile=$(mktemp)
  echo "Searching for: ${searchString}"
  selectProfileGroup > $outputFile
  while IFS= read -r line
  do
    if [[ ${line} == **---** ]]; then 
      profileHit=$(echo "${line}" | sed 's/^-*//; s/-*$//; s/  */ /g')
    fi
    if [[ ${line} == **${searchString}** ]] ; then 
      searchHit="true"
      echo -e "\n--- String found in profile: ${profileHit}"
      echo " ${line} "
    fi
    done < $outputFile

    if [[ -z ${searchHit} ]]; then
      echo "!! String not found in AWS PROFILE Group: ${awsProfileGroup}"
    fi

    rm $outputFile

}



function ec2ListInstances() {

# Uitvoeren van het AWS CLI-commando om informatie over EC2-instances op te halen
aws_output=$(aws ec2 describe-instances)

# Extract de gewenste velden uit de JSON-output met behulp van jq
instances=$(echo "$aws_output" | jq -r '.Reservations[].Instances[] | "\(.InstanceId) - \(.PrivateIpAddress) - \(.PublicIpAddress) - \(.Tags[] | select(.Key=="Name") | .Value)"')

# Print de gewenste velden
echo "$instances"
}


while getopts 'as:' opt; do
	case "$opt" in
		a)
			selectProfileGroup
			;;
		s)
			searchString="$OPTARG"
			findSearchPattern
			;;
		?)
			echo -e "Invalid command option."
			exit 1
			;;
	esac
done
shift "$(($OPTIND - 1))"

