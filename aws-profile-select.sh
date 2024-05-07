#!/usr/bin/env bash
# Version: 202404102
# bla

rprompt_config="true"
aws_sso="false"
aws_mfa="false"

mkdir -p ~/.config/aws-profile-select/

if [[ -f ~/.config/aws-profile-select/env ]]; then
  source ~/.config/aws-profile-select/env
  echo "MFA enabled"
fi

# Enable setting of AWS_SDK_LOAD_CONFIG by default
sdk=1
exit_script="false"
unset mfaduration

if [ -n "$ZSH_VERSION" ]; then
  # zsh-handling
  shell_type=zsh
  # echo "!! check zsh_version"
  setopt ksh_arrays
  setopt SH_WORD_SPLIT
  if [ -z "$PROMPTBAK" ]; then
    export PROMPTBAK="$PROMPT"
    cmd_prompt="$PROMPT"
  else
    cmd_prompt="$PROMPTBAK"
  fi
else
  # bash-handling
  # Backup prompt as a separate variable if not already backed up
  if [ -z "$PS1BAK" ]; then
    export PS1BAK="$PS1"
    cmd_prompt="$PS1"
  else
    cmd_prompt="$PS1BAK"
  fi
fi

profiles=$(grep '^\[' <~/.aws/config | sed -E 's/\[profile (.*)/\1/g' | sed 's/\[//; s/\]//')

IFSBAK=$IFS
IFS=$'\n'
profiles=(${profiles})
IFS=$IFSBAK

