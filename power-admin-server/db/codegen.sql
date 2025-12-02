-- ===================================
-- 代码生成器相关表
-- ===================================

-- 代码生成配置表
CREATE TABLE IF NOT EXISTS gen_config (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    table_name VARCHAR(100) NOT NULL COMMENT '表名称',
    table_prefix VARCHAR(50) COMMENT '表前缀',
    business_name VARCHAR(100) COMMENT '业务名称',
    module_name VARCHAR(100) COMMENT '模块名称',
    package_name VARCHAR(200) COMMENT '包路径',
    author VARCHAR(50) COMMENT '作者',
    remark VARCHAR(500) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_table_name (table_name),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代码生成配置表';

-- 表字段信息表
CREATE TABLE IF NOT EXISTS gen_table_column (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '字段ID',
    gen_config_id BIGINT NOT NULL COMMENT '配置ID',
    column_name VARCHAR(100) NOT NULL COMMENT '字段名称',
    column_comment VARCHAR(500) COMMENT '字段注释',
    column_type VARCHAR(100) COMMENT '字段类型(MySQL类型)',
    go_type VARCHAR(100) COMMENT 'Go类型',
    go_field VARCHAR(100) COMMENT 'Go字段名',
    is_pk TINYINT DEFAULT 0 COMMENT '是否主键(1是)',
    is_increment TINYINT DEFAULT 0 COMMENT '是否自增(1是)',
    is_required TINYINT DEFAULT 0 COMMENT '是否必填(1是)',
    is_insert TINYINT DEFAULT 1 COMMENT '是否为插入字段(1是)',
    is_edit TINYINT DEFAULT 1 COMMENT '是否编辑字段(1是)',
    is_list TINYINT DEFAULT 1 COMMENT '是否列表字段(1是)',
    is_query TINYINT DEFAULT 1 COMMENT '是否查询字段(1是)',
    query_type VARCHAR(50) DEFAULT '=' COMMENT '查询方式(=,!=,>,<,LIKE等)',
    html_type VARCHAR(50) COMMENT '显示类型(input,textarea,select,radio,checkbox,datetime等)',
    dict_type VARCHAR(100) COMMENT '字典类型',
    sort INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (gen_config_id) REFERENCES gen_config(id) ON DELETE CASCADE,
    INDEX idx_gen_config_id (gen_config_id),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代码生成表字段信息';

-- 代码生成历史表
CREATE TABLE IF NOT EXISTS gen_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '历史ID',
    gen_config_id BIGINT NOT NULL COMMENT '配置ID',
    table_name VARCHAR(100) NOT NULL COMMENT '表名称',
    file_path VARCHAR(500) COMMENT '生成的文件路径',
    file_type VARCHAR(50) COMMENT '文件类型(api,model,handler,logic,repository)',
    content LONGTEXT COMMENT '生成的文件内容',
    status TINYINT DEFAULT 1 COMMENT '状态(1成功 0失败)',
    error_msg TEXT COMMENT '错误信息',
    operator VARCHAR(50) COMMENT '操作人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
    FOREIGN KEY (gen_config_id) REFERENCES gen_config(id) ON DELETE CASCADE,
    INDEX idx_gen_config_id (gen_config_id),
    INDEX idx_table_name (table_name),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代码生成历史表';

-- 插入开发工具菜单
INSERT INTO menus (parent_id, menu_name, menu_path, component, icon, sort, status, menu_type, remark)
VALUES 
(0, '开发工具', '/devtools', '', 'tool', 20, 1, 1, '开发工具菜单'),
(LAST_INSERT_ID(), '代码生成', '/devtools/codegen', 'devtools/codegen/CodeGen', 'code', 1, 1, 1, '代码生成器'),
(LAST_INSERT_ID() - 1, '生成历史', '/devtools/history', 'devtools/codegen/GenHistory', 'history', 2, 1, 1, '代码生成历史');
