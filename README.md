# OweFlow

**OweFlow** is a backend service written in Go that helps users evenly split shared expenses.
Users log purchases in a shared Google Sheet, and the service calculates who owes whom money on a weekly or monthly basis.

## Features

- Reads a shared Google Sheet for expense data
- Splits purchases evenly between all specified users
- Calculates balances and simplifies who owes whom
- Scheduled to run periodically (weekly or monthly)
- (Eventually) Automatic PayPal payment requests

## How It Works

1. Each user logs expenses in a shared Google Sheet.
2. The backend service reads the sheet on a set schedule.
3. It calculates how much each person owes or is owed.
4. (Eventually) Sends PayPal requests to settle debts automatically.

## Example Spreadsheet Format

Subject to change.

| Date       | Description | Amount | Paid By | Split Between |
| ---------- | ----------- | ------ | ------- | ------------- |
| 2025-06-30 | Groceries   | 80.00  | Alice   | Alice, Bob    |
| 2025-06-30 | Utilities   | 120.00 | Bob     | Alice, Bob    |

## ToDo

- [ ] Create a test Google Sheet
- [ ] Research how to get API access to the Google Sheet (Google Sheets API + service account)
- [ ] Write a simple “Hello World” in Go that reads some dummy data from the test sheet
- [ ] Test the ability to make a new sheet once the current sheet is "processed"
