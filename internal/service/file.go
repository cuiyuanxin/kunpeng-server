package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/config"
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// FileService 文件服务
type FileService struct {
	config *config.Config
	db     *gorm.DB
}

// NewFileService 创建文件服务
func NewFileService(cfg *config.Config, db *gorm.DB) *FileService {
	return &FileService{
		config: cfg,
		db:     db,
	}
}

// Upload 上传文件
func (s *FileService) Upload(file *multipart.FileHeader) (*model.FileInfo, error) {
	// 检查文件大小
	maxSize := int64(10 << 20) // 10MB
	if file.Size > maxSize {
		return nil, fmt.Errorf("file size exceeds limit: %d bytes", maxSize)
	}

	// 检查文件类型
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("file type not allowed: %s", ext)
	}

	// 生成文件名
	filename := s.generateFilename(file.Filename)

	// 创建上传目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 创建文件信息
	fileInfo := &model.FileInfo{
		OriginalName: file.Filename,
		Filename:     filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		FileExt:      ext,
		URL:          fmt.Sprintf("/uploads/%s", filename),
	}

	// 保存到数据库
	if s.db != nil {
		if err := s.db.Create(fileInfo).Error; err != nil {
			return nil, fmt.Errorf("failed to save file info: %w", err)
		}
	}

	return fileInfo, nil
}

// GetList 获取文件列表
func (s *FileService) GetList(page, limit int) ([]model.FileInfo, int64, error) {
	if s.db == nil {
		return nil, 0, fmt.Errorf("database not initialized")
	}

	var files []model.FileInfo
	var total int64

	// 获取总数
	if err := s.db.Model(&model.FileInfo{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get files: %w", err)
	}

	return files, total, nil
}

// GetByID 根据ID获取文件信息
func (s *FileService) GetByID(id uint) (*model.FileInfo, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var fileInfo model.FileInfo
	if err := s.db.First(&fileInfo, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &fileInfo, nil
}

// Delete 删除文件
func (s *FileService) Delete(id uint) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 获取文件信息
	fileInfo, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(fileInfo.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %w", err)
	}

	// 删除数据库记录
	if err := s.db.Delete(&model.FileInfo{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

// generateFilename 生成文件名
func (s *FileService) generateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	uuid := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s%s", timestamp, uuid, ext)
}

// GetFileURL 获取文件访问URL
func (s *FileService) GetFileURL(filename string) string {
	return fmt.Sprintf("/uploads/%s", filename)
}

// IsImageFile 判断是否为图片文件
func (s *FileService) IsImageFile(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	ext = strings.ToLower(ext)
	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}