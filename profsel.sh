#!/usr/bin/env bash

# INCLUDE LIB
thisdir="$(dirname "$0")"
source "$thisdir/_bmclib.sh"

##### PLACE YOUR COMMANDS BELOW #####

make_command "cli" "Open account in AWS console."
cli(){

  checkdeps "jsonify-aws-dotfiles"
  checkdeps "jq"
  checkdeps "aws"
  checkdeps "gum"
  checkdeps "aws-mfa"

  checkOS
  selectAWSProfile
  setDates
  setMFA
  export AWS_PROFILE=${selectedProfileName}
}

make_command "console" "Open account in AWS console."
console(){

  # CHECK deps
  loadConfig
  checkdeps "jsonify-aws-dotfiles"
  checkdeps "jq"
  checkdeps "awk"
  checkdeps "assumego"
  checkdeps "gum"

  checkOS
  selectAWSProfile
  setDates
  setMFA

  GRANTED_ALIAS_CONFIGURED="true" assumego -c $selectedProfileName
}

##### PLACE YOUR COMMANDS ABOVE #####

runme
