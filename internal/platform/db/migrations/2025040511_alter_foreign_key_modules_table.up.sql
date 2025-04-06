ALTER TABLE modules ADD CONSTRAINT fk_modules_course_id FOREIGN KEY (course_id) REFERENCES courses(id);
