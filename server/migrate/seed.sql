-- record_type に初期データを挿入（支出、収入、振替）
INSERT IGNORE INTO record_type (id, name) VALUES
    (1, '支出'),
    (2, '収入'),
    (3, '振替');
