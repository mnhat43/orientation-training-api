-- Remove the foreign key constraint from the app_feedbacks table
ALTER TABLE app_feedbacks
DROP CONSTRAINT IF EXISTS fk_app_feedbacks_user_id;