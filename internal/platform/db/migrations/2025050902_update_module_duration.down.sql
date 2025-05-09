-- Rollback: Drop triggers and function
DROP TRIGGER IF EXISTS trg_update_module_duration_after_upsert ON module_items;
DROP TRIGGER IF EXISTS trg_update_module_duration_after_delete ON module_items;
DROP TRIGGER IF EXISTS trg_update_module_duration_after_soft_delete ON module_items;
DROP FUNCTION IF EXISTS update_module_duration;

-- Recreate original function without deleted_at condition
CREATE OR REPLACE FUNCTION update_module_duration()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE modules
  SET duration = (
    SELECT COALESCE(SUM(required_time), 0)
    FROM module_items
    WHERE module_id = COALESCE(NEW.module_id, OLD.module_id)
  )
  WHERE id = COALESCE(NEW.module_id, OLD.module_id);

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate original triggers
CREATE TRIGGER trg_update_module_duration_after_upsert
AFTER INSERT OR UPDATE OF required_time, module_id
ON module_items
FOR EACH ROW
EXECUTE FUNCTION update_module_duration();

CREATE TRIGGER trg_update_module_duration_after_delete
AFTER DELETE
ON module_items
FOR EACH ROW
EXECUTE FUNCTION update_module_duration();
