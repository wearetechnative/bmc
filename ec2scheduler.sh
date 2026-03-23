#!/usr/bin/env bash

# Query AWS for all EC2 instances
aws_output=$(aws ec2 describe-instances)

if [ $? -ne 0 ]; then
  echo "!! Error: Can't query instances. Check error above."
  exit 1
fi

# Flatten the Reservations structure and filter out terminated instances
instances_json=$(echo "$aws_output" | jq '[.Reservations[].Instances[] | select(.State.Name != "terminated")]')

# Check if any instances were found
instance_count=$(echo "$instances_json" | jq 'length')
if [ "$instance_count" -eq 0 ]; then
  echo "-- No instances found in current AWS profile."
  exit 0
fi

# Build the table with instance details, Scheduler status, and Ignore_scheduler status
header=$(echo "$instances_json" | jq -r '["InstanceId", "Name", "State", "Scheduler", "IgnoreUntil"] | @csv')
instances=$(echo "$instances_json" | jq -r '.[] |
  ((.Tags[]? | select(.Key=="InstanceScheduler") | .Value) // "") as $instanceScheduler |
  ((.Tags[]? | select(.Key=="Ignore_scheduler") | .Value) // "N/A") as $ignoreUntil |
  (if $instanceScheduler != "" then "yes" else "no" end) as $schedulerStatus |
  [
    .InstanceId,
    ((.Tags[]? | select(.Key=="Name") | .Value) // "N/A"),
    .State.Name,
    $schedulerStatus,
    $ignoreUntil
  ] | @csv')

formatted_instances=$(echo -e "$header\n$instances")

# Let user select an instance
selected_line=$(echo "$formatted_instances" | gum table -w 20,30,12,10,30 --height 20)

if [ -z "$selected_line" ]; then
  echo "-- No instance selected. Exiting."
  exit 0
fi

# Extract instance ID from the selected line (first field in CSV)
instance_id=$(echo "$selected_line" | awk -F, '{print $1}' | tr -d '"')

# Get current tag details and region for this instance
instance_data=$(echo "$instances_json" | jq -r --arg INSTANCE_ID "$instance_id" '
  .[] | select(.InstanceId == $INSTANCE_ID) |
  {
    instanceScheduler: ((.Tags[]? | select(.Key=="InstanceScheduler") | .Value) // ""),
    ignoreScheduler: ((.Tags[]? | select(.Key=="Ignore_scheduler") | .Value) // ""),
    availabilityZone: .Placement.AvailabilityZone
  }')

instance_scheduler_value=$(echo "$instance_data" | jq -r '.instanceScheduler // empty')
ignore_scheduler_value=$(echo "$instance_data" | jq -r '.ignoreScheduler // empty')
availability_zone=$(echo "$instance_data" | jq -r '.availabilityZone')

# Extract region from availability zone (remove last character, e.g., eu-central-1a -> eu-central-1)
region="${availability_zone%?}"

# Handle instances without scheduler tags
if [ -z "$instance_scheduler_value" ] || [ "$instance_scheduler_value" = "null" ]; then
  echo ""
  echo "Instance: ${instance_id}"
  echo "Status: No scheduler tag configured"
  echo ""

  if ! gum confirm "This instance does not have a scheduler tag. Would you like to add one?"; then
    echo "-- Cancelled. No changes made."
    exit 0
  fi

  echo ""
  echo "=== Instructions to Add InstanceScheduler Tag ==="
  echo ""
  echo "To add the scheduler tag, you need to:"
  echo "1. Go to the AWS EC2 Console"
  echo "2. Find instance: ${instance_id}"
  echo "3. Add tag with:"
  echo "   - Key: InstanceScheduler"
  echo "   - Value: <your schedule definition>"
  echo ""
  echo "Example schedule values:"
  echo "  - office-hours"
  echo "  - weekdays-9to5"
  echo "  - 24x7"
  echo ""

  if gum confirm "Would you like to open the AWS Console now?"; then
    echo ""
    echo "-- Opening AWS Console to instance details page..."
    echo "   Profile: ${AWS_PROFILE:-default}"
    echo "   Instance: ${instance_id}"
    echo "   Region: ${region}"

    # Construct the console URL directly to the instance details page
    console_url="https://${region}.console.aws.amazon.com/ec2/home?region=${region}#InstanceDetails:instanceId=${instance_id}"

    # Use assumego to open the console with the proper profile and destination
    GRANTED_ALIAS_CONFIGURED="true" GRANTED_ENABLE_AUTO_REASSUME=true \
      assumego --duration 3600s -c "${AWS_PROFILE}" --console-destination "${console_url}"

    exit 0
  else
    echo ""
    echo "-- No problem. You can add the tag later through the AWS Console."
    echo "-- Run 'bmc ec2scheduler' again after adding the tag."
    exit 0
  fi
fi

# Show current Ignore_scheduler status
echo ""
echo "Instance: ${instance_id}"
if [ -n "$ignore_scheduler_value" ] && [ "$ignore_scheduler_value" != "null" ]; then
  echo "Current ignore override: ${ignore_scheduler_value}"
else
  echo "Current ignore override: Not set"
fi
echo ""

# Present action menu
action=$(gum choose "Set ignore until time" "Remove ignore override" "Cancel")

if [ -z "$action" ] || [ "$action" = "Cancel" ]; then
  echo "-- Cancelled. No changes made."
  exit 0
fi

# Handle menu selection
case "$action" in
  "Set ignore until time")
    # Prompt for time in HH:MM format
    while true; do
      time_input=$(gum input --placeholder "Example: 22:00" --prompt "Enter time (HH:MM): ")

      if [ -z "$time_input" ]; then
        echo "-- Cancelled. No changes made."
        exit 0
      fi

      # Validate time format (HH:MM with 24-hour format)
      if echo "$time_input" | grep -qE '^([01][0-9]|2[0-3]):[0-5][0-9]$'; then
        break
      else
        echo "!! Invalid time format. Please use HH:MM (24-hour format, e.g., 22:00, 08:30)"
        echo ""
      fi
    done

    # Prompt for timezone
    timezone_input=$(gum input --placeholder "Examples: Europe/Amsterdam, America/New_York, UTC" --prompt "Enter timezone: ")

    if [ -z "$timezone_input" ]; then
      echo "-- Cancelled. No changes made."
      exit 0
    fi

    # Combine time and timezone
    ignore_value="${time_input} ${timezone_input}"

    echo ""
    echo "-- Setting ignore override for instance ${instance_id}..."
    echo "   Ignore until: ${ignore_value}"

    # Create or update Ignore_scheduler tag
    aws ec2 create-tags --resources "$instance_id" --tags "Key=Ignore_scheduler,Value=${ignore_value}" >/dev/null 2>&1
    if [ $? -ne 0 ]; then
      echo "!! Error: Failed to set Ignore_scheduler tag. Check error above."
      exit 1
    fi

    echo "-- Successfully set ignore override"
    echo "   Instance ${instance_id} will ignore scheduled stops until ${ignore_value}"
    echo "   The tag will be automatically removed after this time."
    ;;

  "Remove ignore override")
    # Check if tag exists
    if [ -z "$ignore_scheduler_value" ] || [ "$ignore_scheduler_value" = "null" ]; then
      echo "-- No ignore override is currently set."
      exit 0
    fi

    echo ""
    echo "-- Removing ignore override for instance ${instance_id}..."

    # Delete Ignore_scheduler tag
    aws ec2 delete-tags --resources "$instance_id" --tags "Key=Ignore_scheduler" >/dev/null 2>&1
    if [ $? -ne 0 ]; then
      echo "!! Error: Failed to remove Ignore_scheduler tag. Check error above."
      exit 1
    fi

    echo "-- Successfully removed ignore override"
    echo "   Instance ${instance_id} will resume normal schedule."
    ;;
esac
