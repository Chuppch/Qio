-- 关闭外键检查（安全起见，先关掉）
SET FOREIGN_KEY_CHECKS = 0;

-- 按功能模块分组删除（顺序其实无所谓，这里只是方便阅读）

-- 管理员与智能体
DROP TABLE IF EXISTS `admin_user`;
DROP TABLE IF EXISTS `ai_agent`;
DROP TABLE IF EXISTS `ai_agent_draw_config`;
DROP TABLE IF EXISTS `ai_agent_flow_config`;
DROP TABLE IF EXISTS `ai_agent_task_schedule`;

-- 客户端配置相关
DROP TABLE IF EXISTS `ai_client`;
DROP TABLE IF EXISTS `ai_client_advisor`;
DROP TABLE IF EXISTS `ai_client_api`;
DROP TABLE IF EXISTS `ai_client_config`;
DROP TABLE IF EXISTS `ai_client_model`;
DROP TABLE IF EXISTS `ai_client_rag_order`;
DROP TABLE IF EXISTS `ai_client_system_prompt`;
DROP TABLE IF EXISTS `ai_client_tool_mcp`;

-- 打开外键检查
SET FOREIGN_KEY_CHECKS = 1;