#!/usr/bin/env bash
# Version: 2024030401 - WvdT
# To-Do:
unset TF_VAR TF_VARS

# Extract the base directory containing '*.tfvars' files
base_directory=$(pwd)
if [[ "$base_directory" == *"stack"* ]]; then
  base_directory=$(dirname "${base_directory%stack*}stack")
fi

# Set TF_ENV variable
# Find all '*.tfvars' files in the base directory
TF_VARS=($(find "${base_directory}" -type f -name "*.tfvars.*" | sort))
TF_VARS_BASE=($(find "${base_directory}" -type f -name "*.tfvars.*" -exec basename {} \; | sort))
TF_VARS_LEN=${#TF_VARS[*]}

function multiple_vars() {

  PS3="Select a var-file number: "
  # echo "-: Unset"
  select item in "${TF_VARS_BASE[@]}"; do
    if [[ -n ${item} ]]; then
      # echo "Using var-file: $item"
      ((REPLY--))
      TF_VAR=${TF_VARS[${REPLY}]} #${item}
      break
    fi
  done
}

if [[ ${TF_VARS_LEN} -eq 1 ]]; then
  echo "just one backend"
  TF_VAR=${TF_VARS}
fi

TF_BACKEND=$(cat .terraform/tfbackend.state 2>/dev/null)

if [[ ! -z ${TF_BACKEND} ]]; then
  TF_ENV=$(basename $(echo $TF_BACKEND) | awk -F '.' '{print $1}' 2>&1)

  for var in "${TF_VARS[@]}"; do
    if [[ $var == *"${TF_ENV}"* ]]; then
      TF_VAR="$var"
    fi
  done
  if [[ -z ${TF_VAR} ]]; then
    multiple_vars
  fi
fi

if [[ ${TF_VARS_LEN} -ge 2 && -z ${TF_VAR} ]]; then
  echo "multiple backends, not matching"
  multiple_vars
fi

echo "Using TF variable-file: ${TF_VAR}"
