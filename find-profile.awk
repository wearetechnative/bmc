#!/usr/bin/awk -f

# Gebruik: awk -v id=599192675977 -f find-profile.awk ~/.aws/config

BEGIN {
  found = 0
}

/^\[profile / {
  profile_line = $0
  match(profile_line, /^\[profile (.+)\]$/, m)
  if (m[1] != "") {
    current_profile = m[1]
  }
}

/^\s*aws_account_id\s*=\s*[0-9]+/ {
  match($0, /[0-9]+/, idmatch)
  if (idmatch[0] == id) {
    print current_profile
    found = 1
    exit
  }
}

END {
  if (!found) {
    # exit code 1 zodat je in shell fout kunt afhandelen
    exit 1
  }
}

