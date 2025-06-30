package test

import (
	"testing"
	"github.com/cuiyuanxin/kunpeng/pkg/utils"
)

// TestStringUtils 测试字符串工具函数
func TestStringUtils(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		if !utils.IsEmpty("") {
			t.Error("Expected empty string to be empty")
		}
		if !utils.IsEmpty("   ") {
			t.Error("Expected whitespace string to be empty")
		}
		if utils.IsEmpty("hello") {
			t.Error("Expected non-empty string to not be empty")
		}
	})

	t.Run("Capitalize", func(t *testing.T) {
		result := utils.Capitalize("hello")
		expected := "Hello"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		result := utils.Reverse("hello")
		expected := "olleh"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("ContainsIgnoreCase", func(t *testing.T) {
		if !utils.ContainsIgnoreCase("Hello World", "WORLD") {
			t.Error("Expected case-insensitive match")
		}
		if utils.ContainsIgnoreCase("Hello", "xyz") {
			t.Error("Expected no match")
		}
	})
}