# 数据模型 — RBAC（用户 / 角色 / 权限 / 日志）

> 域文档，从 `运营平台-P1数据模型草案.md` 拆出。P1 第 1 块。
> 范围：sys_user / sys_role / sys_menu / sys_permission / sys_role_menu / sys_role_permission / sys_login_log / sys_operation_log。
> 目标：支持**多终端**（运营平台 / 商户后台 / 代理后台）、**角色基权限**（菜单 + 操作点）、**内置角色保护**、提现审核等操作权限可分配。P1 只实现**运营平台**登录与菜单；商户后台/代理后台结构先就位、功能后置。

最后更新：2026-07-15
状态：⬜ 待评审

---

## 1. 设计原则

1. **多终端**：同一套用户/角色表服务三个终端（运营平台=1 / 商户后台=2 / 代理后台=3）。`sys_user.terminal` 决定登录入口；`sys_menu.terminal` 决定菜单归属；`sys_role.role_type` 与终端域对应。
2. **角色基权限**：权限挂在角色上（角色→菜单 M:N、角色→操作权限 M:N），用户→角色。P1 用户单角色（`sys_user.role_id`）；结构允许后续升 M:N。
3. **内置角色受保护**：`is_builtin=1` 的角色（如 `administrator`）不可删、不可改 role_type/code；`is_admin=1` 角色默认拥有全部菜单。
4. **数据范围（P1 简化）**：运营平台用户看全量；商户后台用户只看 `merchant_id` 范围；代理后台用户只看 `agent_id` 范围（P1 运营平台无需数据范围，结构预留）。
5. **审计字段每表显式列出**：`created_at / updated_at / created_by / updated_by`；业务主体表带 `deleted_at`（软删）。日志表只追加。
6. **密码**：`login_password_hash` 用 bcrypt/argon2，**禁存明文**；用户默认密码来自系统配置「用户默认密码」。

---

## 2. 表清单

### 2.1 sys_user 用户

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| login_account | varchar(64) | 登录账号，**unique** |
| login_password_hash | varchar(255) | 密码哈希 |
| nickname | varchar(64) | 昵称 |
| terminal | tinyint | 1运营平台 / 2商户后台 / 3代理后台 |
| merchant_id | bigint | 所属商户（terminal=2 时填） |
| agent_id | bigint | 所属代理（terminal=3 时填） |
| role_id | bigint FK→sys_role | 角色（P1 单角色） |
| google_secret | varchar(255) | 谷歌验证密钥（加密；登录谷歌验证开启时用） |
| register_ip | varchar(64) | 注册IP |
| last_login_at | datetime | 最后登录 |
| status | tinyint | 1启用 / 0禁用 |
| remark | varchar(255) | |
| created_at / updated_at / deleted_at | datetime | |
| created_by / updated_by | bigint | 用户ID |

> 索引：`login_account`、`(terminal, status)`、`(merchant_id)`、`(agent_id)`。
> 谷歌验证：登录页是否展示谷歌验证码字段，由系统配置「登录谷歌验证」+ 用户是否绑定 google_secret 共同决定（PRD §1 关系待确认 → 此处落库支持两种）。
> **商户/代理用户**：商户后台用户 = `sys_user`(terminal=2, merchant_id, role_type=商户系统)；代理后台用户 = `sys_user`(terminal=3, agent_id, role_type=代理系统)。**创建商户/代理时联动建对应 sys_user**，登录账号唯一性在此校验；merchant/agent 表不再存登录账号。

### 2.2 sys_role 角色

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| code | varchar(64) | 角色编码，**unique** |
| name | varchar(64) | 角色名称 |
| role_type | tinyint | 1运营平台 / 2超级管理 / 3商户系统 / 4代理系统 |
| is_admin | tinyint | 是否管理员（1拥有全部菜单） |
| is_builtin | tinyint | 内置角色（1受保护，不可删/不可改 code、role_type） |
| sort | int | 排序号 |
| status | tinyint | 1正常 / 0停用 |
| remark | varchar(255) | |
| created_at / updated_at / deleted_at | datetime | |
| created_by / updated_by | bigint | |

> 约束：`is_builtin=1` 的行禁止物理删除（只能停用）。删除接口需校验。

### 2.3 sys_menu 菜单（树）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| parent_id | bigint | 父菜单（0=根） |
| name | varchar(64) | 菜单名称 |
| path | varchar(128) | 路由路径 |
| component | varchar(128) | 前端组件（目录节点为空） |
| icon | varchar(64) | 图标 |
| menu_type | tinyint | 1目录 / 2菜单 / 3按钮（操作点占位） |
| permission_key | varchar(128) | 权限标识（按钮型必填，如 `withdraw:audit`） |
| terminal | tinyint | 归属终端（同 sys_user.terminal） |
| sort | int | 排序 |
| visible | tinyint | 是否显示 |
| status | tinyint | 1启用 / 0停用 |
| created_at / updated_at | datetime | |
| created_by / updated_by | bigint | |

