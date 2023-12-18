#!/bin/zsh
awsconfigfile=${HOME}/.aws/config
awsOnlineID=$(aws sts get-caller-identity 2>/dev/null |grep -i account |awk  '{print $2}' |tr -d "\",")
awsprofiles=($(grep -i sso_session ${awsconfigfile} |awk '{print $3}'))
select item in "${awsprofiles[@]}" Quit
do
    case $REPLY in
        1|2)
            awsprofileid=$(grep -i -B1 ${item}  ~/.aws/config |grep profile|awk -F "-" '{print $2}'|tr -d "]")
            export AWS_PROFILE=$(grep -i -B1 ${item}  ~/.aws/config |grep profile |awk '{print $2}' |tr -d "]")
            export RPROMPT=${item}
            if [[ ${awsOnlineID} != ${awsprofileid} ]];
            then
                aws sso login --sso-session ${item}
            fi
            break
    ;;
        3)
	unset RPROMPT
	unset AWS_PROFILE
	aws sso logout
        break
        ;;
        $((${#items[@]}+1))) echo "We're done!"; break;;
        *) echo "Ooops - unknown choice $REPLY";;
    esac
break
done


awsID=$(echo ${AWS_PROFILE})
