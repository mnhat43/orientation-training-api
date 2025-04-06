ALTER TABLE user_courses ADD CONSTRAINT fk_user_courses_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_courses ADD CONSTRAINT fk_user_courses_course_id FOREIGN KEY (course_id) REFERENCES courses(id);
