#!/usr/bin/env bash


# INCLUDE LIB
thisdir="$(dirname "$0")"
source "$thisdir/_bmclib.sh"


checkdeps "jsonify-aws-dotfiles"
checkdeps "jq"
checkdeps "aws"
checkdeps "gum"
checkdeps "aws-mfa"


while [[ $# -gt 0 ]]; do
  case "$1" in
    -l|--list)
      printAWSProfiles
      break
      ;;
    *)
      echo "Unkown option: $1"
      ;;
  esac
done

if [[ -z $1 ]]; then
    loadConfig
    checkOS
    selectAWSProfile
    setDates
    setMFA
    export AWS_PROFILE=${selectedProfileName}
fi

# END OF SCRIPT
