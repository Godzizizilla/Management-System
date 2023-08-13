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
//		 Bearer:
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

// swagger:route POST /users Public register
// 学生注册接口
// responses:
//	  default: response
//    200: tokenResponse

// swagger:route PUT /users Users updateUser
// 用户修改个人信息接口
// responses:
//	  default: response
// Security:
// - Bearer: []

// swagger:route DELETE /users/{id} Users deleteUser
// 删除学生账户接口
// parameters:
//   + name: id
//     in: path
//     description: me/学生ID
//     type: string
//     required: true
// responses:
//	  default: response
// Security:
// - Bearer: []

// swagger:route GET /users/{id} Users getUser
// 获取学生/管理员信息接口
// Parameters:
//   + name: id
//     in: path
//     description: me: 获取当前用户信息, all: 获取所有学生信息, {student_id}: 获取指定学生信息
//     type: string
//     required: true
// responses:
//	  default: response
//    200: getUserResponse
//    200: getAdminResponse
//    200: getUserListResponse
// Security:
// - Bearer: []

// swagger:parameters login
type loginResponseWrapper struct {
	// 登录需要提供 用户名(学生ID / 管理员账户名) + 密码
	// in:body
	Body models.LoginRequest
}

// swagger:parameters register
type registerResponseWrapper struct {
	// 注册需要提供 姓名 + 密码 + 年级 + 学生ID
	// in:body
	Body models.RegisterRequest
}

// swagger:parameters updateUser
type updateInfoRequestWrapper struct {
	// 修改姓名不需要提供密码, 修改密码需要提供旧密码
	// in:body
	Body models.UpdateInfoRequest
}

// swagger:parameters deleteUser
type deleteRequestWrapper struct {
	// 学生删除自己账户需要提供密码, 管理员删除学生账户不需要提供密码
	// in:body
	Body models.DeleteRequest
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

// 获取管理员信息返回的Response
// swagger:response getAdminResponse
type getAdminResponseWrapper struct {
	// in:body
	Body models.GetAdminResponse
}

// 获取学生列表返回的Response
// swagger:response getUserListResponse
type getUserListResponseWrapper struct {
	// in:body
	Body models.GetUserListResponse
}
