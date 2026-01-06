package jwts

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/models/enum"
	"strings"
	"time"
)

// Claims 自定义载荷（Payload）
// 作用：决定了 Token 里面存哪些用户信息。
// 这里的 UserID, Username, Role 是为了让后端不用查数据库就能知道这是谁，有什么权限。
type Claims struct {
	UserID   uint          `json:"userID"`
	Username string        `json:"username"`
	Role     enum.RoleType `json:"role"`
}

// MyClaims 最终的 JWT 结构
// 作用：把我们自定义的 Claims 和 JWT 标准字段（如过期时间、签发人）组合在一起。
type MyClaims struct {
	Claims
	jwt.StandardClaims
}

// GetUser 从 Token 信息反查数据库获取完整用户对象
// 作用：如果你觉得 Token 里的信息不够（比如你需要用户的头像、邮箱），可以用这个方法回数据库查。
// 接收者：(m MyClaims) 表示必须要在拿到了 Claims 之后才能调这个方法。
func (m MyClaims) GetUser() (user models.UserModel, err error) {
	err = global.DB.Take(&user, m.UserID).Error
	return
}

// GetToken 生成 Token (造票)
// 场景：通常在【用户登录】接口中使用。
// 输入：用户的基本信息 (ID, 用户名, 角色)。
// 输出：一个长长的加密字符串 (Token)。
func GetToken(claims Claims) (string, error) {
	cla := MyClaims{
		Claims: claims,
		StandardClaims: jwt.StandardClaims{
			// 设置过期时间：当前时间 + 配置文件里写的过期小时数
			ExpiresAt: time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour).Unix(),
			// 设置签发人
			Issuer: global.Config.Jwt.Issuer,
		},
	}
	// 使用 HS256 算法进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	// 使用配置文件里的 Secret (密钥) 生成最终字符串
	return token.SignedString([]byte(global.Config.Jwt.Secret))
}

// ParseToken 解析并校验 Token (检票 - 核心逻辑)
// 场景：通常在【中间件】中使用，验证用户传来的 Token 是否合法。
// 输入：Token 字符串。
// 输出：解析出来的用户信息结构体 (MyClaims) 或 错误信息。
func ParseToken(tokenString string) (*MyClaims, error) {
	if tokenString == "" {
		return nil, errors.New("请登录")
	}

	// 核心解析函数，会验证签名是否正确
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil // 必须使用相同的密钥才能解密
	})

	// 错误处理：将晦涩的英文错误翻译成中文提示
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token过期")
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return nil, errors.New("token无效")
		}
		if strings.Contains(err.Error(), "token contains an invalid") {
			return nil, errors.New("token非法")
		}
		return nil, err
	}

	// 如果解析成功且 Token 有效，返回 claims 对象
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseTokenByGin 从 HTTP 请求中提取并解析 Token (检票 - 便捷版)
// 作用：自动去 Header 或 Query 参数里找 Token，然后调用 ParseToken。
// 场景：在某些不需要经过中间件，但又要手动校验 Token 的接口里用。
func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	// 1. 优先从请求头 Header 中获取
	token := c.GetHeader("token")
	// 2. 如果 Header 里没有，尝试从 URL 参数中获取 (?token=xxx)
	if token == "" {
		token = c.Query("token")
	}

	return ParseToken(token)
}

// GetClaims 从 Gin 上下文中获取用户信息 (读票)
// 场景：在【Controller/API】层使用。
// 前提：必须已经经过了 JWT 中间件（中间件会把解析好的 claims 塞入 context）。
// 作用：不用重新解析 Token，直接拿中间件存好的数据，比如 jwts.GetClaims(c).UserID。
func GetClaims(c *gin.Context) (claims *MyClaims) {
	_claims, ok := c.Get("claims") // "claims" 这个 key 必须和中间件里 Set 的 key 一致
	if !ok {
		return
	}
	// 类型断言：把 interface{} 转回 MyClaims 指针
	claims, ok = _claims.(*MyClaims)
	if !ok {
		return
	}
	return
}
