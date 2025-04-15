CREATE TABLE IF NOT EXISTS user_progresses (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    course_id INT NOT NULL,
    module_position INT NOT NULL,
    module_item_position INT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);
