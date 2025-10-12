CREATE TABLE csv_data (
id INTEGER PRIMARY KEY,
category TEXT,
amount REAL, 
expense_date TEXT,
comment TEXT,
UNIQUE(category, amount, expense_date)

);

CREATE TABLE budget_summary (
id INTEGER PRIMARY KEY,
month TEXT UNIQUE NOT NULL,
monthly_total_expense REAL, 
gross_month_income REAL, 
total_savings REAL
);


CREATE TABLE settings (
id INTEGER PRIMARY KEY, 
widget_name TEXT,
bg_color TEXT,
font_color TEXT,
border_color TEXT
);

CREATE TABLE stocks (
id INTEGER PRIMARY KEY,
symbol TEXT, 
trade_date TEXT,
open_price REAL,
close_price REAL
);

CREATE TABLE favorite_stocks (
id INTEGER PRIMARY KEY,
symbol_text TEXT UNIQUE,
added_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
