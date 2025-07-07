package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// StatService 统计服务
type StatService struct {
	db *gorm.DB
}

// NewStatService 创建统计服务
func NewStatService(db *gorm.DB) *StatService {
	return &StatService{
		db: db,
	}
}

// DashboardStats 仪表板统计数据
type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	TotalRoles       int64 `json:"total_roles"`
	TotalDepartments int64 `json:"total_departments"`
	TotalFiles       int64 `json:"total_files"`
	TodayLogins      int64 `json:"today_logins"`
}

// UserStats 用户统计
type UserStats struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetDashboardStats 获取仪表板统计数据
func (s *StatService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 总用户数
	if err := s.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	// 活跃用户数（最近30天登录过的用户）
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := s.db.Model(&model.User{}).Where("last_login_at > ?", thirtyDaysAgo).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// 总角色数
	if err := s.db.Model(&model.Role{}).Count(&stats.TotalRoles).Error; err != nil {
		return nil, fmt.Errorf("failed to count total roles: %w", err)
	}

	// 总部门数
	if err := s.db.Model(&model.Department{}).Count(&stats.TotalDepartments).Error; err != nil {
		return nil, fmt.Errorf("failed to count total departments: %w", err)
	}

	// 总文件数
	if err := s.db.Model(&model.FileInfo{}).Count(&stats.TotalFiles).Error; err != nil {
		return nil, fmt.Errorf("failed to count total files: %w", err)
	}

	// 今日登录数
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	todayEnd := todayStart.Add(24 * time.Hour)
	if err := s.db.Model(&model.User{}).Where("last_login_at >= ? AND last_login_at < ?", todayStart, todayEnd).Count(&stats.TodayLogins).Error; err != nil {
		return nil, fmt.Errorf("failed to count today logins: %w", err)
	}

	return stats, nil
}

// GetUserRegistrationStats 获取用户注册统计（最近7天）
func (s *StatService) GetUserRegistrationStats() ([]UserStats, error) {
	var stats []UserStats

	// 获取最近7天的数据
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dateStart, _ := time.Parse("2006-01-02", dateStr)
		dateEnd := dateStart.Add(24 * time.Hour)

		var count int64
		if err := s.db.Model(&model.User{}).Where("created_at >= ? AND created_at < ?", dateStart, dateEnd).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count user registrations for %s: %w", dateStr, err)
		}

		stats = append(stats, UserStats{
			Date:  dateStr,
			Count: count,
		})
	}

	return stats, nil
}

// GetUserLoginStats 获取用户登录统计（最近7天）
func (s *StatService) GetUserLoginStats() ([]UserStats, error) {
	var stats []UserStats

	// 获取最近7天的数据
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dateStart, _ := time.Parse("2006-01-02", dateStr)
		dateEnd := dateStart.Add(24 * time.Hour)

		var count int64
		if err := s.db.Model(&model.User{}).Where("last_login_at >= ? AND last_login_at < ?", dateStart, dateEnd).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count user logins for %s: %w", dateStr, err)
		}

		stats = append(stats, UserStats{
			Date:  dateStr,
			Count: count,
		})
	}

	return stats, nil
}

// GetRoleDistribution 获取角色分布统计
func (s *StatService) GetRoleDistribution() ([]map[string]interface{}, error) {
	type RoleCount struct {
		RoleName string `json:"role_name"`
		Count    int64  `json:"count"`
	}

	var roleCounts []RoleCount
	if err := s.db.Table("users").
		Select("roles.name as role_name, COUNT(users.id) as count").
		Joins("LEFT JOIN user_roles ON users.id = user_roles.user_id").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Group("roles.name").
		Scan(&roleCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get role distribution: %w", err)
	}

	var result []map[string]interface{}
	for _, rc := range roleCounts {
		result = append(result, map[string]interface{}{
			"name":  rc.RoleName,
			"value": rc.Count,
		})
	}

	return result, nil
}

// GetDepartmentUserCount 获取各部门用户数量统计
func (s *StatService) GetDepartmentUserCount() ([]map[string]interface{}, error) {
	type DeptCount struct {
		DeptName string `json:"dept_name"`
		Count    int64  `json:"count"`
	}

	var deptCounts []DeptCount
	if err := s.db.Table("users").
		Select("departments.name as dept_name, COUNT(users.id) as count").
		Joins("LEFT JOIN departments ON users.department_id = departments.id").
		Group("departments.name").
		Scan(&deptCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get department user count: %w", err)
	}

	var result []map[string]interface{}
	for _, dc := range deptCounts {
		result = append(result, map[string]interface{}{
			"name":  dc.DeptName,
			"value": dc.Count,
		})
	}

	return result, nil
}

// GetFileTypeStats 获取文件类型统计
func (s *StatService) GetFileTypeStats() ([]map[string]interface{}, error) {
	var stats []map[string]interface{}

	// 这里可以根据实际需求实现文件类型统计
	// 示例数据
	stats = append(stats, map[string]interface{}{
		"type":  "PDF",
		"count": 150,
	})
	stats = append(stats, map[string]interface{}{
		"type":  "Word",
		"count": 120,
	})
	stats = append(stats, map[string]interface{}{
		"type":  "Excel",
		"count": 80,
	})
	stats = append(stats, map[string]interface{}{
		"type":  "Image",
		"count": 200,
	})

	return stats, nil
}

// GetUserStats 获取用户统计信息
func (s *StatService) GetUserStats() (*DashboardStats, error) {
	// 复用现有的GetDashboardStats方法
	return s.GetDashboardStats()
}