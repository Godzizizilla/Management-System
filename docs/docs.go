// Package docs Management System API.
//
// 学生信息管理系统.
//
// Version: 1.0.0
// Schemes: http
// Host: localhost:7890
// BasePath: /v1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
//			SecurityDefinitions:
//		 bearer:
//		      type: apiKey
//		      name: Authorization
//	       in: header
//
// swagger:meta
package docs

import "github.com/Godzizizilla/Management-System/models"

// swagger:route POST /login Public login
// 学生, 管理员登录接口
// responses:
//	  default: response
//    200: tokenResponse

// swagger:parameters login
type loginResponseWrapper struct {
	// 登录需要提供 用户名(学生ID / 管理员账户名) + 密码 + 角色(student / admin)
	// in:body
	Body models.LoginRequest
}

// swagger:route POST /users Public register
// 学生注册接口
// responses:
//	  default: response
//    200: tokenResponse

// swagger:parameters register
type registerResponseWrapper struct {
	// 注册需要提供 姓名 + 密码 + 年级 + 学生ID
	// in:body
	Body models.RegisterRequest
}

// swagger:route PUT /users Users updateUser
// 学生修改个人信息接口
// responses:
//	  default: response
// Security:
// - bearer: []

// swagger:parameters updateUser
type updateInfoRequestWrapper struct {
	// 修改姓名不需要提供密码, 修改密码需要提供旧密码
	// in:body
	Body models.UpdateInfoRequest
}

// swagger:route DELETE /users Users deleteUser
// 删除学生账户接口
// responses:
//	  default: response
// Security:
// - bearer: []

// swagger:route DELETE /users/{id} Users deleteUser
// 删除学生账户接口
// Parameters:
//   + name: id
//     in: path
//     description: 学生ID
//     type: integer
//     required: true
// responses:
//	  default: response
// Security:
// - bearer: []

// swagger:parameters deleteUser
type deleteRequestWrapper struct {
	// 学生删除自己账户需要提供密码, 管理员删除学生账户不需要提供密码
	// in:body
	Body models.DeleteRequest
}

// swagger:route GET /users Users getUser
// 获取学生信息接口
// responses:
//	  default: response
//    200: getUserResponse
//    200: getUserListResponse
// Security:
// - bearer: []

// swagger:route GET /users/{id} Users getUser
// 获取学生信息接口
// Parameters:
//   + name: id
//     in: path
//     description: 学生ID
//     type: integer
//     required: true
// responses:
//	  default: response
//    200: getUserResponse
//    200: getUserListResponse
// Security:
// - bearer: []

// swagger:route PUT /admin Admin updateAdmin
// 管理员修改个人信息接口
// responses:
//	  default: response
// Security:
// - bearer: []

// swagger:parameters updateAdmin
type updateAdminRequestWrapper struct {
	// 修改密码需要提供旧密码
	// in:body
	Body models.UpdateInfoRequest
}

// 注册, 登录, 修改密码将返回包含token的Response
// swagger:response tokenResponse
type tokenResponseWrapper struct {
	// in:body
	Body models.TokenResponse
}

// 默认返回的的Response
// swagger:response response
type responseWrapper struct {
	// in:body
	Body models.Response
}

// 获取学生信息返回的Response
// swagger:response getUserResponse
type getUserResponseWrapper struct {
	// in:body
	Body models.GetUserResponse
}

// 获取学生列表返回的Response
// swagger:response getUserListResponse
type getUserListResponseWrapper struct {
	// in:body
	Body models.GetUserListResponse
}
