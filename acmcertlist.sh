#!/usr/bin/env bash 


# Standaardwaarde voor dagen tot verval
days_until_expiry=0
show_only_expiring=false

# Opties verwerken
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --ed) days_until_expiry="$2"; show_only_expiring=true; shift ;;
        *) echo "Onbekende parameter doorgegeven: $1"; exit 1 ;;
    esac
    shift
done

# Titel van de tabel
printf "%-50s %-50s %-30s\n" "Domainname" "Expiration Date" "ARN"
printf "%s\n" "----------------------------------------------------------------------------------------------------"

# Alle certificaat-ARN's ophalen
cert_arns=$(aws acm list-certificates --query "CertificateSummaryList[*].[CertificateArn]" --output text)

# Loop door de lijst van ARN's en controleer vervaldatums
for cert_arn in $cert_arns; do
    # Haal de vervaldatum en de domeinnaam op
    cert_details=$(aws acm describe-certificate --certificate-arn $cert_arn --query "Certificate.{NotAfter:NotAfter,DomainName:DomainName}" --output json)
    expiry_date=$(echo $cert_details | jq -r '.NotAfter')
    domain_name=$(echo $cert_details | jq -r '.DomainName')
    
    expiry_timestamp=$(date -d $expiry_date +%s)
    current_timestamp=$(date +%s)
    difference=$(( (expiry_timestamp - current_timestamp) / 86400 ))
    
    if [ "$show_only_expiring" = true ] && [ $difference -le $days_until_expiry ]; then
        printf "%-50s %-50s %-30s\n" "$domain_name" "$expiry_date" "$cert_arn" 
    elif [ "$show_only_expiring" = false ]; then
        printf "%-50s %-50s %-30s\n" "$domain_name" "$expiry_date" "$cert_arn"
    fi
done

