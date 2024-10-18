CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);


CREATE TABLE runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target TEXT NOT NULL,
    html TEXT,
    accuracy REAL,
    wpm REAL,
    awpm REAL,
    time REAL
);

CREATE TABLE run_inputs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER,
    value TEXT NOT NULL,
    time REAL,
    FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
);
