#!/usr/bin/env bash
# Version: 2024020601

rprompt_config="true"
aws_sso="false"
aws_mfa="false"

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
  set_prompt

  if [[ ${aws_mfa} == "true" ]]; then mfa; fi
  # if [[ ${aws_sso} == "true" ]]; then sso; fi

}

# read_aws_config() {
#   oldCol=$COLUMNS
#   COLUMNS=6
#   config_file="$HOME/.aws/config"

#   if [ -f "$config_file" ]; then
#     while IFS= read -r line; do
#       if [[ $line == \[* ]]; then
#         current_profile=$(echo "$line" | sed 's/\[\(.*\)\]/\1/')
#       elif [[ $line == role_arn* ]]; then
#         role_arn=$(echo "$line" | awk -F'=' '{print $2}' | tr -d '[:space:]')
#         account_number=$(echo "$role_arn" | cut -d':' -f5)

#         if [[ "$account_number" =~ ^[0-9]+$ ]]; then
#           printf "%s : %s\n" "$account_number" "$current_profile"
#         fi
#       fi
#     done <"$config_file"
#   else
#     echo "AWS config file not found: $config_file"
#   fi
#   COLUMNS=$oldCol
# }

function usage {
  echo "Usage: aps [-n|--no-sdk] [-h|--help]"
  echo "For normal usage, just run aps and make your selection followed by Enter."
  echo "If you do not want the AWS_SDK_LOAD_CONFIG environment variable to be set to true, append -n or --no-sdk to the command"
  # exit 1
}

function mfa {
  if [[ ${exit_script} == "true" ]]; then return; fi
  # Check for valid aws-mfa session
  unset expiration_date

  # search long-term
  # echo "!!! MFA: AWS_PROFILE: ${AWS_PROFILE}"
  source_profile=$(sed -n -e "/\[.*${AWS_PROFILE}\]/,/^$/ s/^[[:space:]]*source_profile[[:space:]]*=[[:space:]]*\(.*\)/\1/p" ${HOME}/.aws/config)
  # echo "!!! MFA: source_profile: ${source_profile}"

  if [[ -z ${source_profile} ]]; then
    source_profile_longterm="${AWS_PROFILE}-long-term"
    source_profile=${AWS_PROFILE}
  else
    source_profile_longterm="${source_profile}-long-term"
  fi

  expiration_date=$(date -j -f "%Y-%m-%d %H:%M:%S" "$(sed -n -e "/\[${source_profile}\]/,/^$/ s/^[[:space:]]*expiration[[:space:]]*=[[:space:]]*\(.*\)/\1/p" ${HOME}/.aws/credentials)" "+%s" 2>/dev/null)
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
    echo "MFA Validi, until: $(date -j -f "%s" ${expiration_date} "%Y-%m-%d %H:%M:%S")"
  fi
}

# function selection_menu {
#   # echo ${profiles[*]}
#   # echo "-: Unset Profile"
#   # for ((i = 0; i < $profiles_len; i++)); do
#   #   echo "$i: ${profiles[$i]}"
#   # done
#   # read_selection
# }

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

function set_prompt {
  new_prompt="${cmd_prompt}aps:(${AWS_PROFILE}): "

  # set prompt for ZSH or BASH
  if [[ $shell_type == "zsh" ]]; then
    # ZSH
    if [[ ${rprompt_config} == "true" ]]; then
      # export RPROMPT=${AWS_PROFILE}-${TF_BACKEND}
      export RPROMPT=${AWS_PROFILE}-$(basename ${TF_BACKEND})
    else
      export PROMPT="$new_prompt"
    fi
  else
    #BASH
    export PS1="$new_prompt"
  fi
}

function list_config_profiles {
  config_profiles=$(grep -E '\[profile .+\]' ~/.aws/config | sed 's/\[profile \(.*\)\]/\1/')
  max_account_length=0
  max_profile_length=0
  max_region_length=0

  # Bepaal de maximale lengte van elk veld
  for config_profile in ${config_profiles}; do
    role_arn=$(grep -A3 "\[profile ${config_profile}\]" ~/.aws/config | grep role_arn | awk -F' = ' '{print $2}')
    account_number=$(echo ${role_arn} | awk -F'::' '{print $2}' | awk -F':' '{print $1}')
    region=$(grep -A3 "\[profile ${config_profile}\]" ~/.aws/config | grep region | awk -F' = ' '{print $2}')
    if [[ ! -z ${role_arn} ]]; then
      if [ ${#account_number} -gt $max_account_length ]; then
        max_account_length=${#account_number}
      fi
      if [ ${#config_profile} -gt $max_profile_length ]; then
        max_profile_length=${#config_profile}
      fi
      if [ ${#region} -gt $max_region_length ]; then
        max_region_length=${#region}
      fi
    fi
  done

  # Uitlijning en weergave van gegevens
  for config_profile in ${config_profiles}; do
    role_arn=$(grep -A3 "\[profile ${config_profile}\]" ~/.aws/config | grep role_arn | awk -F' = ' '{print $2}')
    account_number=$(echo ${role_arn} | awk -F'::' '{print $2}' | awk -F':' '{print $1}')
    region=$(grep -A3 "\[profile ${config_profile}\]" ~/.aws/config | grep region | awk -F' = ' '{print $2}')
    if [[ ! -z ${role_arn} ]]; then
      printf "%-${max_account_length}s : %-${max_profile_length}s : %-${max_region_length}s\n" "${account_number}" "${config_profile}" "${region}"
    fi
  done
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
