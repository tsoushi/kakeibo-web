-- record_type に初期データを挿入（支出、収入、振替）
INSERT IGNORE INTO record_type (name) VALUES
    ('EXPENSE'),
    ('INCOME'),
    ('TRANSFER');
