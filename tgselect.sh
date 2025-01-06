#!/usr/bin/env bash

# Haal de lijst op met description en project
task_list=$(toggl ls -f description,project -s $(date +"%F") | sed 's/  \+/;/g;s/ (.*$//g;s/;$//g')

# Zet de lijst om in een multidimensionale array
declare -a tasks
while IFS=';' read -r description project; do
    tasks+=("$description:$project")
done <<< "$task_list"

# Maak een tabel voor gum
table_output=""
for task in "${tasks[@]}"; do
    IFS=':' read -r desc proj <<< "$task"
    table_output+="$desc\t$proj\n"
done

# Gebruik gum om een tabel te tonen en een selectie te maken
selected_task=$(echo -e "$table_output" | gum table -w 200,270 --height 25)

# Splits de geselecteerde taak in description en project
IFS=$'\t' read -r description project <<< "$selected_task"

# Toon het toggl start commando
echo toggl start -o "$project" "$description"

