ALTER TABLE
    user_progresses DROP CONSTRAINT IF EXISTS fk_user_progresses_user_id;

ALTER TABLE
    user_progresses DROP CONSTRAINT IF EXISTS fk_user_progresses_course_id;