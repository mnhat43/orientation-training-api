ALTER TABLE courses ADD CONSTRAINT fk_courses_created_by FOREIGN KEY (created_by) REFERENCES users(id);
