ALTER TABLE
    quiz_submissions
ADD
    CONSTRAINT fk_quiz_submissions_reviewed_by FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE CASCADE;