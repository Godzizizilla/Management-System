## Setup

### 1. Clone项目并添加依赖
```shell
# Clone项目
git clone https://github.com/Godzizizilla/Management-System.git
# 添加依赖
cd Management-System
go mid tidy
```

### 2. 配置`config.yml`
配置文件在./config/config.yml

### 3. 添加环境变量`JWT_SECRET_KEY`
```shell
# 方式1-手动生成并添加
export JWT_SECRET_KEY=your_secret_key

# 方式2-自动生成
go run cmd/generate_secret_key/main.go -dir .
```
### 4. 启动服务
```shell
go run cmd/server/main.go 
```

### 5. 设置管理员账户
根据终端提示设置管理员账户

## API
### Public
- /login POST 学生/管理员登录
- /users POST 登录
### Protected
- /users/{me, all, student_id} GET 学生: 获取自己的信息; 管理员: 获取自己的信息, 获取指定学生信息, 获取学生信息列表
- /users PUT 学生, 管理员修改自己的信息
- /users/{me, student_id} DELETE 学生删除自己, 管理员删除指定学生账号

## 安全性
使用Redis缓存{key: jti, val: role}, 同时设定过期时间

jti而不是studentID或adminID, 二者可能重复, 也不应该是token本身(太长)

JWT中间件根据payload中的jti和role, 确保权限正确

修改密码, 删除旧{key: jti, val: role}, 生成新{key: jti, val: role}

需要考虑:
1. 确保jti的唯一性
2. 服务器负载会不会过大, 权衡安全性和性能的考虑
3. 一个用户获取多个token, 允许吗?