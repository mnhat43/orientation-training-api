ALTER TABLE
    quiz_submissions
ADD
    CONSTRAINT fk_quiz_submissions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE
    quiz_submissions
ADD
    CONSTRAINT fk_quiz_submissions_quiz_id FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE;

ALTER TABLE
    quiz_submissions
ADD
    CONSTRAINT fk_quiz_submissions_quiz_question_id FOREIGN KEY (quiz_question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE;