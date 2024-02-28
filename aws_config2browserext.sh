#!/bin/bash

find_blocks_with_role_arn_and_region() {
    local config_file="$1"
    local profile_header=""
    local role_arn=""
    local region=""

    while IFS= read -r line; do
        if [[ $line == \[profile\ * ]]; then
            if [ -n "$profile_header" ] && [ -n "$role_arn" ]; then
                echo "$profile_header"
                echo "role_arn = $role_arn"
                echo "region = $region"
                echo
            fi
            profile_header="$line"
            unset role_arn region
        fi

        if [[ $line == role_arn* ]]; then
            role_arn_value="${line#role_arn = }"
            if [ -n "$role_arn_value" ]; then
                role_arn="$role_arn_value"
            fi
        elif [[ $line == region* ]]; then
            region_value="${line#region = }"
            if [ -n "$region_value" ]; then
                region="$region_value"
            fi
        fi
    done < "$config_file"

    if [ -n "$profile_header" ] && [ -n "$role_arn" ]; then
        echo "$profile_header"
        echo "role_arn = $role_arn"
        echo "region = $region"
    fi
}

config_file="${HOME}/.aws/config"

find_blocks_with_role_arn_and_region "$config_file"

