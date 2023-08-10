## API

### Login
- /login POST 学生/管理员登录
### User
- /users POST 登录
- /users GET 学生获取自己的信息, 管理员获取学生信息列表
- /users/{id} GET 管理员获取指定学生信息
- /users PUT 学生修改自己的信息
- /users DELETE 学生删除自己的账号
- /users/{id} DELETE 管理员删除指定学生账号
### Admin
- /admin PUT 管理员修改自己的信息

## JWT Token
- payload中添加role参数, 区分学生与管理员
- 修改密码后返回新的token, 旧token缓存到redis中, 7天过期(token有效期), JWT中间件先检测token是否是bad token

## 管理员账户
- 通过命令行添加, 数据库中存在管理员->跳过, 不存在->提示添加

## Header带Authorization的请求如何实现跨域访问
~~1. Origin不能是localhost~~
~~2. AllowHeaders包含Authorization (cors.Default()的AllowHeaders为空)~~
玄学!

## TODO
1. 数据库配置, Redis配置
2. 完善swagger的Response