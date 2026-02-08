package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"gorm.io/gorm"
)

// ProjectRepository provides data access for Project and Committee entities.
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project/committee repository.
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("soft_deleted = ?", false)
}

// ─── Project ─────────────────────────────────────────────────────────────────

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id string) (*performance.Project, error) {
	var p performance.Project
	err := r.base(ctx).
		Preload("ProjectMembers").
		Preload("ProjectWorkProducts").
		Preload("ProjectObjectives").
		Preload("ProjectAssignedWorkProducts").
		First(&p, "project_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetProjectByID: %w", err)
	}
	return &p, nil
}

func (r *ProjectRepository) GetProjectsByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.Project, error) {
	var results []performance.Project
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("ProjectMembers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetProjectsByReviewPeriod: %w", err)
	}
	return results, nil
}

func (r *ProjectRepository) GetProjectsByDepartment(ctx context.Context, deptID int) ([]performance.Project, error) {
	var results []performance.Project
	err := r.base(ctx).
		Where("department_id = ?", deptID).
		Preload("ProjectMembers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetProjectsByDepartment: %w", err)
	}
	return results, nil
}

func (r *ProjectRepository) GetProjectsByMemberStaffID(ctx context.Context, staffID string) ([]performance.Project, error) {
	var results []performance.Project
	err := r.base(ctx).
		Joins("JOIN pms.project_members pm ON pm.project_id = projects.project_id").
		Where("pm.staff_id = ? AND pm.soft_deleted = ?", staffID, false).
		Preload("ProjectMembers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetProjectsByMemberStaffID: %w", err)
	}
	return results, nil
}

// ─── ProjectMember ───────────────────────────────────────────────────────────

func (r *ProjectRepository) GetProjectMembersByProject(ctx context.Context, projectID string) ([]performance.ProjectMember, error) {
	var results []performance.ProjectMember
	err := r.base(ctx).
		Where("project_id = ?", projectID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetProjectMembersByProject: %w", err)
	}
	return results, nil
}

// ─── Committee ───────────────────────────────────────────────────────────────

func (r *ProjectRepository) GetCommitteeByID(ctx context.Context, id string) (*performance.Committee, error) {
	var c performance.Committee
	err := r.base(ctx).
		Preload("CommitteeMembers").
		Preload("CommitteeWorkProducts").
		Preload("CommitteeObjectives").
		Preload("CommitteeAssignedWorkProducts").
		First(&c, "committee_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetCommitteeByID: %w", err)
	}
	return &c, nil
}

func (r *ProjectRepository) GetCommitteesByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.Committee, error) {
	var results []performance.Committee
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("CommitteeMembers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetCommitteesByReviewPeriod: %w", err)
	}
	return results, nil
}

func (r *ProjectRepository) GetCommitteesByMemberStaffID(ctx context.Context, staffID string) ([]performance.Committee, error) {
	var results []performance.Committee
	err := r.base(ctx).
		Joins("JOIN pms.committee_members cm ON cm.committee_id = committees.committee_id").
		Where("cm.staff_id = ? AND cm.soft_deleted = ?", staffID, false).
		Preload("CommitteeMembers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetCommitteesByMemberStaffID: %w", err)
	}
	return results, nil
}

// ─── CommitteeMember ─────────────────────────────────────────────────────────

func (r *ProjectRepository) GetCommitteeMembersByCommittee(ctx context.Context, committeeID string) ([]performance.CommitteeMember, error) {
	var results []performance.CommitteeMember
	err := r.base(ctx).
		Where("committee_id = ?", committeeID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("projectRepo.GetCommitteeMembersByCommittee: %w", err)
	}
	return results, nil
}

// ─── Transaction helper ──────────────────────────────────────────────────────

func (r *ProjectRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *ProjectRepository) DB() *gorm.DB {
	return r.db
}
