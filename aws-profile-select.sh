#!/usr/bin/env bash


# INCLUDE LIB
thisdir="$(dirname "$0")"
source "$thisdir/_bmclib.sh"


checkdeps "jsonify-aws-dotfiles"
checkdeps "jq"
checkdeps "aws"
checkdeps "gum"
checkdeps "aws-mfa"
deps_missing

loadConfig
checkOS
selectAWSProfile "$@"
setDates
setMFA
export AWS_PROFILE=${selectedProfileName}

# END OF SCRIPT
