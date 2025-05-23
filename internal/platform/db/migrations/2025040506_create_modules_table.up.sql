CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
    course_id INT NOT NULL,
    position INT DEFAULT 1,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);