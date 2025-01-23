-- Table to store categories
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- Table to store principles
CREATE TABLE principles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
);

-- Table to store relationships between principles
CREATE TABLE principle_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    principle_id INTEGER NOT NULL,
    related_id INTEGER NOT NULL,
    relation_type TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (principle_id) REFERENCES principles (id) ON DELETE CASCADE,
    FOREIGN KEY (related_id) REFERENCES principles (id) ON DELETE CASCADE
);
