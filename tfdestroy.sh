#!/bin/bash

TF_ENV=$(echo $TF_BACKEND | awk -F '.' '{print $1}')

if [[ $1 == "" ]]; then
    objects="$(for object in $(terraform state list | grep -vE "^data\." | grep -vE "backend|dynamodb|kms"); do echo -target="${object}"\ ; done)"
    IFSBAK=$IFS
    IFS=$'\n'
    objects=($objects)
    IFS=$IFSBAK
    objects_len=${#objects[*]}
else
    unset objects
fi

if [[ -f ${TF_ENV}.tfvars ]]; then
    # terraform destroy -var-file=${TF_ENV}.tfvars ${objects[*]} $#
    terraform destroy -var-file=${TF_ENV}.tfvars ${objects[*]} $@
else
    terraform destroy ${objects[*]} $@
fi
