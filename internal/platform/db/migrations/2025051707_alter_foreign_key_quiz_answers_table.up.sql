ALTER TABLE
    quiz_answers
ADD
    CONSTRAINT fk_quiz_answers_quiz_question_id FOREIGN KEY (quiz_question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE;