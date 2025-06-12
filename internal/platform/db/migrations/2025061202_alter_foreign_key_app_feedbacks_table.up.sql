-- Add foreign key constraint to the user_id column in app_feedbacks table
ALTER TABLE app_feedbacks ADD CONSTRAINT fk_app_feedbacks_user_id FOREIGN KEY (user_id) REFERENCES users (id);