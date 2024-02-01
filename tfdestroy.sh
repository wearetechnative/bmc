#!/bin/bash

TF_ENV=$(echo $TF_BACKEND | awk -F '.' '{print $1}' 2>&1)
TF_VARS=$(find . -type f -name "*.tfvars")


if [[ -f ${TF_ENV}.tfvars ]]; then
    terraform destroy -var-file=${TF_ENV}.tfvars $@
elif [[ ! -z ${TF_VARS} ]]; then

    TF_VARS=(${TF_VARS})
    TF_VARS_len=${#TF_VARS[*]}


    echo "--- Choose vars-file"
    echo "-: Quit"
    for ((i = 0; i < $TF_VARS_len; i++)); do
        echo "$i: ${TF_VARS[$i]}"
    done
    echo
    printf 'Selection: '
    read choice
    COLUMN=6
    case $choice in
    '' | *[!0-9\-]*)
        clear
        echo Invalid selection. Make a valid selection from the list above or press ctrl+c to exit
        echo '-> Error: Not a number, and not "-"'
        echo
        break
        ;;
    esac
    in_range=false
    while [ $in_range != true ]; do
        if [[ $choice == '-' ]]; then
            break
        elif (($choice >= 0)) && (($choice <= (${TF_VARS_len}-1))); then

            echo ${TF_VARS[choice]}
            in_range=true
            terraform destroy -var-file=${TF_VARS[choice]} $@
            break
        else
            echo "!! Not a valid option"
            break
        fi
    done

fi
