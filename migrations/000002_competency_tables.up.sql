-- Competency Schema Migration
-- Creates all competency domain tables in CoreSchema.
-- Mirrors .NET EF Core CompetencyCoreDbContext DbSets.

-- ============================================================
-- COMPETENCY CORE TABLES (CoreSchema)
-- ============================================================

CREATE TABLE IF NOT EXISTS "CoreSchema".ratings (
    rating_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    value INT NOT NULL,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".review_types (
    review_type_id SERIAL PRIMARY KEY,
    review_type_name TEXT NOT NULL UNIQUE,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".training_types (
    training_type_id SERIAL PRIMARY KEY,
    training_type_name TEXT NOT NULL,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competency_categories (
    competency_category_id SERIAL PRIMARY KEY,
    category_name TEXT NOT NULL UNIQUE,
    is_technical BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competencies (
    competency_id SERIAL PRIMARY KEY,
    competency_category_id INT NOT NULL REFERENCES "CoreSchema".competency_categories(competency_category_id),
    competency_name TEXT NOT NULL UNIQUE,
    description TEXT,
    -- BaseWorkFlowData fields
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT,
    approved_by TEXT,
    date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE,
    is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT,
    rejection_reason TEXT,
    date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competency_category_gradings (
    competency_category_grading_id SERIAL PRIMARY KEY,
    competency_category_id INT NOT NULL REFERENCES "CoreSchema".competency_categories(competency_category_id),
    review_type_id INT NOT NULL REFERENCES "CoreSchema".review_types(review_type_id),
    weight_percentage DECIMAL(18,2) NOT NULL,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competency_rating_definitions (
    competency_rating_definition_id SERIAL PRIMARY KEY,
    competency_id INT NOT NULL REFERENCES "CoreSchema".competencies(competency_id),
    rating_id INT NOT NULL REFERENCES "CoreSchema".ratings(rating_id),
    definition TEXT NOT NULL,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

-- ============================================================
-- JOB ROLE & GRADE TABLES (CoreSchema)
-- ============================================================

CREATE TABLE IF NOT EXISTS "CoreSchema".job_roles (
    job_role_id SERIAL PRIMARY KEY,
    job_role_name TEXT NOT NULL,
    description TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".job_grades (
    job_grade_id SERIAL PRIMARY KEY,
    grade_code TEXT NOT NULL UNIQUE,
    grade_name TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".job_grade_groups (
    job_grade_group_id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL UNIQUE,
    "order" INT NOT NULL,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".assign_job_grade_groups (
    assign_job_grade_group_id SERIAL PRIMARY KEY,
    job_grade_group_id INT NOT NULL REFERENCES "CoreSchema".job_grade_groups(job_grade_group_id),
    job_grade_id INT NOT NULL REFERENCES "CoreSchema".job_grades(job_grade_id),
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".office_job_roles (
    office_job_role_id SERIAL PRIMARY KEY,
    office_id INT NOT NULL REFERENCES "CoreSchema".offices(office_id),
    job_role_id INT NOT NULL REFERENCES "CoreSchema".job_roles(job_role_id),
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".job_role_competencies (
    job_role_competency_id SERIAL PRIMARY KEY,
    office_id INT NOT NULL REFERENCES "CoreSchema".offices(office_id),
    job_role_id INT NOT NULL REFERENCES "CoreSchema".job_roles(job_role_id),
    competency_id INT NOT NULL REFERENCES "CoreSchema".competencies(competency_id),
    rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT,
    UNIQUE (job_role_competency_id, office_id, job_role_id, competency_id)
);

CREATE TABLE IF NOT EXISTS "CoreSchema".behavioral_competencies (
    behavioral_competency_id SERIAL PRIMARY KEY,
    competency_id INT NOT NULL REFERENCES "CoreSchema".competencies(competency_id),
    job_grade_group_id INT NOT NULL REFERENCES "CoreSchema".job_grade_groups(job_grade_group_id),
    rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT,
    UNIQUE (behavioral_competency_id, competency_id, job_grade_group_id)
);

CREATE TABLE IF NOT EXISTS "CoreSchema".job_role_grades (
    job_role_grade_id SERIAL PRIMARY KEY,
    job_role_id INT NOT NULL REFERENCES "CoreSchema".job_roles(job_role_id),
    grade_id TEXT,
    grade_name TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

-- ============================================================
-- REVIEW PERIOD & REVIEW TABLES (CoreSchema)
-- ============================================================

CREATE TABLE IF NOT EXISTS "CoreSchema".review_periods (
    review_period_id SERIAL PRIMARY KEY,
    bank_year_id INT NOT NULL REFERENCES "CoreSchema".bank_years(bank_year_id),
    name TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    -- BaseWorkFlowData fields
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT,
    approved_by TEXT,
    date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE,
    is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT,
    rejection_reason TEXT,
    date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competency_reviews (
    competency_review_id SERIAL PRIMARY KEY,
    employee_number TEXT NOT NULL,
    review_period_id INT NOT NULL REFERENCES "CoreSchema".review_periods(review_period_id),
    competency_id INT NOT NULL REFERENCES "CoreSchema".competencies(competency_id),
    review_type_id INT REFERENCES "CoreSchema".review_types(review_type_id),
    expected_rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    actual_rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    actual_rating_value INT DEFAULT 0,
    review_date TIMESTAMPTZ,
    reviewer_id TEXT,
    reviewer_name TEXT,
    employee_name TEXT,
    employee_grade TEXT,
    employee_department TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".competency_review_profiles (
    competency_review_profile_id SERIAL PRIMARY KEY,
    review_period_id INT NOT NULL REFERENCES "CoreSchema".review_periods(review_period_id),
    average_rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    expected_rating_id INT REFERENCES "CoreSchema".ratings(rating_id),
    average_score DECIMAL(18,2) DEFAULT 0,
    competency_gap INT DEFAULT 0,
    have_gap BOOLEAN DEFAULT FALSE,
    employee_number TEXT NOT NULL,
    employee_full_name TEXT,
    competency_id INT REFERENCES "CoreSchema".competencies(competency_id),
    competency_name TEXT,
    competency_category INT,
    competency_category_name TEXT,
    office_id TEXT,
    office_name TEXT,
    division_id TEXT,
    division_name TEXT,
    department_id TEXT,
    department_name TEXT,
    job_role_id TEXT,
    job_role_name TEXT,
    grade_name TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE IF NOT EXISTS "CoreSchema".development_plans (
    development_plan_id SERIAL PRIMARY KEY,
    competency_review_profile_id INT REFERENCES "CoreSchema".competency_review_profiles(competency_review_profile_id),
    training_type_name TEXT,
    activity TEXT NOT NULL,
    employee_number TEXT NOT NULL,
    target_date TIMESTAMPTZ NOT NULL,
    completion_date TIMESTAMPTZ,
    task_status TEXT,
    learning_resource TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

-- ============================================================
-- STAFF JOB ROLES (CoreSchema)
-- ============================================================

CREATE TABLE IF NOT EXISTS "CoreSchema".staff_job_roles (
    staff_job_role_id SERIAL PRIMARY KEY,
    employee_id TEXT NOT NULL,
    full_name TEXT,
    department_id INT REFERENCES "CoreSchema".departments(department_id),
    division_id INT REFERENCES "CoreSchema".divisions(division_id),
    office_id INT REFERENCES "CoreSchema".offices(office_id),
    supervisor_id TEXT,
    job_role_id INT REFERENCES "CoreSchema".job_roles(job_role_id),
    job_role_name TEXT,
    soa_status BOOLEAN DEFAULT FALSE,
    soa_response TEXT,
    -- HrdWorkFlowData fields
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT,
    approved_by TEXT,
    date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE,
    is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT,
    rejection_reason TEXT,
    date_rejected TIMESTAMPTZ,
    hrd_approved_by TEXT,
    hrd_date_approved TIMESTAMPTZ,
    hrd_is_approved BOOLEAN DEFAULT FALSE,
    hrd_is_rejected BOOLEAN DEFAULT FALSE,
    hrd_rejected_by TEXT,
    hrd_rejection_reason TEXT,
    hrd_date_rejected TIMESTAMPTZ
);
