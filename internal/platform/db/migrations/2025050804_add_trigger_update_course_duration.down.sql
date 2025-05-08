-- Drop triggers
DROP TRIGGER IF EXISTS trg_update_course_duration_after_upsert ON modules;
DROP TRIGGER IF EXISTS trg_update_course_duration_after_delete ON modules;

-- Drop function
DROP FUNCTION IF EXISTS update_course_duration;