-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    acc_id INTEGER REFERENCES accounts(acc_id),
    session_metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expiry_datetime TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_acc_id ON sessions(acc_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry_datetime);
