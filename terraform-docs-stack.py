#!/usr/bin/env python3
import os
import subprocess

# Configurable variables
FILENAME = "README.md"
STACK_DIRECTORY = "stack"

# Define the tags
BEGIN_TAG = "<!-- BEGIN_TERRAFORM_DOMAINS -->"
END_TAG = "<!-- END_TERRAFORM_DOMAINS -->"

# Find all directories one level deep in the stack directory
map_array = [domain_dir for domain_dir in os.listdir(STACK_DIRECTORY) if os.path.isdir(os.path.join(STACK_DIRECTORY, domain_dir))]

# Array to store results
results = []

# Loop over the directories and execute the command
for domain_dir in map_array:
    # Add a header for the stack directory
    header = f"## Domain: {domain_dir}\n"
    results.append(header)

    # Execute the command in the directory
    result = subprocess.run(["terraform-docs", "markdown", os.path.join(STACK_DIRECTORY, domain_dir)], capture_output=True, text=True)
    results.append(result.stdout)

# Check if README.md exists, if not, create it
if not os.path.exists(FILENAME):
    with open(FILENAME, "w") as f:
        f.write("")

# Read the file
with open(FILENAME, "r") as f:
    content = f.read()

# Find the position of the tags
begin_position = content.find(BEGIN_TAG)
end_position = content.find(END_TAG)

# If the tags are not found, append them along with the results array to the end of the file
if begin_position == -1 or end_position == -1:
    print("One or both tags not found.")

    # Append the tags and results array to the end of the file
    content += f"\n{BEGIN_TAG}\n"
    content += "\n".join(results)
    content += f"\n{END_TAG}\n"
    print("Tags added to the document.")
else:
    # Determine the lines to insert the results
    insert_start = begin_position + len(BEGIN_TAG) + 1
    insert_end = end_position - 1

    # Remove the current content between the tags
    content = content[:insert_start] + content[insert_end:]

    # Insert the results between the tags
    for result in results:
        content = content[:insert_start] + result + content[insert_start:]

    print("Results inserted between the tags.")

# Write the updated content back to the file
with open(FILENAME, "w") as f:
    f.write(content)
