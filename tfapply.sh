#!/bin/bash



TF_ENV=$(echo $TF_BACKEND |awk -F '.' '{print $1}')

if [[ -f ${TF_ENV}.tfvars ]]; then
    terraform apply -var-file=${TF_ENV}.tfvars $@
else
    terraform apply $@
fi
