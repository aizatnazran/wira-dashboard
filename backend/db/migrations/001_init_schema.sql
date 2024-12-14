-- Create accounts table (combining with users table)
CREATE TABLE IF NOT EXISTS accounts (
    acc_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts(username);
CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

-- Create races table
CREATE TABLE IF NOT EXISTS races (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NOT NULL
);

-- Create classes table
CREATE TABLE IF NOT EXISTS classes (
    id SERIAL PRIMARY KEY,
    race_id INTEGER REFERENCES races(id),
    name VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    combat_type VARCHAR(50) NOT NULL,
    damage INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    difficulty INTEGER NOT NULL,
    speed INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_classes_race_id ON classes(race_id);

-- Create characters table
CREATE TABLE IF NOT EXISTS characters (
    char_id SERIAL PRIMARY KEY,
    acc_id INTEGER REFERENCES accounts(acc_id),
    class_id INTEGER REFERENCES classes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create scores table
CREATE TABLE IF NOT EXISTS scores (
    score_id SERIAL PRIMARY KEY,
    char_id INTEGER REFERENCES characters(char_id),
    reward_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_characters_acc_id ON characters(acc_id);
CREATE INDEX IF NOT EXISTS idx_characters_class_id ON characters(class_id);
CREATE INDEX IF NOT EXISTS idx_scores_char_id ON scores(char_id);
CREATE INDEX IF NOT EXISTS idx_scores_reward_score ON scores(reward_score DESC);

-- Create view for rankings
CREATE OR REPLACE VIEW character_rankings AS
SELECT 
    c.char_id,
    a.username,
    cl.name as class_name,
    s.reward_score,
    RANK() OVER (PARTITION BY c.class_id ORDER BY s.reward_score DESC) as rank
FROM characters c
JOIN accounts a ON c.acc_id = a.acc_id
JOIN classes cl ON c.class_id = cl.id
JOIN scores s ON c.char_id = s.char_id;
