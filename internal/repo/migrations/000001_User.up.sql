

CREATE TABLE users (
    id SERIAL PRIMARY KEY ,
    username VARCHAR(50) UNIQUE NOT NULL,
    password  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);


-- Таблица для хранения сессий/токенов (опционально)
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY ,
    user_id SERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);