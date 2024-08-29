#!/usr/bin/env bash
#(C)2019-2022 Pim Snel - https://github.com/mipmip/RUNME.sh
CMDS=();DESC=();NARGS=$#;ARG1=$1;make_command(){ CMDS+=($1);DESC+=("$2");};usage(){ printf "\nUsage: %s [command]\n\nCommands:\n" $0;line="              ";for((i=0;i<=$(( ${#CMDS[*]} -1));i++));do printf "  %s %s ${DESC[$i]}\n" ${CMDS[$i]} "${line:${#CMDS[$i]}}";done;echo;};runme(){ if test $NARGS -eq 1;then eval "$ARG1"||usage;else usage;fi;}

checkdeps(){
  if ! command -v $1 &> /dev/null
  then
    echo "<$1> could not be found"
    echo "  install this program first"
    exit 1
  fi
}

##### PLACE YOUR COMMANDS BELOW #####

make_command "console" "Open account in AWS console."
console(){

  # CHECK deps
  checkdeps "dasel"
  checkdeps "jsonify-aws-dotfiles"
  checkdeps "jq"
  checkdeps "assumego"
  checkdeps "gum"

  GROUP=`jsonify-aws-dotfiles | jq -r '[.config[].group]| unique| sort| .[]' | grep -v null | gum choose --height 25`

  if [ -z "${GROUP}" ] || [ "${GROUP}" = "user aborted" ]; then
    exit 1
  fi

  #echo $GROUP

  ACCOUNTTMP=`jsonify-aws-dotfiles | jq "[.config | with_entries(select(.value.group == \"$GROUP\")) | to_entries[] | {\"AWS Account\": .key, \"Role\": .value.role_arn}]" | dasel -r json -w csv | gum table -w 30,70 --height 25`

  if [ -z "${ACCOUNTTMP}" ]; then
    exit 1
  fi

  ACCOUNT=`echo $ACCOUNTTMP | cut -d ',' -f 1`

  if [ -z "${ACCOUNT}" ]; then
    exit 1
  fi

  GRANTED_ALIAS_CONFIGURED="true" assumego -c $ACCOUNT
}

##### PLACE YOUR COMMANDS ABOVE #####

runme
