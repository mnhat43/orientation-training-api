-- Create function to update duration in modules
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

-- Trigger: After INSERT or UPDATE
CREATE TRIGGER trg_update_module_duration_after_upsert
AFTER INSERT OR UPDATE OF required_time, module_id
ON module_items
FOR EACH ROW
EXECUTE FUNCTION update_module_duration();

-- Trigger: After DELETE
CREATE TRIGGER trg_update_module_duration_after_delete
AFTER DELETE
ON module_items
FOR EACH ROW
EXECUTE FUNCTION update_module_duration();