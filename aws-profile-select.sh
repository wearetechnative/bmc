## Copyright 2022 Jesse Price

## User configurable setings

# User RPROMPT function, zsh only
rprompt_config="true"
aws_sso="true"
aws_mfa="true"
mfa_profile="technative-userauth"
mfa_device_arn="arn:aws:iam::521402697040:mfa/wouter@technative.nl"

# Enable setting of AWS_SDK_LOAD_CONFIG by default
sdk=1

AWS_PROFILE="technative-userauth"

# # Voer AWS CLI-commando uit om de identiteit op te halen
# aws_sts_output=$(aws sts get-caller-identity --profile $AWS_PROFILE 2>&1)

aws_mfa_session_file=".aws-profile-switch.mfa.session"
aws_mfa_session_file=${aws_mfa_session_file/#\~/$HOME}

# Check for valid aws-mfa session
#while [[ -z ${mfa_session_valid} ]]; do
#  while [[ ! -f ${aws_mfa_session_file} ]]; do
#    aws-mfa --profile ${mfa_profile} --device ${mfa_device_arn} | sed -n -e 's/^.*expire at //p' | tee ${aws_mfa_session_file}
#  done
#  check_datetime=$(<${aws_mfa_session_file})
#  current_datetime=$(date "+%Y-%m-%d %H:%M")
#
#  if [[ ${check_datetime} == ${current_datetime} ]]; then
#    mfa_session_valid=true
#  else
#    aws-mfa --profile ${mfa_profile} --device ${mfa_device_arn} | sed -n -e 's/^.*expire at //p' | tee ${aws_mfa_session_file}
#  fi
#done

# Controleer de uitvoer van het AWS CLI-commando
if [ $? -eq 0 ]; then
  echo "Valid AWS-session found:"
  echo "$aws_sts_output"
else
  echo "No valid AWS-session found."
  aws-mfa --profile ${mfa_profile} -userauth --device ${mfa_device_arn}
fi

if [ -n "$ZSH_VERSION" ]; then
  # zsh-handling
  shell_type=zsh
  echo "!! check zsh_version"
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
profiles=($profiles)
IFS=$IFSBAK

profiles_len=${#profiles[*]}

function main {
  printf "Current value of AWS_SDK_LOAD_CONFIG: ${AWS_SDK_LOAD_CONFIG}\n"
  echo
  echo ------------- AWS Profile Select-O-Matic -------------
  echo
  if [ -z "$AWS_PROFILE" ]; then
    printf "No profile set yet\n\n"
  else
    printf "\nCurrently-selected profile: ${AWS_PROFILE}\n\n"
  fi
  echo "Type the number of the profile you want to use from the list below, and press enter"
  echo

  # Show the menu
  selection_menu
}

function usage {
  echo "Usage: aps [-n|--no-sdk] [-h|--help]"
  echo "For normal usage, just run aps and make your selection followed by Enter."
  echo "If you do not want the AWS_SDK_LOAD_CONFIG environment variable to be set to true, append -n or --no-sdk to the command"
  # exit 1
}

# Takes 1 argument: the variable name which stores the prompt (PS1/PROMPT)
function set_prompt {
  # echo ${profiles[*]}
  echo "-: Unset Profile"
  for ((i = 0; i < $profiles_len; i++)); do
    echo "$i: ${profiles[$i]}"
  done
  read_selection
}

function selection_menu {
  # echo ${profiles[*]}
  echo "-: Unset Profile"
  for ((i = 0; i < $profiles_len; i++)); do
    echo "$i: ${profiles[$i]}"
  done
  read_selection
}

function read_selection {
  echo
  printf 'Selection: '
  read choice
  case $choice in
  '' | *[!0-9\-]*)
    clear
    echo Invalid selection. Make a valid selection from the list above or press ctrl+c to exit
    echo '-> Error: Not a number, and not "-"'
    echo
    selection_menu
    ;;
  esac
  in_range=false
  while [ $in_range != true ]; do
    if [[ $choice == '-' ]]; then
      echo "Deactivating all profiles"
      unset AWS_PROFILE
      if [[ $shell_type == "zsh" ]]; then
        export PROMPT="$PROMPTBAK"
        if [[ ${rprompt_config} == "true" ]]; then
          unset RPROMPT
        fi
      else
        export PS1="$PS1BAK"
      fi
      in_range=true
    elif (($choice >= 0)) && (($choice <= ${profiles_len})); then
      # Set AWS_SDK_LOAD_CONFIG to true to make this useful for tools such as Terraform and Serverless framework
      if (($sdk == 1)); then
        export AWS_SDK_LOAD_CONFIG=1
      else
        export AWS_SDK_LOAD_CONFIG=0
      fi
      echo "Activating profile ${choice}: ${profiles[choice]}"
      export AWS_PROFILE="${profiles[choice]}"
      new_prompt="${cmd_prompt}aps:(${profiles[choice]}): "
      if [[ $shell_type == "zsh" ]]; then
        if [[ ${rprompt_config} == "true" ]]; then
		    # export RPROMPT=${AWS_PROFILE}-${TF_BACKEND}
        export RPROMPT=${AWS_PROFILE}-${TF_BACKEND}
        else
          export PROMPT="$new_prompt"
        fi
      else
        export PS1="$new_prompt"
      fi
      in_range=true
    else
      clear
      echo Invalid selection. Select a valid profile number or press ctrl+c to exit
      echo "-> Error:  Number must be one of 0-"$((${#profiles[@]} - 1))""
      echo
      selection_menu
    fi
  done
}

if [ $# -gt 0 ]; then
  while [ ! $# -eq 0 ]; do
    case "$1" in
    --help | -h)
      usage
      ;;
    --no-sdk | -n)
      sdk=0
      # Kick off the main function:
      main
      ;;
    esac
    shift
  done
else
  # Kick off the main function:
  main
fi
