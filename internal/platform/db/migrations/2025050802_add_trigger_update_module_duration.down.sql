-- Drop triggers
DROP TRIGGER IF EXISTS trg_update_module_duration_after_upsert ON module_items;
DROP TRIGGER IF EXISTS trg_update_module_duration_after_delete ON module_items;

-- Drop function
DROP FUNCTION IF EXISTS update_module_duration;