#!/usr/bin/env python3 
from datetime import datetime
import pytz

# Ask the user to input a date and time
mydate = input("Enter a date and time (e.g., September 12, 2025, 01:59:59 (UTC+02:00)): ")

# Step 1: Remove the timezone information
date_without_timezone = mydate.split(" (")[0]

# Step 2: Parse the date without the timezone
parsed_date = datetime.strptime(date_without_timezone, "%B %d, %Y, %H:%M:%S")

# Step 3: Add the timezone information back
# UTC+02:00 is in the form of 'Etc/GMT-2' because positive offsets are negative in 'Etc/GMT'
timezone_str = mydate.split(" (")[1][:-1]  # Extract the timezone
utc_offset_hours = int(timezone_str[3:6]) * -1  # Convert to Etc/GMT format
timezone = pytz.timezone(f"Etc/GMT{utc_offset_hours}")
localized_date = timezone.localize(parsed_date)

# Step 4: Convert to the desired format
formatted_date = localized_date.strftime("%Y-%m-%d %H:%M")

print("Formatted date and time:", formatted_date)

