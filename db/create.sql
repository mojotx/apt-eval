CREATE TABLE IF NOT EXISTS apartments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    address TEXT NOT NULL,
    visit_date TIMESTAMP,
    notes TEXT,
    rating INTEGER,
    price REAL,
    floor INTEGER,
    is_gated BOOLEAN DEFAULT 0,
    has_garage BOOLEAN DEFAULT 0,
    has_laundry BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
