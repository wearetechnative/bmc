#!/usr/bin/env python3
import argparse
import datetime

def current_week_number(date):
    if date:
        given_date = datetime.datetime.strptime(date, "%Y-%m-%d").date()
        week_number = given_date.isocalendar()[1]
        return week_number, given_date.strftime("%Y-%m-%d")
    else:
        today = datetime.date.today()
        return today.isocalendar()[1], today.strftime("%Y-%m-%d")

def main():
    parser = argparse.ArgumentParser(description="Script to determine the week number.")
    parser.add_argument("-d", "--date", metavar="YYYY-MM-DD", help="Specify a date to determine the week number for that date.")
    args = parser.parse_args()

    week_number, given_date = current_week_number(args.date)
    if args.date:
        print(f"The week number for the given date {given_date} is: {week_number}")
    else:
        print("The current week number is:", week_number)

if __name__ == "__main__":
    main()