profiles_len=${#profiles[*]}

function main {
  # parse_arguments
  printf "Current value of AWS_SDK_LOAD_CONFIG: ${AWS_SDK_LOAD_CONFIG}\n"
  echo ------------- AWS Profile Select-O-Matic -------------
  if [ -z "$AWS_PROFILE" ]; then
    printf "No profile set yet\n\n"
  else
    printf "\nCurrently-selected profile: ${AWS_PROFILE}\n\n"
  fi
  echo "Type the number of the profile you want to use from the list below, and press enter"
  echo

  # Show the menu
  # echo "!! READ SELECTION"
  read_selection
  # echo "!! CHECK DSK"
  check_sdk
  # echo "!! SET PROMPT"

  if [[ ${aws_mfa} == "true" ]]; then mfa; fi
  # if [[ ${aws_sso} == "true" ]]; then sso; fi

}

function usage {
  echo "Usage: aps [-n|--no-sdk] [-h|--help]"
  echo "For normal usage, just run aps and make your selection followed by Enter."
  echo "If you do not want the AWS_SDK_LOAD_CONFIG environment variable to be set to true, append -n or --no-sdk to the command"
  # exit 1
}

function checkOS {
  if [ -f /etc/lsb-release ]; then
    echo "linux"
  elif [ -f /System/Library/CoreServices/SystemVersion.plist ]; then
    echo "macos"
  else
    echo "other"
  fi

}

function mfa {
  if [[ ${exit_script} == "true" ]]; then return; fi
  # Check for valid aws-mfa session
  unset expiration_date

  # search long-term
  # echo "!!! MFA: AWS_PROFILE: ${AWS_PROFILE}"

  # Detecting source_profile to obtain mfa-device
  source_profile=$(sed -n -e "/\[.*${AWS_PROFILE}\]/,/^$/ s/^[[:space:]]*source_profile[[:space:]]*=[[:space:]]*\(.*\)/\1/p" ${HOME}/.aws/config)
  # echo "!!! MFA: source_profile: ${source_profile}"

  if [[ -z ${source_profile} ]]; then
    source_profile_longterm="${AWS_PROFILE}-long-term"
    source_profile=${AWS_PROFILE}
  else
    source_profile_longterm="${source_profile}-long-term"
  fi

  expiration=$(sed -n -e "/\[$source_profile\]/,/^$/ s/^[[:space:]]*expiration[[:space:]]*=[[:space:]]*\(.*\)/\1/p" "$HOME/.aws/credentials")

  osType=$(checkOS)
  if [[ ${osType} == "macos" ]]; then
    expiration_date=$(date -j -f "%Y-%m-%d %H:%M:%S" "${expiration}" "+%s" 2>/dev/null)
    dateCmd="date -j -f "
  elif [[ ${osType} == "linux" ]]; then
    expiration_date=$(date -d "$expiration" +%s 2>/dev/null)
  else
    echo "!! ostype not linux or macos"
  fi

  #  expiration_date=$(${dateCmd}  "%Y-%m-%d %H:%M:%S" "$(sed -n -e "/\[${source_profile}\]/,/^$/ s/^[[:space:]]*expiration[[:space:]]*=[[:space:]]*\(.*\)/\1/p" ${HOME}/.aws/credentials)" "+%s" 2>/dev/null)
  date_now=$(date +%s)
  mfa_arn=$(sed -n -e "/\[${source_profile_longterm}\]/,/^$/ s/^[[:space:]]*aws_mfa_device[[:space:]]*=[[:space:]]*\(.*\)/\1/p" ${HOME}/.aws/credentials)

  if [[ ${expiration_date} -lt ${date_now} ]]; then
    if [[ ! -z ${mfa_arn} ]]; then
      echo aws-mfa --profile ${source_profile} --force --device ${mfa_arn}
      aws-mfa --profile ${source_profile} --force --device ${mfa_arn}
    else
      echo "!! MFA_arn not found. Can't renew session"
    fi
  else
    echo "MFA Valid, until: ${expiration}"
  fi
}

function read_selection {
  PS3="Select a profile number: "
  echo "-: Unset profile"
  select item in "${profiles[@]}"; do
    if [[ -n ${item} ]]; then
      echo "AWS_PROFILE set to $item"
      export AWS_PROFILE=${item}
      break
    fi
    if [[ ${REPLY} == "-" ]]; then
      echo "Unsetting profile"
      unset AWS_PROFILE
      if [[ $shell_type == "zsh" ]]; then
        export PROMPT="$PROMPTBAK"
        if [[ ${rprompt_config} == "true" ]]; then
          unset RPROMPT
        fi
      else
        export PS1="$PS1BAK"
      fi
      exit_script=true
      break
    else
      echo Invalid selection. Select a valid profile number or press ctrl+c to exit
    fi
  done
}

function check_sdk {
  # Set AWS_SDK_LOAD_CONFIG to true to make this useful for tools such as Terraform and Serverless framework
  if (($sdk == 1)); then
    export AWS_SDK_LOAD_CONFIG=1
  else
    export AWS_SDK_LOAD_CONFIG=0
  fi
}

function list_config_profiles {

  config_profiles=($(grep -E '\[profile .+\]' ~/.aws/config | sed 's/\[profile \(.*\)\]/\1/'))

  max_account_length=0
  max_profile_length=0
  max_region_length=0

  # Bepaal de maximale lengte van elk veld
  for config_profile in "${config_profiles[@]}"; do
    osType="macos"
    if [[ ${osType} == "macos" ]]; then
      profile_data=($(
        profile_name="${config_profile}"
        sed -n -e "/^\[profile $profile_name\]/, /^\[/ {/^\[/d; p;}" ~/.aws/config | awk 'NF'
      ))
    else
	    profile_data=($(profile_name="${config_profile}"; awk -v profile="$config_profile" '$0 ~ "[profile " profile "]"{p=1} p && NF && !/^\[/{print; p=0}' ~/.aws/config))
    fi

    for field in "${profile_data[@]}"; do
      key=$(echo "$field" | cut -d '=' -f 1)
      value=$(echo "$field" | cut -d '=' -f 2-)

      if [ "$key" = "role_arn" ]; then
        account_number=$(echo "$value" | awk -F ':' '{print $5}')
        role_arn="$value"
      elif [ "$key" = "region" ]; then
        region="$value"
      fi
    done

    printf "%-20s | %-15s | %-10s | %s\n" "$config_profile" "$account_number" "$region" "$role_arn"
  done | column -t -s '|'
}

if [ $# -gt 0 ]; then
  while [ ! $# -eq 0 ]; do
    case "$1" in
    -l)
      list_config_profiles
      break
      ;;
    --help | -h)
      usage
      ;;
    --no-sdk | -n)
      sdk=0
      # Kick off the main function:
      main
      ;;
    *)
      echo "Invalid script parameters"
      main
      ;;
    esac
    main
  done
else
  # Kick off the main function:
  main
fi