> 索引：`(parent_id, sort)`、`(terminal, status)`。
> 按钮型菜单（menu_type=3）用于页面内操作点（如「提现审核」「一键结算」），通过 permission_key 与角色关联。

### 2.4 sys_permission 操作权限点（可选，与按钮型菜单二选一）

> 若操作权限完全用 sys_menu(menu_type=3) 表达，则本表可省；若需独立维护「功能权限点」，用本表。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| code | varchar(128) | 权限编码，**unique**（如 `withdraw:audit`） |
| name | varchar(64) | 名称（提现审核） |
| module | varchar(32) | 所属模块（提现管理） |
| created_at / updated_at | datetime | |

> **决策点 A**：操作权限用「sys_menu.menu_type=3」还是「独立 sys_permission」？建议 P1 用前者（少一张表，菜单树即权限树）。

### 2.5 sys_role_menu 角色-菜单

| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | bigint FK→sys_role | |
| menu_id | bigint FK→sys_menu | |

> 主键：`(role_id, menu_id)`。`is_admin=1` 角色查询时返回全部，不入此表（或入全量，二选一）。

### 2.6 sys_role_permission 角色-权限点（仅当启用 sys_permission）

| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | bigint FK→sys_role | |
| permission_id | bigint FK→sys_permission | |

> 主键：`(role_id, permission_id)`。

### 2.7 sys_login_log 登录日志

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| login_account | varchar(64) | 登录账号（冗余，便于无用户也记录） |
| user_id | bigint | 用户ID（可空，账号不存在时） |
| status | tinyint | 1成功 / 0失败 |
| message | varchar(64) | 操作消息：成功/密码错误/验证码错误/账号停用/账号不存在等 |
| ip | varchar(64) | |
| location | varchar(64) | 地点 |
| os | varchar(64) | |
| access_time | datetime | 访问时间 |
| detail | json | 扩展（UA 等） |

> 索引：`(login_account, access_time)`、`(access_time)`。只追加。

### 2.8 sys_operation_log 操作日志

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | |
| module | varchar(32) | 模块（商户管理/通道管理/提现管理/…，含代理管理） |
| content | varchar(255) | 操作内容（编辑商户/一键结算/…） |
| method | varchar(128) | 方法名 |
| ip / location / os | varchar | |
| operator_id | bigint | 操作人ID |
| operator_name | varchar(64) | 操作人（冗余） |
| detail | json | 变更前后（before/after） |
| created_at | datetime | |

> 索引：`(module, created_at)`、`(operator_id, created_at)`。只追加。
> **所有资金类操作（调额/改费率/结算/提现审核）必须落此表**。

---

## 3. 关键流程

### 3.1 登录认证
1. 校验 `login_account` + `login_password_hash` + 图形验证码。
2. 若系统配置「登录谷歌验证」开 **且** 用户已绑 google_secret → 校验谷歌码。
3. 写 `sys_login_log`（成功/失败 + message）。
4. 成功 → 签发 token，载荷含 `user_id / terminal / role_id`。

### 3.2 鉴权校验（每次请求）
1. 按 `role_id` 取 `sys_role_menu` → 当前角色可见菜单 + 可用操作点（permission_key）。
2. `is_admin=1` → 放行全部。
3. 后端接口用 permission_key 校验（如提现审核需 `withdraw:audit`）；前端用菜单树渲染导航 + 按钮显隐。

### 3.3 权限分配（角色权限弹窗）
- 勾选菜单树 → 写 `sys_role_menu`；勾选操作点 → 写 `sys_role_permission`（或同 menu_type=3）。
- 内置角色禁止改权限（或允许但不可删角色）。
- 提现审核权限：在角色权限弹窗勾选 `withdraw:audit`，分配给「财务」角色。

---

## 4. 内置数据（种子）

- **sys_role**：`administrator`（role_type=超级管理, is_admin=1, is_builtin=1）；可选 `运营专员`、`财务`。
- **sys_menu**：按 PRD §0 侧栏初始化运营平台菜单树（首页/商户管理/代理管理/通道管理/通道类型/供应商管理/订单管理/异常订单/商户提现/代理提现/通道提现/数据统计(账户记录·日终统计)/系统管理(系统配置·登录日志·操作日志·用户·角色·机器人配置)）+ 各页按钮型操作点。

---

## 5. 待评审开放点

- **A** 操作权限用 menu_type=3 还是独立 sys_permission？（建议前者）
- **B** 用户-角色：P1 单角色（role_id）还是 M:N（sys_user_role）？（建议 P1 单角色）
- **C** `is_admin=1` 角色：不入 sys_role_menu（查询时返回全部）还是入全量？（建议不入，查询特判）
- **D** 谷歌验证：落库 google_secret 已支持；登录页是否默认显示谷歌字段 = 系统配置开关 + 是否绑定，二者关系确认。
- **E** 商户后台/代理后台用户的 RBAC 功能是否进 P1？（默认 P1 只做运营平台）
