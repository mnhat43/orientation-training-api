ALTER TABLE
    module_items
ADD
    CONSTRAINT fk_module_items_quiz_id FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE;