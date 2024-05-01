#!/usr/bin/env python3
import datetime

def current_week_number():
    today = datetime.date.today()
    return today.isocalendar()[1]

week_number = current_week_number()
print("Current weeknummer: ", week_number)

