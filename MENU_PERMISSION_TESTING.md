# 菜单权限功能快速测试指南

## 测试前准备

### 1. 启动服务器
```bash
cd power-admin-server
go build -o power-admin.exe ./power.go
./power-admin.exe
```

服务器应在 `0.0.0.0:8888` 启动

### 2. 前端应用（可选）
```bash
cd power-admin-web
npm install
npm run dev
```

## 测试场景

### 场景1：用户有角色且角色绑定了菜单

**步骤：**

1. 登录用户（假设用户已绑定"Admin"角色）
   ```bash
   curl -X POST http://localhost:8888/api/admin/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin",
       "password": "admin"
     }'
   ```

2. 获取返回的token，然后请求菜单列表
   ```bash
   curl -X GET http://localhost:8888/api/admin/system/menus \
     -H "Authorization: Bearer <YOUR_TOKEN>"
   ```

**预期结果：**
- 返回200 OK
- 返回data数组包含该用户角色绑定的菜单
- 菜单按sort字段升序排列

**响应示例：**
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "parentId": 0,
      "menuName": "系统管理",
      "menuPath": "/system",
      "component": "Layout",
      "icon": "setting",
      "sort": 1,
      "status": 1,
      "menuType": 1,
      "createdAt": "2025-11-30 11:00:00",
      "children": [
        {
          "id": 2,
          "parentId": 1,
          "menuName": "用户管理",
          "menuPath": "/system/users",
          "component": "system/User",
          "icon": "user",
          "sort": 1,
          "status": 1,
          "menuType": 1,
          "createdAt": "2025-11-30 11:00:00",
          "children": []
        }
      ]
    }
  ],
  "total": 1
}
```

### 场景2：用户没有分配任何角色

**步骤：**

1. 创建一个没有分配角色的新用户
2. 用该用户登录并获取token
3. 请求菜单列表

**预期结果：**
- 返回200 OK
- data为空数组 []
- total为 0

**响应示例：**
```json
{
  "code": 200,
  "data": [],
  "total": 0
}
```

### 场景3：用户有多个角色，角色绑定不同菜单

**步骤：**

1. 创建用户并分配多个角色（如：Admin + Editor）
2. 分别为这两个角色绑定不同的菜单
3. 用该用户登录并请求菜单列表

**预期结果：**
- 返回所有角色绑定菜单的并集
- 避免重复（使用了DISTINCT）
- 菜单去重后按sort排序

### 场景4：菜单被禁用（status=0）

**步骤：**

1. 创建角色和用户，绑定菜单关系
2. 禁用该菜单（UPDATE menus SET status=0 WHERE id=xxx）
3. 请求菜单列表

**预期结果：**
- 即使角色绑定了该菜单，也不会显示
- 返回的菜单列表中不包含被禁用的菜单

### 场景5：没有提供Authorization头

**步骤：**
```bash
curl -X GET http://localhost:8888/api/admin/system/menus
```

**预期结果：**
- 返回401 Unauthorized
- 错误信息："missing authorization header"

### 场景6：提供无效的token

**步骤：**
```bash
curl -X GET http://localhost:8888/api/admin/system/menus \
  -H "Authorization: Bearer invalid_token_here"
```

**预期结果：**
- 返回401 Unauthorized
- 错误信息："invalid or expired token"

## 数据库查询示例

### 查看用户的所有角色
```sql
SELECT r.* FROM roles r
LEFT JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = 1 AND r.status = 1;
```

### 查看角色绑定的菜单
```sql
SELECT DISTINCT m.* FROM menus m
LEFT JOIN role_menus rm ON rm.menu_id = m.id
WHERE rm.role_id IN (1, 2) AND m.status = 1
ORDER BY m.sort;
```

### 查看某个用户能访问的所有菜单
```sql
SELECT DISTINCT m.* FROM menus m
LEFT JOIN role_menus rm ON rm.menu_id = m.id
LEFT JOIN user_roles ur ON ur.role_id = rm.role_id
WHERE ur.user_id = 1 AND m.status = 1
ORDER BY m.sort;
```

## 常见问题排查

### Q1: 返回空菜单列表
**原因可能：**
- 用户没有被分配角色
- 用户的角色没有绑定菜单
- 菜单被禁用（status=0）
- 数据库中role_menus表为空

**检查方法：**
```sql
-- 检查用户的角色
SELECT * FROM user_roles WHERE user_id = <USER_ID>;

-- 检查角色绑定的菜单
SELECT * FROM role_menus WHERE role_id IN (<ROLE_IDS>);

-- 检查菜单状态
SELECT * FROM menus WHERE status = 1;
```

### Q2: 返回了不应该显示的菜单
**原因可能：**
- 其他角色也绑定了该菜单（DISTINCT可能有问题）
- 菜单状态过滤有问题

**检查方法：**
- 查看SQL日志
- 检查DISTINCT是否生效
- 验证WHERE条件中的status = 1

### Q3: 树形结构不正确
**原因可能：**
- 菜单的parentId设置不正确
- Logic层构建树形结构的算法有问题

**检查方法：**
```sql
-- 检查菜单的parentId
SELECT id, parent_id, menu_name FROM menus ORDER BY parent_id, sort;
```

### Q4: 401错误 - 用户未认证
**解决方案：**
- 确保Authorization header格式正确：`Bearer <TOKEN>`
- 确保token未过期
- 检查token是否有效

### Q5: 查询性能问题
**优化建议：**
- 在role_menus表的(role_id, menu_id)上建立索引
- 在user_roles表的(user_id, role_id)上建立索引
- 考虑使用Redis缓存用户的菜单列表

## 验证清单

- [ ] 服务器成功启动
- [ ] 用户认证功能正常（能获取token）
- [ ] 有角色有菜单的用户能获取菜单列表
- [ ] 无角色的用户返回空菜单
- [ ] 禁用的菜单不显示
- [ ] 多角色用户返回所有角色菜单的并集
- [ ] 树形结构正确（父子关系）
- [ ] 菜单按sort排序
- [ ] 无Authorization header返回401
- [ ] 无效token返回401

## 性能测试

### 模拟大数据场景
```sql
-- 创建大量菜单
INSERT INTO menus (parent_id, menu_name, menu_path, component, icon, sort, status) 
VALUES (0, CONCAT('menu_', ROW_NUMBER()), CONCAT('/menu_', ROW_NUMBER()), 'Layout', 'icon', 1, 1);

-- 创建大量角色
INSERT INTO roles (name, description) 
VALUES (CONCAT('role_', ROW_NUMBER()), 'Test Role');

-- 创建大量关联
INSERT INTO role_menus (role_id, menu_id) 
VALUES ...; -- 大量插入

-- 查询性能
EXPLAIN SELECT DISTINCT menus.* FROM menus
LEFT JOIN role_menus ON role_menus.menu_id = menus.id
WHERE role_menus.role_id IN (1, 2, 3, 4, 5) AND menus.status = 1
ORDER BY menus.sort;
```

## 日志查看

### 后端日志
查看`power-admin-server`的标准输出，寻找：
- 数据库查询错误
- 权限验证失败日志
- 业务逻辑错误

### 数据库日志
如果启用了MySQL查询日志，可以查看实际执行的SQL：
```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';
SELECT * FROM mysql.general_log ORDER BY event_time DESC LIMIT 100;
```
