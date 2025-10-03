package task

import (
	"log"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/service"
	serviceImpl "github.com/cuiyuanxin/kunpeng/internal/service"
)

// TokenCleanupTask token黑名单清理任务
type TokenCleanupTask struct {
	tokenBlacklistService service.TokenBlacklistService
}

// NewTokenCleanupTask 创建token清理任务
func NewTokenCleanupTask() *TokenCleanupTask {
	return &TokenCleanupTask{
		tokenBlacklistService: serviceImpl.GetTokenBlacklistService(),
	}
}

// Run 执行清理任务
func (t *TokenCleanupTask) Run() {
	log.Println("开始清理过期的token黑名单记录...")

	err := t.tokenBlacklistService.CleanExpiredTokens()
	if err != nil {
		log.Printf("清理过期token黑名单记录失败: %v", err)
		return
	}

	log.Println("清理过期的token黑名单记录完成")
}

// StartCleanupScheduler 启动定时清理任务
// 每天凌晨2点执行一次清理任务
func (t *TokenCleanupTask) StartCleanupScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		// 计算到下一个凌晨2点的时间
		now := time.Now()
		next2AM := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
		if now.Hour() < 2 {
			next2AM = time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
		}

		// 等待到凌晨2点
		time.Sleep(time.Until(next2AM))

		// 立即执行一次
		t.Run()

		// 然后每24小时执行一次
		for range ticker.C {
			t.Run()
		}
	}()

	log.Println("token黑名单清理定时任务已启动，将在每天凌晨2点执行")
}
