-- Drop old triggers and function if they exist
DROP TRIGGER IF EXISTS trg_update_course_duration_after_upsert ON modules;
DROP TRIGGER IF EXISTS trg_update_course_duration_after_delete ON modules;
DROP TRIGGER IF EXISTS trg_update_course_duration_after_soft_delete ON modules;
DROP FUNCTION IF EXISTS update_course_duration;

-- Recreate function with deleted_at condition
CREATE OR REPLACE FUNCTION update_course_duration()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE courses
  SET duration = (
    SELECT COALESCE(SUM(duration), 0)
    FROM modules
    WHERE course_id = COALESCE(NEW.course_id, OLD.course_id)
      AND deleted_at IS NULL
  )
  WHERE id = COALESCE(NEW.course_id, OLD.course_id);

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers: Insert / Update / Delete / Soft-delete
CREATE TRIGGER trg_update_course_duration_after_upsert
AFTER INSERT OR UPDATE OF duration, course_id
ON modules
FOR EACH ROW
EXECUTE FUNCTION update_course_duration();

CREATE TRIGGER trg_update_course_duration_after_delete
AFTER DELETE
ON modules
FOR EACH ROW
EXECUTE FUNCTION update_course_duration();

CREATE TRIGGER trg_update_course_duration_after_soft_delete
AFTER UPDATE OF deleted_at
ON modules
FOR EACH ROW
EXECUTE FUNCTION update_course_duration();
