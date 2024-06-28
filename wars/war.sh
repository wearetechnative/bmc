#!/usr/bin/env bash

# Check if jq is installed
if ! command -v jq &> /dev/null
then
    echo "program 'jq' not found. Please install it before executing this script."
    exit
fi

# Default JSON file
config_file="config.json"

# Process the optional -i parameter
while getopts "i:" opt; do
    case $opt in
        i)
            config_file=$OPTARG
            ;;
        *)
            echo "Usage: $0 [-i JSON-file]"
            exit 1
            ;;
    esac
done

# Read the JSON variables
customerName=$(jq -r '.customerName' "$config_file")
role_name=$(jq -r '.role_name' "$config_file")
aws_regions=$(jq -r '.awsRegions[]' "$config_file")
environments=$(jq -c '.environments[]' "$config_file")

# Create the customerName directory
mkdir -p "$customerName"
timestamp=$(date +"%Y%m%d-%H%M%S")
errorLogFile=${customerName}-errors-${timestamp}.log
exec 2> "$errorLogFile"

# Loop through each environment and execute the commands
for environment in $environments; do
    awsAccountNumber=$(echo "$environment" | jq -r '.awsAccountNumber')
    environmentName=$(echo "$environment" | jq -r '.environmentName')

    # Set AWS_PROFILE
    export AWS_PROFILE="war-${customerName}-${environmentName}"

    # Check if AWS_PROFILE exists in $HOME/.aws/config
    if ! grep -q "\[profile ${AWS_PROFILE}\]" "$HOME/.aws/config"; then
        echo "!! Profile ${AWS_PROFILE} not found at ${timestamp}" |tee -a ${errorLogFile}
        continue
    fi

    # Loop through each AWS region
    for region in $aws_regions; do
        # Create the environment + region directory
        environment_dir="$customerName/$awsAccountNumber-$environmentName/$region"
        mkdir -p "$environment_dir"

        # Change to the environment + region directory
        cd "$environment_dir" || exit

        echo "-- RUNNING PROWLER for: ${AWS_PROFILE} - ${region}"
        prowler aws -f ${region}  --compliance aws_well_architected_framework_security_pillar_aws --output-directory ./security_pillar_output/
        prowler aws -f ${region} --compliance aws_well_architected_framework_reliability_pillar_aws --output-directory ./reliability_pillar_output/
        prowler aws -f ${region} --compliance aws_foundational_security_best_practices_aws --output-directory ./best_practices_output/

        # Change back to the main directory
        cd - || exit
    done
done

