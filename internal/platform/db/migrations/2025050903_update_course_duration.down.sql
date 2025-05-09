-- Rollback: Drop new triggers and function
DROP TRIGGER IF EXISTS trg_update_course_duration_after_upsert ON modules;
DROP TRIGGER IF EXISTS trg_update_course_duration_after_delete ON modules;
DROP TRIGGER IF EXISTS trg_update_course_duration_after_soft_delete ON modules;
DROP FUNCTION IF EXISTS update_course_duration;

-- Recreate original function without deleted_at filter
CREATE OR REPLACE FUNCTION update_course_duration()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE courses
  SET duration = (
    SELECT COALESCE(SUM(duration), 0)
    FROM modules
    WHERE course_id = COALESCE(NEW.course_id, OLD.course_id)
  )
  WHERE id = COALESCE(NEW.course_id, OLD.course_id);

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate original triggers (no soft-delete needed)
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
