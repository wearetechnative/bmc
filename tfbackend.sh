#!/bin/bash
backends="$(ls *.tfbackend)"

IFSBAK=$IFS
IFS=$'\n'
backends=($backends)
IFS=$IFSBAK
backends_len=${#backends[*]}
echo $backends
echo $backends_len

sdk=1

function main {
  echo ------------- Select Backend -------------
  # if [ -z "$TF_BACKEND" ]; then
  #   TF_BACKEND=$(terraform show | grep -A2 "module.terraformbackend.module.state_lock.aws_dynamodb_table.this"|tail -1|awk -F ":" '{print $5}')
  # else
  #   printf "\nCurrently-selected backend: ${TF_BACKEND}\n\n"
  # fi
  echo "Type the number of the backend you want to use from the list below, and press enter"
  echo

  # Show the menu
  selection_menu
}

function usage {
  echo "Usage: $(basename $0) [-h|--help]"
  echo "For normal usage, just run aps and make your selection followed by Enter."
  # exit 1
}

function set_tfbackend {
  backend_file=$1
  if [[ -z ${backend_file} ]]; then echo "!!! Error backend-file"; exit 1;fi
  terraform  init -backend-config="${backend_file}" -reconfigure
  TF_BACKEND=$(terraform show | grep -A2 "module.terraformbackend.module.state_lock.aws_dynamodb_table.this"|tail -1|awk -F ":" '{print $5}')
}

function set_prompt {
  echo "-: Unset TF_BACKEND"
  for ((i = 0; i < $backends_len; i++)); do
    echo "$i: ${backends[$i]}"
  done
  read_selection
}

function selection_menu {
  # echo ${profiles[*]}
  echo "-: Unset backend"
  for ((i = 0; i < $backends_len; i++)); do
    echo "$i: ${backends[$i]}"
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
      echo "Deactivating backend"
      unset TF_BACKEND
      if [[ $shell_type == "zsh" ]]; then
        export PROMPT="$PROMPTBAK"
        if [[ ${rprompt_config} == "true" ]]; then
          unset RPROMPT
        fi
      else
        export PS1="$PS1BAK"
      fi
      in_range=true
    elif (($choice >= 0)) && (($choice <= ${backends_len})); then
      # Set AWS_SDK_LOAD_CONFIG to true to make this useful for tools such as Terraform and Serverless framework

      if (($sdk == 1)); then
        export AWS_SDK_LOAD_CONFIG=1
      else
        export AWS_SDK_LOAD_CONFIG=0
      fi
      echo "Activating backend ${choice}: ${backends[choice]}"
      export TF_BACKEND="${backends[choice]}"
      set_tfbackend ${backends[choice]} ${choice}
      new_prompt="${cmd_prompt}-(${backends[choice]}): "
      if [[ $shell_type == "zsh" ]]; then
        if [[ ${rprompt_config} == "true" ]]; then
		export RPROMPT=${profiles[choice]}-${TF_BACKEND}
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
      echo "-> Error:  Number must be one of 0-"$((${#backends[@]} - 1))""
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
