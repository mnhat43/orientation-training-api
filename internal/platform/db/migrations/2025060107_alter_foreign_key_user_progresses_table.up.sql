ALTER TABLE
    user_progresses
ADD
    CONSTRAINT fk_user_progresses_reviewed_by FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE CASCADE;