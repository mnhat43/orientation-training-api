ALTER TABLE module_items ADD CONSTRAINT fk_module_items_module_id FOREIGN KEY (module_id) REFERENCES modules(id);
