ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users(id);
