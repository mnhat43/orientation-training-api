ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role_id; ALTER TABLE users ADD CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES user_roles(id);
