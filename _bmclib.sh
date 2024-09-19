function loadConfig(){
  conffile="${HOME}/.config/bmc/config.env"
  if [ -f ${conffile} ]; then
    source $conffile
  fi
}

function checkdeps(){
  if ! command -v $1 &> /dev/null
  then
    echo "<$1> could not be found"
    echo "  install this program first"
    ## exit 1
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
      fi
    else
      echo "Current MFA Session Valid, until: $(convertTime ${currentMFASessionExpirationDate})"
    fi
  fi
}

