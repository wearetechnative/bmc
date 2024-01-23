#!/bin/bash

TF_ENV=$(echo $TF_BACKEND |awk -F '.' '{print $1}')

if [[ -f ${TF_ENV}.tfvars ]]; then
    terraform plan -var-file=${TF_ENV}.tfvars $@
else
    terraform plan $@
fi

