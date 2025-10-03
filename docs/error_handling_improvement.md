# GORM 错误处理改进

## 概述

本次改进主要针对项目中使用 GORM 查询数据时 `record not found` 错误没有统一处理的问题。通过引入统一的错误处理机制，确保所有数据库查询错误都能被正确识别和处理。

## 改进内容

### 1. 基础仓储错误处理方法

在 `BaseRepository` 中新增了 `HandleDBError` 方法：

```go
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
```

### 2. 更新的 Repository 文件

以下 repository 实现文件已更新使用统一的错误处理：

- `user_repository_impl.go`
- `role_repository_impl.go`
- `menu_repository_impl.go`
- `api_repository_impl.go`
- `post_repository_impl.go`
- `dept_repository_impl.go`
- `dict_repository_impl.go`

### 3. 错误码映射

- `gorm.ErrRecordNotFound` → `kperrors.ErrDBNotFound` (10104)
- 其他数据库错误 → `kperrors.ErrDatabase` (10013)

## 使用示例

### 修改前
```go
func (r *UserRepositoryImpl) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, errors.New(errors.ErrDatabase, err) // 所有错误都被当作数据库错误
    }
    return &user, nil
}
```

### 修改后
```go
func (r *UserRepositoryImpl) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, r.HandleDBError(err) // 统一错误处理，区分记录不存在和其他数据库错误
    }
    return &user, nil
}
```

## 优势

1. **统一错误处理**：所有 repository 层的数据库错误都通过统一方法处理
2. **错误类型区分**：能够区分记录不存在错误和其他数据库错误
3. **便于维护**：错误处理逻辑集中在一个地方，便于后续维护和扩展
4. **更好的用户体验**：前端可以根据不同的错误码提供更准确的错误提示

## 注意事项

- Service 层已经有部分地方正确处理了 `gorm.ErrRecordNotFound` 错误，这些地方保持不变
- 新的错误处理机制向下兼容，不会影响现有的业务逻辑
- 建议在新增 repository 方法时都使用 `HandleDBError` 方法处理错误

## 测试

可以通过以下方式测试错误处理是否正常工作：

```go
// 测试记录不存在错误
user, err := userRepo.FindByID(99999) // 不存在的ID
if kpErr, ok := err.(*errors.Error); ok {
    if kpErr.Code == errors.ErrDBNotFound {
        // 正确处理了记录不存在错误
    }
}
```