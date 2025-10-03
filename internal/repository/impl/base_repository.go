package impl

import (
	"errors"

	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// BaseRepository 基础仓储
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository() BaseRepository {
	return BaseRepository{
		db: database.GetDB(),
	}
}

// NewBaseRepositoryPtr 创建基础仓储指针
func NewBaseRepositoryPtr() *BaseRepository {
	return &BaseRepository{
		db: database.GetDB(),
	}
}

// HandleDBError 处理数据库错误
func (r *BaseRepository) HandleDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return kperrors.New(kperrors.ErrDBNotFound, err)
	}
	return kperrors.New(kperrors.ErrDatabase, err)
}
