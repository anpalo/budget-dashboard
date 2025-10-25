import csv, sqlite3, re
from datetime import datetime

DB = "budget.db"
CSV = "CurrentBudget.csv"
MONTH_MAP = {
    "Jan": "01",
    "Feb": "02",
    "Mar": "03",
    "Apr": "04",
    "May": "05",
    "Jun": "06",
    "Jul": "07",
    "Aug": "08",
    "Sep": "09",
    "Oct": "10",
    "Nov": "11",
    "Dec": "12",
}

# https://www.shortcutfoo.com/app/dojos/regex/cheatsheet
def to_float(s):
    if s is None:
        return None
    s = s.strip()
    if s == "" or s.isspace():
        return None
    s = re.sub(r'[^\d\.-]', '', s)

    try:
        return float(s)
    except:
        return None

# https://docs.python.org/3/library/datetime.html

def is_date(s):
    try:
        default_year = "2025"
        s = datetime.strptime(f"{s}-{default_year}", "%d-%b-%Y" )
        return True
    except:
        return False

def parse_date(s):
    try:
        default_year = "2025"
        dt = datetime.strptime(f"{s}-{default_year}", "%d-%b-%Y")
        return dt
    except:
        return None


connect = sqlite3.connect(DB)
cur = connect.cursor() # to use SQL commands eg. SELECT INSERT UPDATE DELETE


SKIP_CATEGORIES = {"Daily Sum","Monthly Sum","Adjusted Monthly Sum","Income","Net Month Income","Total Savings"}

# https://docs.python.org/3/library/csv.html
with open(CSV, newline = '', encoding='utf-8') as f:
    reader = csv.DictReader(f)
    to_insert = []
    to_insert_summary = []
    monthly_data = {}
    total_savings = 0.0
    for r in reader:
        
        date_str = r.get("Date", "").strip()
        if not is_date(date_str):
            continue
        dt = parse_date(date_str)
        month_key = dt.strftime("%Y-%m") 
  
        savings = r.get("Total Savings")
        if savings:
            total_savings = to_float(savings)


        if month_key not in monthly_data:
            monthly_data[month_key] = {"expenses": 0.0, "income": 0.0}

        for col, val in r.items():
            if col == "Date" or col in SKIP_CATEGORIES:
                continue

            amt = to_float(val)
            if amt is None or amt == 0:
                continue 
            monthly_data[month_key]["expenses"] += amt

            category = col.strip()
            comment = ""
            to_insert.append((category, amt, date_str, comment))

        income_val = r.get("Income")
        if income_val:
            monthly_data[month_key]["income"] += to_float(income_val)

    for month, data in monthly_data.items():
        to_insert_summary.append((
            month,
            data["expenses"],
            data["income"],
            total_savings
        ))


cur.executemany('''
    INSERT OR IGNORE INTO csv_data (category, amount, expense_date, comment)
    VALUES (?, ?, ?, ?)
    ''', to_insert)

cur.executemany('''  
    INSERT OR IGNORE INTO budget_summary (month, monthly_total_expense, gross_month_income, total_savings)
    VALUES (?, ?, ?, ?)
    ''', to_insert_summary)

connect.commit()
print(f'Attempted to insert {len(to_insert)} rows — duplicates ignored.')
connect.close()

print(to_insert_summary)