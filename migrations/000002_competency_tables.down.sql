-- Reverse competency tables migration (drop in dependency order)

DROP TABLE IF EXISTS "CoreSchema".staff_job_roles;
DROP TABLE IF EXISTS "CoreSchema".development_plans;
DROP TABLE IF EXISTS "CoreSchema".competency_review_profiles;
DROP TABLE IF EXISTS "CoreSchema".competency_reviews;
DROP TABLE IF EXISTS "CoreSchema".review_periods;
DROP TABLE IF EXISTS "CoreSchema".job_role_grades;
DROP TABLE IF EXISTS "CoreSchema".behavioral_competencies;
DROP TABLE IF EXISTS "CoreSchema".job_role_competencies;
DROP TABLE IF EXISTS "CoreSchema".office_job_roles;
DROP TABLE IF EXISTS "CoreSchema".assign_job_grade_groups;
DROP TABLE IF EXISTS "CoreSchema".job_grade_groups;
DROP TABLE IF EXISTS "CoreSchema".job_grades;
DROP TABLE IF EXISTS "CoreSchema".job_roles;
DROP TABLE IF EXISTS "CoreSchema".competency_rating_definitions;
DROP TABLE IF EXISTS "CoreSchema".competency_category_gradings;
DROP TABLE IF EXISTS "CoreSchema".competencies;
DROP TABLE IF EXISTS "CoreSchema".competency_categories;
DROP TABLE IF EXISTS "CoreSchema".training_types;
DROP TABLE IF EXISTS "CoreSchema".review_types;
DROP TABLE IF EXISTS "CoreSchema".ratings;
