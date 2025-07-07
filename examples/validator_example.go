package examples

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// UserRequest 用户请求结构体
type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" label:"用户名"`
	Email    string `json:"email" validate:"required,email" label:"邮箱"`
	Password string `json:"password" validate:"required,min=6" label:"密码"`
	Age      int    `json:"age" validate:"required,min=1,max=120" label:"年龄"`
	Phone    string `json:"phone" validate:"omitempty,len=11" label:"手机号"`
}

// ProductRequest 产品请求结构体
type ProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" label:"产品名称"`
	Price       float64 `json:"price" validate:"required,gt=0" label:"价格"`
	Description string  `json:"description" validate:"omitempty,max=500" label:"描述"`
	Category    string  `json:"category" validate:"required,oneof=electronics clothing books" label:"分类"`
}

// RunValidatorExample 运行验证器示例
func RunValidatorExample() {
	// 创建验证器实例
	v := validator.New()
	
	// 设置语言为中文
	v.SetLanguage("zh")
	
	// 创建Gin引擎
	r := gin.Default()
	
	// 用户注册接口示例
	r.POST("/register", func(c *gin.Context) {
		var req UserRequest
		
		// 使用自定义验证器进行绑定和验证
		valid, errs := v.BindJSONAndValid(c, &req)
		if !valid {
			// 验证失败，返回错误信息
			response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
			return
		}
		
		// 验证成功，处理业务逻辑
		response.Success(c, gin.H{
			"message": "用户注册成功",
			"user":    req,
		})
	})
	
	// 产品创建接口示例
	r.POST("/products", func(c *gin.Context) {
		var req ProductRequest
		
		// 使用自定义验证器进行绑定和验证
		valid, errs := v.BindJSONAndValid(c, &req)
		if !valid {
			// 验证失败，返回错误信息
			response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
			return
		}
		
		// 验证成功，处理业务逻辑
		response.Success(c, gin.H{
			"message": "产品创建成功",
			"product": req,
		})
	})
	
	// 表单绑定示例（支持form-data和query参数）
	r.GET("/search", func(c *gin.Context) {
		type SearchRequest struct {
			Keyword  string `form:"keyword" validate:"required,min=1" label:"关键词"`
			Page     int    `form:"page" validate:"omitempty,min=1" label:"页码"`
			PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100" label:"每页数量"`
		}
		
		var req SearchRequest
		
		// 使用BindAndValid支持多种绑定方式
		valid, errs := v.BindAndValid(c, &req)
		if !valid {
			response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
			return
		}
		
		response.Success(c, gin.H{
			"message": "搜索成功",
			"params":  req,
		})
	})
	
	// 直接验证结构体示例
	r.POST("/validate-only", func(c *gin.Context) {
		// 手动创建结构体实例
		user := UserRequest{
			Username: "test",
			Email:    "invalid-email", // 故意设置无效邮箱
			Password: "123",           // 故意设置过短密码
			Age:      25,
		}
		
		// 直接验证结构体
		valid, errs := v.Validate(&user)
		if !valid {
			response.Error(c, http.StatusBadRequest, "数据验证失败", errs)
			return
		}
		
		response.Success(c, gin.H{
			"message": "数据验证通过",
			"data":    user,
		})
	})
	
	// 自定义验证规则示例
	r.POST("/custom-validation", func(c *gin.Context) {
		// 注册自定义验证规则：检查用户名是否为admin
		v.RegisterValidation("not_admin", func(fl validator.FieldLevel) bool {
			return fl.Field().String() != "admin"
		})
		
		// 注册自定义翻译
		v.RegisterTranslation("not_admin", "{0}不能为admin")
		
		type CustomRequest struct {
			Username string `json:"username" validate:"required,not_admin" label:"用户名"`
		}
		
		var req CustomRequest
		valid, errs := v.BindJSONAndValid(c, &req)
		if !valid {
			response.Error(c, http.StatusBadRequest, "参数验证失败", errs)
			return
		}
		
		response.Success(c, gin.H{
			"message": "自定义验证通过",
			"data":    req,
		})
	})
	
	// 语言切换示例
	r.POST("/set-language/:lang", func(c *gin.Context) {
		lang := c.Param("lang")
		err := v.SetLanguage(lang)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "语言设置失败", err.Error())
			return
		}
		
		response.Success(c, gin.H{
			"message": fmt.Sprintf("语言已切换为: %s", lang),
		})
	})
	
	fmt.Println("验证器示例服务启动在 :8081")
	fmt.Println("测试接口:")
	fmt.Println("POST /register - 用户注册")
	fmt.Println("POST /products - 产品创建")
	fmt.Println("GET /search?keyword=test&page=1 - 搜索")
	fmt.Println("POST /validate-only - 直接验证")
	fmt.Println("POST /custom-validation - 自定义验证")
	fmt.Println("POST /set-language/en - 切换语言")
	
	r.Run(":8081")
}