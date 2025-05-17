ALTER TABLE
    quiz_questions
ADD
    CONSTRAINT fk_quiz_questions_quiz_id FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE;