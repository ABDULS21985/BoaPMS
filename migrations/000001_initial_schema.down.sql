-- Reverse the initial schema migration

DROP TABLE IF EXISTS pms.pms_configurations;
DROP TABLE IF EXISTS pms.settings;
DROP TABLE IF EXISTS pms.sequence_numbers;
DROP TABLE IF EXISTS pmsaudit.auditable_attributes;
DROP TABLE IF EXISTS pmsaudit.auditable_entities;
DROP TABLE IF EXISTS pmsaudit.audit_logs;
DROP TABLE IF EXISTS pms.period_scores;
DROP TABLE IF EXISTS pms.work_products;
DROP TABLE IF EXISTS pms.office_objectives;
DROP TABLE IF EXISTS pms.division_objectives;
DROP TABLE IF EXISTS pms.department_objectives;
DROP TABLE IF EXISTS pms.enterprise_objectives;
DROP TABLE IF EXISTS pms.objective_categories;
DROP TABLE IF EXISTS pms.performance_review_periods;
DROP TABLE IF EXISTS pms.strategic_themes;
DROP TABLE IF EXISTS pms.strategies;
DROP TABLE IF EXISTS "CoreSchema".role_permissions;
DROP TABLE IF EXISTS "CoreSchema".permissions;
DROP TABLE IF EXISTS "CoreSchema".bank_years;
DROP TABLE IF EXISTS "CoreSchema".asp_net_users;
DROP TABLE IF EXISTS "CoreSchema".asp_net_roles;
DROP TABLE IF EXISTS "CoreSchema".offices;
DROP TABLE IF EXISTS "CoreSchema".divisions;
DROP TABLE IF EXISTS "CoreSchema".departments;
DROP TABLE IF EXISTS "CoreSchema".directorates;

DROP SCHEMA IF EXISTS pmsaudit;
DROP SCHEMA IF EXISTS pms;
DROP SCHEMA IF EXISTS "CoreSchema";
