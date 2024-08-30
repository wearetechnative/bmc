checkdeps(){
  if ! command -v $1 &> /dev/null
  then
    echo "<$1> could not be found"
    echo "  install this program first"
    exit 1
  fi
}

##### PLACE YOUR COMMANDS BELOW #####

cli(){

  # CHECK deps
  checkdeps "dasel"
  checkdeps "jsonify-aws-dotfiles"
  checkdeps "jq"
  checkdeps "gum"

  GROUP=`jsonify-aws-dotfiles | jq -r '[.config[].group]| unique| sort| .[]' | grep -v null | gum choose --height 25`

  if [ -z "${GROUP}" ] || [ "${GROUP}" = "user aborted" ]; then
    exit 1
  fi

  echo $GROUP

  ACCOUNTTMP=`jsonify-aws-dotfiles | jq "[.config | with_entries(select(.value.group == \"$GROUP\")) | to_entries[] | {\"AWS Account\": .key, \"Role\": .value.role_arn}]" | dasel -r json -w csv | gum table -w 30,70 --height 25`

  if [ -z "${ACCOUNTTMP}" ]; then
    exit 1
  fi

  ACCOUNT=`echo $ACCOUNTTMP | cut -d ',' -f 1`

  if [ -z "${ACCOUNT}" ]; then
    exit 1
  fi

  echo "AWS_PROFILE set to ${ACCOUNT}"
  export AWS_PROFILE="${ACCOUNT}"
}

##### PLACE YOUR COMMANDS ABOVE #####

cli

