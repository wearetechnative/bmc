#!/usr/bin/env bash

# INCLUDE LIB
thisdir="$(dirname "$0")"
source "$thisdir/_bmclib.sh"

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

