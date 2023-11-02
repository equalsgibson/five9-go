#!/bin/bash
## Following domains are allowed to be included in files.
WHITELISTEDDOMAINS=("example.com" "api.prod.us.five9.net" "app.five9.com" "cobertura.sourceforge.net")

## Remove the trailing \| from the joined string
printf -v joinedWhitelist '%s\|' "${WHITELISTEDDOMAINS[@]}"
DOMAINSTRING=$(echo "${joinedWhitelist%\\|}")

if [ "$(grep -rI -Eo '(www\.|:\/\/|@)([a-zA-Z][a-zA-Z0-9\-_]{1,61}).?([a-zA-Z][a-zA-Z0-9\-_]{1,61}\.?)+\/?' five9 | grep -vc ${DOMAINSTRING} )" -gt 0 ]; 
then 
    echo "Found a non whitelisted domain:";
    grep -rI -Eo '(www\.|:\/\/|@)([a-zA-Z][a-zA-Z0-9\-_]{1,61}).?([a-zA-Z][a-zA-Z0-9\-_]{1,61}\.?)+\/?' five9 | grep -v ${DOMAINSTRING}
    exit 1; 
fi
exit 0;