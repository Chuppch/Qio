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
@Service("step2PrecisionExecutorNode")
public class Step2PrecisionExecutorNode extends AbstractExecuteSupport {
    @Override
    protected String doApply(ExecuteCommandEntity requestParameter, DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext) throws Exception {
        log.info("\n 阶段2: 精准任务执行");

        // 从动态上下文中获取分析结果
        String analysisResult = dynamicContext.getValue("analysisResult");
        if (analysisResult == null || analysisResult.isEmpty()) {
            log.warn(" 分析结果为空，使用默认执行策略");
            analysisResult = "执行当前任务步骤";
        }

        // 获取客户端配置
        AiAgentClientFlowConfigVO aiAgentClientFlowConfigVO = dynamicContext.getAiAgentClientFlowConfigVOMap().get(AiClientTypeEnumVO.PRECISION_EXECUTOR_CLIENT.getCode());

        // 构建执行提示词
        String executionPrompt = String.format(aiAgentClientFlowConfigVO.getStepPrompt(),
                requestParameter.getMessage(),
                analysisResult
        );

        // 获取 AgentClient 客户端
        ChatClient chatClient = getChatClientByClientId(aiAgentClientFlowConfigVO.getClientId());

        // 调用客户端执行 - 获取分析结果
        String executionResult = chatClient
                .prompt(executionPrompt)
                .advisors(a -> a
                        .param(CHAT_MEMORY_CONVERSATION_ID_KEY, requestParameter.getSessionId())
                        .param(CHAT_MEMORY_RETRIEVE_SIZE_KEY, 1024))
                .call().content();

        assert executionResult != null;

        parseExecutionResult(dynamicContext, executionResult, requestParameter.getSessionId());

        dynamicContext.setValue("executionResult", executionResult);

        // 注意：此阶段不追加历史记录
        // 历史记录将在 Step3QualitySupervisorNode 中统一追加（包含完整的三个阶段结果）

        return router(requestParameter, dynamicContext);
    }

    @Override
    public StrategyHandler<ExecuteCommandEntity, DefaultAutoAgentExecuteStrategyFactory.DynamicContext, String> get(ExecuteCommandEntity executeCommandEntity, DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext) throws Exception {
        return getBean("step3QualitySupervisorNode");
    }

    /**
     * 解析执行结果
     */
    private void parseExecutionResult(DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext, String executionResult, String sessionId) {
        int step = dynamicContext.getStep();
        log.info("\n⚡ === 第 {} 步执行结果 ===", step);

        String[] lines = executionResult.split("\n");
        String currentSection = "";
        StringBuilder sectionContent = new StringBuilder();

        for (String line : lines) {
            line = line.trim();
            if (line.isEmpty()) continue;

            if (line.contains("执行目标:")) {
                // 识别到 "执行目标:" 章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"execution_target"
                // 3. 清空累积器，准备收集新章节内容
                sendExecutionSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "execution_target";
                sectionContent = new StringBuilder();
                log.info("\n 执行目标:");
                continue;
            } else if (line.contains("执行过程:")) {
                // 识别到 "执行过程:" 章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"execution_process"
                // 3. 清空累积器，准备收集新章节内容
                sendExecutionSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "execution_process";
                sectionContent = new StringBuilder();
                log.info("\n 执行过程:");
                continue;
            } else if (line.contains("执行结果:")) {
                // 识别到 "执行结果:" 章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"execution_result"
                // 3. 清空累积器，准备收集新章节内容
                sendExecutionSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "execution_result";
                sectionContent = new StringBuilder();
                log.info("\n 执行结果:");
                continue;
            } else if (line.contains("质量检查:")) {
                // 识别到 "质量检查:" 章节标题：
                // 1. 发送上一个章节的累积内容（如果有）
                // 2. 切换章节类型标识为"quality_check"
                // 3. 清空累积器，准备收集新章节内容
                sendExecutionSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
                currentSection = "execution_quality";
                sectionContent = new StringBuilder();
                log.info("\n 质量检查:");
                continue;
            }
            // 收集当前section的内容
            if (!currentSection.isEmpty()) {
                sectionContent.append(line).append("\n");
                switch (currentSection) {
                    case "execution_target":
                        log.info("   🎯 {}", line);
                        break;
                    case "execution_process":
                        log.info("   ⚙️ {}", line);
                        break;
                    case "execution_result":
                        log.info("   📊 {}", line);
                        break;
                    case "execution_quality":
                        log.info("   ✅ {}", line);
                        break;
                    default:
                        log.info("   📝 {}", line);
                        break;
                }
            }
        }
        // 发送最后一个section的内容
        sendExecutionSubResult(dynamicContext, currentSection, sectionContent.toString(), sessionId);
    }

    /**
     * 发送执行阶段细分结果到流式输出
     */
    private void sendExecutionSubResult(DefaultAutoAgentExecuteStrategyFactory.DynamicContext dynamicContext,
                                        String subType, String content, String sessionId) {
        // 抽取的通用判断逻辑
        if (!subType.isEmpty() && !content.isEmpty()) {
            AutoAgentExecuteResultEntity result = AutoAgentExecuteResultEntity.createExecutionSubResult(
                    dynamicContext.getStep(), subType, content, sessionId);
            sendSseResult(dynamicContext, result);
        }
    }
}
