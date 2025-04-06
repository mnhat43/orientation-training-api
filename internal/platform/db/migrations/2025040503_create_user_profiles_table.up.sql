CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    avatar VARCHAR(100),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    birthday DATE,
    phone_number VARCHAR(20),
    personal_email VARCHAR(100),
    company_joined_date DATE,
    introduce TEXT,
    gender INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);