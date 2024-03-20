#!/usr/bin/env bash

# INCLUDE LIB
thisdir="$(dirname "$0")"
source "$thisdir/_get_var_file.sh"

if [[ ! -z ${TF_VAR} ]]; then
  terraform apply -var-file=${TF_VAR} $@
else
  terraform apply $@
fi
