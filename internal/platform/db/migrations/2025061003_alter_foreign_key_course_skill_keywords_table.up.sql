ALTER TABLE course_skill_keywords ADD CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE;


ALTER TABLE course_skill_keywords ADD CONSTRAINT fk_skill_keyword FOREIGN KEY (skill_keyword_id) REFERENCES skill_keywords (id) ON DELETE CASCADE;