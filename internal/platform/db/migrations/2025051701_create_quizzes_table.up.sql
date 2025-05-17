CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    difficulty INT NOT NULL,
    total_score FLOAT NOT NULL,
    time_limit INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);