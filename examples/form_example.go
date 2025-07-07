package examples

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/pkg/app"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" label:"用户名"`
	Email    string `json:"email" validate:"required,email" label:"邮箱"`
	Password string `json:"password" validate:"required,min=6" label:"密码"`
	Age      int    `json:"age" validate:"required,min=1,max=120" label:"年龄"`
	Phone    string `json:"phone" validate:"omitempty,len=11" label:"手机号"`
}

// UserSearchRequest 用户搜索请求
type UserSearchRequest struct {
	Keyword string `form:"keyword" validate:"omitempty,min=2" label:"关键词"`
	Page    int    `form:"page" validate:"omitempty,min=1" label:"页码"`
	Size    int    `form:"size" validate:"omitempty,min=1,max=100" label:"每页数量"`
	Status  string `form:"status" validate:"omitempty,oneof=active inactive" label:"状态"`
}

// FormProductRequest 产品请求
type FormProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" label:"产品名称"`
	Price       float64 `json:"price" validate:"required,gt=0" label:"价格"`
	Description string  `json:"description" validate:"omitempty,max=500" label:"描述"`
	Category    string  `json:"category" validate:"required,oneof=electronics clothing books" label:"分类"`
}

// RunFormExample 运行表单验证示例
func RunFormExample() {
	// 创建表单验证器实例
	formValidator := app.NewFormValidator()
	
	// 创建Gin引擎
	r := gin.Default()
	
	// 设置语言切换中间件
	r.Use(func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = c.Query("lang")
		}
		if lang != "" {
			formValidator.SetLanguage(lang)
		}
		c.Next()
	})
	
	// 示例1: 使用MustBindJSONAndValidate - 自动处理错误
	r.POST("/users", func(c *gin.Context) {
		var req UserCreateRequest
		
		// 使用Must方法，验证失败会自动返回错误响应
		if !formValidator.MustBindJSONAndValidate(c, &req) {
			return // 验证失败，已自动返回错误响应
		}
		
		// 验证成功，处理业务逻辑
		response.Success(c, gin.H{
			"message": "用户创建成功",
			"user":    req,
		})
	})
	
	// 示例2: 手动处理验证错误
	r.PUT("/users/:id", func(c *gin.Context) {
		var req UserCreateRequest
		
		// 手动绑定和验证
		errs := formValidator.BindJSONAndValidate(c, &req)
		if len(errs) > 0 {
			// 自定义错误处理
			formValidator.HandleValidationErrorsWithMessage(c, errs, "用户更新参数验证失败")
			return
		}
		
		userID := c.Param("id")
		response.Success(c, gin.H{
			"message": "用户更新成功",
			"user_id": userID,
			"user":    req,
		})
	})
	
	// 示例3: 查询参数验证
	r.GET("/users/search", func(c *gin.Context) {
		var req UserSearchRequest
		
		// 绑定并验证查询参数
		if !formValidator.MustBindQueryAndValidate(c, &req) {
			return
		}
		
		// 设置默认值
		if req.Page == 0 {
			req.Page = 1
		}
		if req.Size == 0 {
			req.Size = 10
		}
		
		response.Success(c, gin.H{
			"message": "搜索成功",
			"params":  req,
			"results": []string{"用户1", "用户2", "用户3"},
		})
	})
	
	// 示例4: 分离绑定和验证
	r.POST("/products", func(c *gin.Context) {
		var req FormProductRequest
		
		// 先绑定参数
		if err := formValidator.ShouldBindJSON(c, &req); err != nil {
			response.BadRequest(c, "JSON格式错误")
			return
		}
		
		// 业务逻辑处理（例如设置默认值）
		if req.Description == "" {
			req.Description = "暂无描述"
		}
		
		// 再验证参数
		errs := formValidator.Validate(&req)
		if len(errs) > 0 {
			formValidator.HandleValidationErrors(c, errs)
			return
		}
		
		response.Success(c, gin.H{
			"message": "产品创建成功",
			"product": req,
		})
	})
	
	// 示例5: 验证单个字段
	r.POST("/validate-email", func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		
		if err := formValidator.ShouldBindJSON(c, &req); err != nil {
			response.BadRequest(c, "JSON格式错误")
			return
		}
		
		// 验证单个字段
		errs := formValidator.ValidateVar(req.Email, "required,email")
		if len(errs) > 0 {
			formValidator.HandleValidationErrors(c, errs)
			return
		}
		
		response.Success(c, gin.H{
			"message": "邮箱格式正确",
			"email":   req.Email,
		})
	})
	
	// 示例6: 自定义验证规则
	r.POST("/register-custom-validation", func(c *gin.Context) {
		// 注册自定义验证规则
		formValidator.RegisterValidation("custom_username", func(fl app.FieldLevel) bool {
			username := fl.Field().String()
			// 自定义规则：用户名不能包含特殊字符
			for _, char := range username {
				if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
					return false
				}
			}
			return true
		})
		
		// 注册自定义翻译
		formValidator.RegisterTranslation("custom_username", "用户名只能包含字母、数字和下划线")
		
		response.Success(c, gin.H{
			"message": "自定义验证规则注册成功",
		})
	})
	
	// 示例7: 测试自定义验证
	r.POST("/test-custom-validation", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" validate:"required,custom_username"`
		}
		
		if !formValidator.MustBindJSONAndValidate(c, &req) {
			return
		}
		
		response.Success(c, gin.H{
			"message":  "用户名验证通过",
			"username": req.Username,
		})
	})
	
	// 示例8: 语言切换
	r.POST("/set-language/:lang", func(c *gin.Context) {
		lang := c.Param("lang")
		err := formValidator.SetLanguage(lang)
		if err != nil {
			response.BadRequest(c, "不支持的语言")
			return
		}
		
		response.Success(c, gin.H{
			"message":  "语言设置成功",
			"language": formValidator.GetLanguage(),
		})
	})
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":   "ok",
			"language": formValidator.GetLanguage(),
		})
	})
	
	// 显示所有可用的接口
	r.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message": "Form Validation Example API",
			"endpoints": []string{
				"POST /users - 创建用户（自动错误处理）",
				"PUT /users/:id - 更新用户（手动错误处理）",
				"GET /users/search - 搜索用户（查询参数验证）",
				"POST /products - 创建产品（分离绑定和验证）",
				"POST /validate-email - 验证邮箱（单字段验证）",
				"POST /register-custom-validation - 注册自定义验证规则",
				"POST /test-custom-validation - 测试自定义验证",
				"POST /set-language/:lang - 设置语言（zh/en）",
				"GET /health - 健康检查",
			},
			"tips": []string{
				"可以通过 Accept-Language 头部或 ?lang=zh/en 参数切换语言",
				"支持的验证标签: required, min, max, email, oneof 等",
				"可以注册自定义验证规则和翻译",
			},
		})
	})
	
	fmt.Println("Form Validation Example Server starting on :8082")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  / - API文档")
	fmt.Println("  POST /users - 创建用户")
	fmt.Println("  PUT  /users/:id - 更新用户")
	fmt.Println("  GET  /users/search - 搜索用户")
	fmt.Println("  POST /products - 创建产品")
	fmt.Println("  POST /validate-email - 验证邮箱")
	fmt.Println("  POST /register-custom-validation - 注册自定义验证")
	fmt.Println("  POST /test-custom-validation - 测试自定义验证")
	fmt.Println("  POST /set-language/:lang - 设置语言")
	fmt.Println("  GET  /health - 健康检查")
	
	r.Run(":8082")
}