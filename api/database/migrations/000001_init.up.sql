CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    off_budget BOOLEAN NOT NULL,
    category TEXT NOT NULL,
    adapter TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    notes VARCHAR(512),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    group_id UUID,
    name VARCHAR(255),
    notes VARCHAR(512),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    UNIQUE (group_id, name)
);

CREATE TABLE payees (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    account_id UUID,
    category_id UUID,
    payee_id UUID,
    name VARCHAR(255),
    credit DOUBLE PRECISION,
    debit DOUBLE PRECISION,
    notes VARCHAR(512),
    cleared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    FOREIGN KEY (payee_id) REFERENCES payees(id) ON DELETE SET NULL,
    UNIQUE (account_id, notes, credit, debit, cleared_at)
);

CREATE TABLE budgets (
    id UUID PRIMARY KEY,
    category_id UUID,
    year INTEGER,
    month INTEGER,
    budgeted DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE (year, month, category_id)
);
