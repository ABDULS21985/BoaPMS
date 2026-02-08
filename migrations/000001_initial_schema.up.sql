-- PMS Initial Schema Migration
-- Creates all core schemas and tables for the Performance Management System.
-- Mirrors the .NET EF Core model with PostgreSQL-native types.

-- ============================================================
-- SCHEMAS
-- ============================================================
CREATE SCHEMA IF NOT EXISTS "CoreSchema";
CREATE SCHEMA IF NOT EXISTS pms;
CREATE SCHEMA IF NOT EXISTS pmsaudit;

-- ============================================================
-- ORGANOGRAM (CoreSchema)
-- ============================================================

CREATE TABLE "CoreSchema".directorates (
    directorate_id SERIAL PRIMARY KEY,
    directorate_name TEXT NOT NULL,
    directorate_code TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE "CoreSchema".departments (
    department_id SERIAL PRIMARY KEY,
    directorate_id INT REFERENCES "CoreSchema".directorates(directorate_id),
    department_name TEXT NOT NULL,
    department_code TEXT,
    is_branch BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE "CoreSchema".divisions (
    division_id SERIAL PRIMARY KEY,
    department_id INT NOT NULL REFERENCES "CoreSchema".departments(department_id),
    division_name TEXT NOT NULL,
    division_code TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE "CoreSchema".offices (
    office_id SERIAL PRIMARY KEY,
    division_id INT NOT NULL REFERENCES "CoreSchema".divisions(division_id),
    office_name TEXT NOT NULL,
    office_code TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

-- ============================================================
-- IDENTITY (CoreSchema)
-- ============================================================

CREATE TABLE "CoreSchema".asp_net_roles (
    id VARCHAR(450) PRIMARY KEY,
    name VARCHAR(256),
    normalized_name VARCHAR(256) UNIQUE,
    concurrency_stamp TEXT
);

CREATE TABLE "CoreSchema".asp_net_users (
    id VARCHAR(450) PRIMARY KEY,
    user_name VARCHAR(256) UNIQUE,
    normalized_user_name VARCHAR(256) UNIQUE,
    email VARCHAR(256),
    normalized_email VARCHAR(256),
    email_confirmed BOOLEAN DEFAULT FALSE,
    password_hash TEXT,
    security_stamp TEXT,
    concurrency_stamp TEXT,
    phone_number TEXT,
    phone_number_confirmed BOOLEAN DEFAULT FALSE,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    lockout_end TIMESTAMPTZ,
    lockout_enabled BOOLEAN DEFAULT TRUE,
    access_failed_count INT DEFAULT 0,
    first_name TEXT,
    last_name TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE "CoreSchema".asp_net_user_roles (
    user_id VARCHAR(450) NOT NULL REFERENCES "CoreSchema".asp_net_users(id),
    role_id VARCHAR(450) NOT NULL REFERENCES "CoreSchema".asp_net_roles(id),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE "CoreSchema".bank_years (
    bank_year_id SERIAL PRIMARY KEY,
    year_name VARCHAR(10) NOT NULL UNIQUE,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE "CoreSchema".permissions (
    permission_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

CREATE TABLE "CoreSchema".role_permissions (
    role_permission_id SERIAL PRIMARY KEY,
    permission_id INT NOT NULL REFERENCES "CoreSchema".permissions(permission_id),
    role_id VARCHAR(450) NOT NULL REFERENCES "CoreSchema".asp_net_roles(id),
    created_by VARCHAR(75) DEFAULT 'SYSTEM',
    date_created TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    status VARCHAR(25),
    soft_deleted BOOLEAN DEFAULT FALSE,
    date_updated TIMESTAMPTZ,
    updated_by TEXT
);

-- ============================================================
-- PMS CORE TABLES (pms schema)
-- ============================================================

CREATE TABLE pms.strategies (
    strategy_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    smd_reference_code TEXT,
    description TEXT,
    bank_year_id INT NOT NULL REFERENCES "CoreSchema".bank_years(bank_year_id),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    file_image TEXT,
    -- BaseWorkFlow fields
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT,
    date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE,
    is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT,
    rejection_reason TEXT,
    date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.strategic_themes (
    strategic_theme_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    strategy_id TEXT NOT NULL REFERENCES pms.strategies(strategy_id),
    file_image TEXT,
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.performance_review_periods (
    period_id TEXT PRIMARY KEY,
    year INT NOT NULL,
    range INT,
    range_value INT,
    name TEXT NOT NULL,
    description TEXT,
    short_name TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    allow_objective_planning BOOLEAN DEFAULT FALSE,
    allow_work_product_planning BOOLEAN DEFAULT FALSE,
    allow_work_product_evaluation BOOLEAN DEFAULT FALSE,
    max_points DECIMAL(18,2) DEFAULT 250,
    min_no_of_objectives INT DEFAULT 1,
    max_no_of_objectives INT,
    strategy_id TEXT REFERENCES pms.strategies(strategy_id),
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.objective_categories (
    objective_category_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.enterprise_objectives (
    enterprise_objective_id TEXT PRIMARY KEY,
    type INT DEFAULT 1,
    enterprise_objectives_category_id TEXT NOT NULL REFERENCES pms.objective_categories(objective_category_id),
    strategic_theme_id TEXT REFERENCES pms.strategic_themes(strategic_theme_id),
    strategy_id TEXT NOT NULL REFERENCES pms.strategies(strategy_id),
    name TEXT NOT NULL,
    smd_reference_code TEXT,
    description TEXT,
    kpi TEXT NOT NULL,
    target TEXT,
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.department_objectives (
    department_objective_id TEXT PRIMARY KEY,
    department_id INT NOT NULL REFERENCES "CoreSchema".departments(department_id),
    enterprise_objective_id TEXT NOT NULL REFERENCES pms.enterprise_objectives(enterprise_objective_id),
    name TEXT NOT NULL, smd_reference_code TEXT, description TEXT, kpi TEXT NOT NULL, target TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.division_objectives (
    division_objective_id TEXT PRIMARY KEY,
    division_id INT NOT NULL REFERENCES "CoreSchema".divisions(division_id),
    department_objective_id TEXT NOT NULL REFERENCES pms.department_objectives(department_objective_id),
    name TEXT NOT NULL, smd_reference_code TEXT, description TEXT, kpi TEXT NOT NULL, target TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.office_objectives (
    office_objective_id TEXT PRIMARY KEY,
    office_id INT NOT NULL REFERENCES "CoreSchema".offices(office_id),
    division_objective_id TEXT NOT NULL REFERENCES pms.division_objectives(division_objective_id),
    job_grade_group_id INT NOT NULL,
    name TEXT NOT NULL, smd_reference_code TEXT, description TEXT, kpi TEXT NOT NULL, target TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE pms.work_products (
    work_product_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    max_point DECIMAL(18,2) NOT NULL,
    work_product_type INT,
    is_self_created BOOLEAN DEFAULT FALSE,
    staff_id TEXT NOT NULL,
    acceptance_comment TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    deliverables TEXT,
    final_score DECIMAL(18,2) DEFAULT 0,
    no_returned INT DEFAULT 0,
    completion_date TIMESTAMPTZ,
    approver_comment TEXT,
    re_evaluation_re_initiated BOOLEAN DEFAULT FALSE,
    remark TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);
CREATE INDEX idx_work_products_staff ON pms.work_products(staff_id);

CREATE TABLE pms.period_scores (
    period_score_id TEXT PRIMARY KEY,
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    staff_id TEXT NOT NULL,
    final_score DECIMAL(18,2) DEFAULT 0,
    score_percentage DECIMAL(18,2) DEFAULT 0,
    final_grade INT NOT NULL,
    end_date TIMESTAMPTZ,
    office_id INT REFERENCES "CoreSchema".offices(office_id),
    min_no_of_objectives INT,
    max_no_of_objectives INT,
    strategy_id TEXT REFERENCES pms.strategies(strategy_id),
    staff_grade TEXT,
    location_id TEXT,
    hrd_deducted_points DECIMAL(18,2) DEFAULT 0,
    is_under_performing BOOLEAN DEFAULT FALSE,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX idx_period_scores_staff ON pms.period_scores(staff_id);

-- ============================================================
-- AUDIT TABLES (pmsaudit schema)
-- ============================================================

CREATE TABLE pmsaudit.audit_logs (
    id SERIAL PRIMARY KEY,
    user_name TEXT,
    audit_event_date_utc TIMESTAMPTZ NOT NULL,
    audit_event_type INT NOT NULL,
    table_name TEXT,
    record_id TEXT,
    field_name TEXT,
    original_value TEXT,
    new_value TEXT,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX idx_audit_logs_table ON pmsaudit.audit_logs(table_name);
CREATE INDEX idx_audit_logs_date ON pmsaudit.audit_logs(audit_event_date_utc);

CREATE TABLE pmsaudit.auditable_entities (
    entity_name TEXT PRIMARY KEY,
    enable_audit BOOLEAN DEFAULT TRUE,
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE pmsaudit.auditable_attributes (
    id SERIAL PRIMARY KEY,
    auditable_entity_id INT,
    attribute_name TEXT,
    enable_audit BOOLEAN DEFAULT TRUE,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE pms.sequence_numbers (
    description TEXT PRIMARY KEY,
    sequence_number_type INT,
    prefix TEXT,
    next_number BIGINT DEFAULT 1,
    use_prefix BOOLEAN DEFAULT FALSE,
    id SERIAL,
    record_status TEXT DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE,
    status TEXT,
    updated_at TIMESTAMPTZ,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE pms.settings (
    setting_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT,
    type TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE pms.pms_configurations (
    pms_configuration_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT,
    type TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
