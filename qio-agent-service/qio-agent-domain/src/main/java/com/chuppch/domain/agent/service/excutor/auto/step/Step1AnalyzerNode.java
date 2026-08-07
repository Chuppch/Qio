package com.chuppch.domain.agent.service.excutor.auto.step;

import cn.bugstack.wrench.design.framework.tree.StrategyHandler;
import com.chuppch.domain.agent.model.entity.AutoAgentExecuteResultEntity;
import com.chuppch.domain.agent.model.entity.ExecuteCommandEntity;
import com.chuppch.domain.agent.model.valobj.AiAgentClientFlowConfigVO;
import com.chuppch.domain.agent.model.valobj.enums.AiClientTypeEnumVO;
import com.chuppch.domain.agent.service.excutor.auto.step.factory.DefaultAutoAgentExecuteStrategyFactory;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.stereotype.Service;

/**
 * @author chuppch
 * @description
 * @create 2025/12/19
 */
@Service
public class Step1AnalyzerNode extends AbstractExecuteSupport {

    @Override
    protected String doApply(ExecuteCommandEntity requestParameter, DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext) throws Exception {
        log.info("\n === 执行第 {} 步 ===", dynamicContext.getStep());

        // 根据客户端类型 - 获取配置信息
        AiAgentClientFlowConfigVO aiAgentClientFlowConfigVO = dynamicContext.getAiAgentClientFlowConfigVOMap().get(AiClientTypeEnumVO.TASK_ANALYZER_CLIENT.getCode());

        // 构建任务分析提示词
        log.info("\n 阶段1: 任务状态分析");
        String analysisPrompt = String.format(aiAgentClientFlowConfigVO.getStepPrompt(),
                requestParameter.getMessage(),
                dynamicContext.getStep(),
                dynamicContext.getMaxStep(),
                dynamicContext.getExecutionHistoryManager().getHistory(),
                dynamicContext.getCurrentTask()
        );

        // 获取 AgentClient 客户端
        ChatClient chatClient = getChatClientByClientId(aiAgentClientFlowConfigVO.getClientId());

        // 调用客户端执行 - 获取分析结果
        String analysisResult = chatClient
                .prompt(analysisPrompt)
                .advisors(a -> a
                        .param(CHAT_MEMORY_CONVERSATION_ID_KEY, requestParameter.getSessionId())
                        .param(CHAT_MEMORY_RETRIEVE_SIZE_KEY, 1024))
                .call().content();

        assert analysisResult != null;

        // 解析分析结果 - 流式输出给前端
        parseAnalysisResult(dynamicContext, analysisResult, requestParameter.getSessionId());

        // 将分析结果保存到动态上下文中，供下一步使用
        dynamicContext.setValue("analysisResult", analysisResult);

        // 检查是否已经完成
        if (analysisResult.contains("任务状态: COMPLETED") ||
                analysisResult.contains("完成度评估: 100%")) {
            dynamicContext.setCompleted(true);
            log.info("\n 任务已完成，任务描述: {}", dynamicContext.getCurrentTask());
        }

        return router(requestParameter, dynamicContext);
    }

    @Override
    public StrategyHandler<ExecuteCommandEntity, DefaultAutoAgentExecuteStrategyFactory.DynamicContext, String> get(ExecuteCommandEntity executeCommandEntity, DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext) throws Exception {
        // 如果任务已完成或达到最大步数，进入总结阶段
        if (dynamicContext.isCompleted() || dynamicContext.getStep() > dynamicContext.getMaxStep()) {
            return getBean("step4LogExecutionSummaryNode");
        }

        // 否则继续执行下一步
        return getBean("step2PrecisionExecutorNode");
    }

    // 状态机
    private void parseAnalysisResult(DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext, String analysisResult, String sessionId) {
        int step = dynamicContext.getStep();
        log.info("\n === 第 {} 步分析结果 ===", step);

        String[] lines = analysisResult.split("\n"); // 切割大模型输出文本 - 以每一行为一个值
        String currentSection = "";
        StringBuilder sectionContent = new StringBuilder();

        for (String line : lines) {
            line = line.trim(); // 去除字符串开头和结尾的空白字符
            if (line.isEmpty()) continue;

            if (line.contains("任务状态分析:")) {
                // 识别到"任务状态分析:"章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"analysis_status"
                // 3. 清空累积器，准备收集新章节内容
                sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "analysis_status";
                sectionContent = new StringBuilder();
                log.info("\n🎯 任务状态分析:");
                continue;

            } else if (line.contains("执行历史评估:")) {
                // 识别到"执行历史评估:"章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"analysis_history"
                // 3. 清空累积器，准备收集新章节内容
                sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "analysis_history";
                sectionContent = new StringBuilder();
                log.info("\n📈 执行历史评估:");
                continue;

            } else if (line.contains("下一步策略:")) {
                // 识别到"下一步策略:"章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"analysis_strategy"
                // 3. 清空累积器，准备收集新章节内容
                sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "analysis_strategy";
                sectionContent = new StringBuilder();
                log.info("\n🚀 下一步策略:");
                continue;

            } else if (line.contains("完成度评估:")) {
                // 识别到"完成度评估:"章节标题（标题行包含关键数据：百分比）：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"analysis_progress"
                // 3. 清空累积器，准备收集新章节内容
                // 4. 提取百分比值用于日志输出
                // 5. 保存标题行（包含关键数据）
                sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "analysis_progress";
                sectionContent = new StringBuilder();

                // 提取百分比值用于日志输出
                String progress = line.substring(line.indexOf(":") + 1).trim();
                log.info("\n📊 完成度评估: {}", progress);

                // 保存包含关键数据的标题行
                sectionContent.append(line).append("\n");
                continue;

            } else if (line.contains("任务状态:")) {
                // 识别到"任务状态:"章节标题（标题行包含关键数据：状态值）：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"analysis_task_status"  ← ✅ 修复错误
                // 3. 清空累积器，准备收集新章节内容
                // 4. 提取状态值用于日志判断和输出
                // 5. 保存标题行（包含关键数据）
                sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "analysis_task_status";  // ← ✅ 修复：从 "analysis_progress" 改为 "analysis_task_status"
                sectionContent = new StringBuilder();

                // 提取状态值用于日志判断
                String status = line.substring(line.indexOf(":") + 1).trim();
                if (status.equals("COMPLETED")) {
                    log.info("\n✅ 任务状态: 已完成");
                } else {
                    log.info("\n🔄 任务状态: 继续执行");
                }

                // 保存包含关键数据的标题行
                sectionContent.append(line).append("\n");
                continue;
            }

            // 收集当前section的内容
            if (!currentSection.isEmpty()) {
                sectionContent.append(line).append("\n");
                switch (currentSection) {
                    case "analysis_status":
                        log.info("   📋 {}", line);
                        break;
                    case "analysis_history":
                        log.info("   📊 {}", line);
                        break;
                    case "analysis_strategy":
                        log.info("   🎯 {}", line);
                        break;
                    default:
                        log.info("   📝 {}", line);
                        break;
                }
            }
        }
        // 发送最后一个section的内容
        sendAnalysisSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
    }

    private void sendAnalysisSubResult(DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext,
                                       String subType, String string, String sessionId) {
        if (!subType.isEmpty() && !string.isEmpty()) {
            AutoAgentExecuteResultEntity result = AutoAgentExecuteResultEntity.createAnalysisSubResult(
                    dynamicContext.getStep(), subType, string, sessionId
            );
            sendSseResult(dynamicContext,result);
        }
    }
}