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

# Build the table with instance details and scheduler status
header=$(echo "$instances_json" | jq -r '["InstanceId", "Name", "State", "SchedulerStatus", "SchedulerValue"] | @csv')
instances=$(echo "$instances_json" | jq -r '.[] |
  ((.Tags[]? | select(.Key=="InstanceScheduler" or .Key=="InstanceScheduler_DISABLED") | .Key) // "none") as $tagKey |
  ((.Tags[]? | select(.Key=="InstanceScheduler" or .Key=="InstanceScheduler_DISABLED") | .Value) // "N/A") as $tagValue |
  (if $tagKey == "InstanceScheduler" then "enabled" elif $tagKey == "InstanceScheduler_DISABLED" then "disabled" else "none" end) as $status |
  [
    .InstanceId,
    ((.Tags[]? | select(.Key=="Name") | .Value) // "N/A"),
    .State.Name,
    $status,
    $tagValue
  ] | @csv')

formatted_instances=$(echo -e "$header\n$instances")

# Let user select an instance
selected_line=$(echo "$formatted_instances" | gum table -w 20,35,12,16,20 --height 20)

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
    tagKey: ((.Tags[]? | select(.Key=="InstanceScheduler" or .Key=="InstanceScheduler_DISABLED") | .Key) // ""),
    tagValue: ((.Tags[]? | select(.Key=="InstanceScheduler" or .Key=="InstanceScheduler_DISABLED") | .Value) // ""),
    availabilityZone: .Placement.AvailabilityZone
  }')

current_tag_key=$(echo "$instance_data" | jq -r '.tagKey // empty')
current_tag_value=$(echo "$instance_data" | jq -r '.tagValue // empty')
availability_zone=$(echo "$instance_data" | jq -r '.availabilityZone')

# Extract region from availability zone (remove last character, e.g., eu-central-1a -> eu-central-1)
region="${availability_zone%?}"

# Handle instances without scheduler tags
if [ -z "$current_tag_key" ] || [ "$current_tag_key" = "null" ]; then
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

# Determine new tag key
if [ "$current_tag_key" = "InstanceScheduler" ]; then
  new_tag_key="InstanceScheduler_DISABLED"
  action_desc="disable scheduler"
  new_status="disabled"
else
  new_tag_key="InstanceScheduler"
  action_desc="enable scheduler"
  new_status="enabled"
fi

# Show current status and ask for confirmation
echo ""
echo "Instance: ${instance_id}"
echo "Current tag: ${current_tag_key}"
echo "Current value: ${current_tag_value}"
echo ""
if ! gum confirm "Do you want to ${action_desc}?"; then
  echo "-- Cancelled. No changes made."
  exit 0
fi

# Perform the tag toggle
echo ""
echo "-- $(echo ${action_desc:0:1} | tr '[a-z]' '[A-Z]')${action_desc:1} for instance ${instance_id}..."

# Delete old tag
aws ec2 delete-tags --resources "$instance_id" --tags "Key=${current_tag_key}" >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "!! Error: Failed to remove old tag '${current_tag_key}'. Check error above."
  exit 1
fi

# Create new tag with same value
aws ec2 create-tags --resources "$instance_id" --tags "Key=${new_tag_key},Value=${current_tag_value}" >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "!! Error: Failed to create new tag '${new_tag_key}'. Check error above."
  exit 1
fi

echo "-- Successfully toggled scheduler tag on instance ${instance_id}"
echo "   From: ${current_tag_key} (value: ${current_tag_value})"
echo "   To:   ${new_tag_key} (value: ${current_tag_value})"
