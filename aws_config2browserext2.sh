#!/bin/bash

find_blocks_with_role_arn_and_region() {
    local config_file="$1"
    local profile_header=""
    local role_arn=""
    local region=""
    local temp_file="$(mktemp)"

    # Tijdelijk bestand aanmaken om gegevens per profile_header op te slaan
    touch "$temp_file"

    while IFS= read -r line; do
        if [[ $line == \[profile\ * ]]; then
            if [ -n "$profile_header" ] && [ -n "$role_arn" ]; then
                printf '%s\n' "$profile_header" >> "$temp_file"
                printf '%s\n' "role_arn = $role_arn" >> "$temp_file"
                printf '%s\n' "region = $region" >> "$temp_file"
                printf '%s\n' "" >> "$temp_file"
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
        printf '%s\n' "$profile_header" >> "$temp_file"
        printf '%s\n' "role_arn = $role_arn" >> "$temp_file"
        printf '%s\n' "region = $region" >> "$temp_file"
    fi

    # Sorteer de georganiseerde gegevens op profile_header
    sort -o "$temp_file" "$temp_file"

    # Print de georganiseerde gegevens gesorteerd op profile_header
    cat "$temp_file"

    # Verwijder het tijdelijke bestand
    rm "$temp_file"
}

config_file="${HOME}/.aws/config"

find_blocks_with_role_arn_and_region "$config_file"

