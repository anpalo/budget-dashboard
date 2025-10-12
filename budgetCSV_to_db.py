import csv, sqlite3, re
from datetime import datetime

DB = "budget.db"
CSV = "CurrentBudget.csv"


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


connect = sqlite3.connect(DB)
cur = connect.cursor() # to use SQL commands eg. SELECT INSERT UPDATE DELETE


SKIP_CATEGORIES = {"Daily Sum","Monthly Sum","Adjusted Monthly Sum","Income","Net Month Income","Yearly Savings Total"}

# https://docs.python.org/3/library/csv.html
with open(CSV, newline = '', encoding='utf-8') as f:
    reader = csv.DictReader(f)
    to_insert = []
    for r in reader:
        if r.get("Date"):
            date = r.get("Date").strip()
        else:
            date = ""
        if not is_date(date):
            continue
        for col, val in r.items():
            if col == "Date" or col in SKIP_CATEGORIES:
                continue
            amt = to_float(val)
            if amt is None or amt == 0:
                continue
            category = col.strip()
            comment = ""
            to_insert.append((category, amt, date, comment) )

cur.executemany('''
INSERT OR IGNORE INTO csv_data (category, amount, expense_date, comment)
VALUES (?, ?, ?, ?)
''', to_insert)

connect.commit()
print(f'Attempted to insert {len(to_insert)} rows — duplicates ignored.')
connect.close()