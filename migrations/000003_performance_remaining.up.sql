-- Performance Remaining Tables Migration
-- Creates all remaining PMS tables not covered in 000001_initial_schema.

-- ============================================================
-- CATEGORY DEFINITIONS & OBJECTIVES (pms schema)
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.category_definitions (
    definition_id TEXT PRIMARY KEY,
    objective_category_id TEXT NOT NULL REFERENCES pms.objective_categories(objective_category_id),
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    weight DECIMAL(18,2) NOT NULL,
    max_no_objectives INT DEFAULT 0,
    max_no_work_product INT DEFAULT 0,
    max_points DECIMAL(18,2) DEFAULT 0,
    is_compulsory BOOLEAN DEFAULT FALSE,
    enforce_work_product_limit BOOLEAN DEFAULT FALSE,
    grade_group_id INT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.period_objectives (
    period_objective_id TEXT PRIMARY KEY,
    objective_id TEXT NOT NULL,
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.period_objective_evaluations (
    period_objective_evaluation_id TEXT PRIMARY KEY,
    total_outcome_score DECIMAL(18,2) DEFAULT 0,
    outcome_score DECIMAL(18,2) DEFAULT 0,
    period_objective_id TEXT NOT NULL REFERENCES pms.period_objectives(period_objective_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.period_objective_department_evaluations (
    period_objective_department_evaluation_id TEXT PRIMARY KEY,
    overall_outcome_scored DECIMAL(18,2) DEFAULT 0,
    allocated_outcome DECIMAL(18,2) DEFAULT 0,
    outcome_score DECIMAL(18,2) DEFAULT 0,
    department_id INT NOT NULL REFERENCES "CoreSchema".departments(department_id),
    period_objective_id TEXT NOT NULL REFERENCES pms.period_objectives(period_objective_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

-- ============================================================
-- INDIVIDUAL PLANNED OBJECTIVES
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.review_period_individual_planned_objectives (
    planned_objective_id TEXT PRIMARY KEY,
    objective_id TEXT NOT NULL,
    staff_id TEXT NOT NULL,
    objective_level INT DEFAULT 3,
    staff_job_role TEXT,
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    no_returned INT DEFAULT 0,
    remark TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_planned_obj_staff ON pms.review_period_individual_planned_objectives(staff_id);

-- ============================================================
-- WORK PRODUCT SUPPORTING TABLES
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.work_product_tasks (
    work_product_task_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    completion_date TIMESTAMPTZ,
    work_product_id TEXT NOT NULL REFERENCES pms.work_products(work_product_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.work_product_evaluations (
    work_product_evaluation_id TEXT PRIMARY KEY,
    work_product_id TEXT NOT NULL REFERENCES pms.work_products(work_product_id),
    timeliness DECIMAL(18,2) DEFAULT 0,
    timeliness_evaluation_option_id TEXT,
    quality DECIMAL(18,2) DEFAULT 0,
    quality_evaluation_option_id TEXT,
    output DECIMAL(18,2) DEFAULT 0,
    output_evaluation_option_id TEXT,
    outcome DECIMAL(18,2) DEFAULT 0,
    evaluator_staff_id TEXT,
    is_re_evaluated BOOLEAN DEFAULT FALSE,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.evaluation_options (
    evaluation_option_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    record_status INT,
    score DECIMAL(18,2) NOT NULL,
    evaluation_type INT NOT NULL,
    id SERIAL, created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.operational_objective_work_products (
    operational_objective_work_product_id TEXT PRIMARY KEY,
    work_product_id TEXT NOT NULL REFERENCES pms.work_products(work_product_id),
    work_product_definition_id TEXT,
    planned_objective_id TEXT REFERENCES pms.review_period_individual_planned_objectives(planned_objective_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms."WorkProductDefinitions" (
    work_product_definition_id TEXT PRIMARY KEY,
    reference_no TEXT,
    name TEXT NOT NULL,
    description TEXT,
    deliverables TEXT,
    objective_id TEXT,
    objective_level TEXT DEFAULT 'Office',
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms."CascadedWorkProducts" (
    cascaded_work_product_id TEXT PRIMARY KEY,
    smd_reference_code TEXT,
    name TEXT NOT NULL,
    description TEXT,
    objective_id TEXT,
    objective_level INT DEFAULT 3,
    staff_job_role TEXT,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

-- ============================================================
-- PROJECT & COMMITTEE TABLES
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.projects (
    project_id TEXT PRIMARY KEY,
    project_manager TEXT,
    name TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    deliverables TEXT,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    department_id INT REFERENCES "CoreSchema".departments(department_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.committees (
    committee_id TEXT PRIMARY KEY,
    chairperson TEXT,
    name TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    deliverables TEXT,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    department_id INT REFERENCES "CoreSchema".departments(department_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.project_members (
    project_member_id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL,
    project_id TEXT NOT NULL REFERENCES pms.projects(project_id),
    planned_objective_id TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.committee_members (
    committee_member_id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL,
    committee_id TEXT NOT NULL REFERENCES pms.committees(committee_id),
    planned_objective_id TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.project_objectives (
    project_objective_id TEXT PRIMARY KEY,
    objective_id TEXT NOT NULL,
    project_id TEXT NOT NULL REFERENCES pms.projects(project_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.committee_objectives (
    committee_objective_id TEXT PRIMARY KEY,
    objective_id TEXT NOT NULL,
    committee_id TEXT NOT NULL REFERENCES pms.committees(committee_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.project_work_products (
    project_work_product_id TEXT PRIMARY KEY,
    work_product_id TEXT NOT NULL REFERENCES pms.work_products(work_product_id),
    project_assigned_work_product_id TEXT,
    project_id TEXT NOT NULL REFERENCES pms.projects(project_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.committee_work_products (
    committee_work_product_id TEXT PRIMARY KEY,
    work_product_id TEXT NOT NULL REFERENCES pms.work_products(work_product_id),
    committee_assigned_work_product_id TEXT,
    committee_id TEXT NOT NULL REFERENCES pms.committees(committee_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.project_assigned_work_products (
    project_assigned_work_product_id TEXT PRIMARY KEY,
    work_product_definition_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    project_id TEXT NOT NULL REFERENCES pms.projects(project_id),
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    deliverables TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.committee_assigned_work_products (
    committee_assigned_work_product_id TEXT PRIMARY KEY,
    work_product_definition_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    committee_id TEXT NOT NULL REFERENCES pms.committees(committee_id),
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    deliverables TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

-- ============================================================
-- REVIEW PERIOD EXTENSIONS & 360 REVIEWS
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.review_period_extensions (
    review_period_extension_id TEXT PRIMARY KEY,
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    target_type INT NOT NULL,
    target_reference TEXT,
    description TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.review_period_360_reviews (
    review_period_360_review_id TEXT PRIMARY KEY,
    review_period_id TEXT NOT NULL REFERENCES pms.performance_review_periods(period_id),
    target_type INT NOT NULL,
    target_reference TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

-- ============================================================
-- FEEDBACK & QUESTIONNAIRES
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.feedback_request_logs (
    feedback_request_log_id TEXT PRIMARY KEY,
    feedback_request_type INT NOT NULL,
    reference_id TEXT,
    time_initiated TIMESTAMPTZ NOT NULL,
    assigned_staff_id TEXT NOT NULL,
    assigned_staff_name TEXT,
    request_owner_staff_id TEXT,
    request_owner_staff_name TEXT,
    time_completed TIMESTAMPTZ,
    request_owner_comment TEXT,
    assigned_staff_comment TEXT,
    has_sla BOOLEAN DEFAULT FALSE,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_feedback_req_ref ON pms.feedback_request_logs(reference_id);
CREATE INDEX IF NOT EXISTS idx_feedback_req_staff ON pms.feedback_request_logs(assigned_staff_id);

CREATE TABLE IF NOT EXISTS pms.pms_competencies (
    pms_competency_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    object_category_id TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.feedback_questionaires (
    feedback_questionaire_id TEXT PRIMARY KEY,
    question TEXT NOT NULL,
    description TEXT,
    pms_competency_id TEXT REFERENCES pms.pms_competencies(pms_competency_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE,
    approved_by TEXT, date_approved TIMESTAMPTZ,
    is_approved BOOLEAN DEFAULT FALSE, is_rejected BOOLEAN DEFAULT FALSE,
    rejected_by TEXT, rejection_reason TEXT, date_rejected TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pms.feedback_questionaire_options (
    feedback_questionaire_option_id TEXT PRIMARY KEY,
    option_statement TEXT NOT NULL,
    description TEXT,
    score DECIMAL(18,2) NOT NULL,
    question_id TEXT NOT NULL REFERENCES pms.feedback_questionaires(feedback_questionaire_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

-- ============================================================
-- COMPETENCY REVIEW FEEDBACK (360)
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.competency_review_feedbacks (
    competency_review_feedback_id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL,
    max_points DECIMAL(18,2) DEFAULT 0,
    final_score DECIMAL(18,2) DEFAULT 0,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_crf_staff ON pms.competency_review_feedbacks(staff_id);

CREATE TABLE IF NOT EXISTS pms.competency_reviewers (
    competency_reviewer_id TEXT PRIMARY KEY,
    review_staff_id TEXT NOT NULL,
    final_rating DECIMAL(18,2) DEFAULT 0,
    competency_review_feedback_id TEXT NOT NULL REFERENCES pms.competency_review_feedbacks(competency_review_feedback_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.competency_reviewer_ratings (
    competency_reviewer_rating_id TEXT PRIMARY KEY,
    pms_competency_id TEXT REFERENCES pms.pms_competencies(pms_competency_id),
    feedback_questionaire_option_id TEXT,
    rating DECIMAL(18,2) DEFAULT 0,
    competency_reviewer_id TEXT NOT NULL REFERENCES pms.competency_reviewers(competency_reviewer_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pms.competency_gap_closures (
    competency_gap_closure_id TEXT PRIMARY KEY,
    staff_id TEXT NOT NULL,
    max_points DECIMAL(18,2) DEFAULT 0,
    final_score DECIMAL(18,2) DEFAULT 0,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    objective_category_id TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_cgc_staff ON pms.competency_gap_closures(staff_id);

-- ============================================================
-- GRIEVANCES
-- ============================================================

CREATE TABLE IF NOT EXISTS pms.grievances (
    grievance_id TEXT PRIMARY KEY,
    grievance_type INT NOT NULL,
    review_period_id TEXT REFERENCES pms.performance_review_periods(period_id),
    subject_id TEXT,
    subject TEXT,
    description TEXT,
    respondent_comment TEXT,
    current_resolution_level INT,
    current_mediator_staff_id TEXT,
    complainant_staff_id TEXT NOT NULL,
    complainant_evidence_upload TEXT,
    respondent_staff_id TEXT NOT NULL,
    respondent_evidence_upload TEXT,
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_grievance_complainant ON pms.grievances(complainant_staff_id);
CREATE INDEX IF NOT EXISTS idx_grievance_respondent ON pms.grievances(respondent_staff_id);

CREATE TABLE IF NOT EXISTS pms.grievance_resolutions (
    grievance_resolution_id TEXT PRIMARY KEY,
    resolution_comment TEXT,
    resolution_level TEXT,
    level INT,
    mediator_staff_id TEXT,
    evidence_upload TEXT,
    respondent_feedback TEXT,
    complainant_feedback TEXT,
    complainant_remark INT,
    respondent_remark INT,
    grievance_id TEXT NOT NULL REFERENCES pms.grievances(grievance_id),
    id SERIAL, record_status TEXT DEFAULT 'Active', created_at TIMESTAMPTZ DEFAULT NOW(),
    soft_deleted BOOLEAN DEFAULT FALSE, status TEXT, updated_at TIMESTAMPTZ,
    created_by VARCHAR(100), updated_by VARCHAR(100), is_active BOOLEAN DEFAULT TRUE
);
