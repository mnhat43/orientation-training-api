ALTER TABLE
    quiz_submissions DROP CONSTRAINT IF EXISTS fk_quiz_submissions_user_id;

ALTER TABLE
    quiz_submissions DROP CONSTRAINT IF EXISTS fk_quiz_submissions_quiz_id;

ALTER TABLE
    quiz_submissions DROP CONSTRAINT IF EXISTS fk_quiz_submissions_quiz_question_id;