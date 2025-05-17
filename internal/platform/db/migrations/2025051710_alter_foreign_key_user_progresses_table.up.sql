ALTER TABLE
    user_progresses
ADD
    CONSTRAINT fk_user_progresses_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE
    user_progresses
ADD
    CONSTRAINT fk_user_progresses_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;