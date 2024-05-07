#!/usr/bin/env bash
# W. van der Toorren - TechNative B>V
# version: 2024043001
# 
 
#set -x

# Init variables
breadcrumbs=""

print_breadcrumbs() {
  ## Function to print breadcrumbs, to give some information about the choices
  printf "\r%s\n" "$1"
}


# Function to update breadcrumbs
update_breadcrumbs() {
  if [[ -z "$breadcrumbs" ]]; then
    breadcrumbs="$1"
  else
    breadcrumbs="$breadcrumbs > $1"
  fi
  print_breadcrumbs "$breadcrumbs"
}

# #Menu to pass ssh-k
# while getopts 'i:' opt; do
#   case "$opt" in
#   i)
#     sshKey="$OPTARG"
#     if [[ -f ${sshKey} ]]; then
#       echo "Using ssh-key: '${OPTARG}' "
#       sshKey="-i ${OPTARG}"
#     else
#       echo "ssh-key not found: '${OPTARG}'. Not using it"
#       unset sshKey
#     fi
#     ;;

#   ?)
#     echo -e "Invalid command option.\nUsage: $(basename $0) [-i path_to/ssh_key]"
#     exit 1
#     ;;
#   esac
# done
# shift "$(($OPTIND - 1))"

# Check AWS credentials via AWS_PROFILE
if [[ -z ${AWS_PROFILE} ]]; then
  echo "!! AWS_PROFILE  not set"
  exit 1
fi

# if [[ -z ${AWS_DEFAULT_REGION} ]]; then
#   echo "!! AWS_DEFAULT_REGION not set"
#   exit 1
# fi

# execute AWS CLI-command to retrieve info about clusters. Exit when command fails
aws_output=$(aws ecs list-clusters 2>/dev/null)
if [[ $? -ne 0 ]]; then
  echo "!! Error listing clusters. Check AWS credentials"
  exit 1
fi

# Extract fields needed from JSON_output with the use of jq-command
echo "-- Select cluster"
clusters=$(echo "$aws_output" | jq -r '.clusterArns[] | split("/")[-1]')
cluster=$(gum choose $clusters)
update_breadcrumbs "$cluster" # update breadcrumbs with clustername
# clustertype=$(aws ecs describe-clusters --cluster ${cluster} | jq -r '.clusters[0].capacityProviders[]')

# execute AWS CLI-command to retrieve info about services in the cluster. Exit when command fails
aws_output=$(aws ecs list-services --cluster ${cluster} 2>/dev/null)
if [[ $? -ne 0 ]]; then
  echo "!! Error listing services"
  exit 1
fi

echo "-- Select service"
services=$(echo "$aws_output" | jq -r '.serviceArns[] | split("/")[-1]')
service=$(gum choose $services)
update_breadcrumbs "$service" # update breadcrumbs with clustername

# execute AWS CLI-command to retrieve info about tasks in the service. Exit when command fails
aws_output=$(aws ecs list-tasks --cluster $cluster --service $service --desired-status RUNNING --query 'taskArns[]')
if [[ $? -ne 0 ]]; then
  echo "!! Error listing tasks."
  exit 1
fi

echo "-- Select task_id"
# Filter taskArns based on lastStatus
taskids=$(echo "$aws_output" |jq -r '.[] | split("/") | last')
taskid=$(gum choose $taskids)
update_breadcrumbs "$taskid"# update breadcrumbs with clustername

aws_output=$(aws ecs describe-tasks --cluster ${cluster} --tasks ${taskids})
if [[ $? -ne 0 ]]; then
  echo "!! Error listing containers. Check AWS credentials"
  exit 1
fi

containers=$(echo "$aws_output" | jq -r '.tasks[0].containers[].name')
echo "-- Select container"
container=$(gum choose $containers)
taskarn=$(echo "$aws_output" | jq -r '.tasks[0].taskArn')
update_breadcrumbs "$container"

aws ecs execute-command --cluster ${cluster} --interactive --container ${container} --command /bin/bash --task ${taskarn}
