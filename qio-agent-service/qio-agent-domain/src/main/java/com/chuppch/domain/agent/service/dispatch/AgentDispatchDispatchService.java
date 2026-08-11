package com.chuppch.domain.agent.service.dispatch;

import com.alibaba.fastjson2.JSON;
import com.chuppch.domain.agent.adapter.repository.IAgentRepository;
import com.chuppch.domain.agent.model.entity.AutoAgentExecuteResultEntity;
import com.chuppch.domain.agent.model.entity.ExecuteCommandEntity;
import com.chuppch.domain.agent.model.valobj.AiAgentVO;
import com.chuppch.domain.agent.service.IAgentDispatchService;
import com.chuppch.domain.agent.service.IExecuteStrategy;
import com.chuppch.types.exception.BizException;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.mvc.method.annotation.ResponseBodyEmitter;

import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ThreadPoolExecutor;

/**
 * @author chuppch
 * @description
 * @create 2026/1/5
 */
@Slf4j
@Service
public class AgentDispatchDispatchService implements IAgentDispatchService {

    /**
     * 策略映射 - 三个大模型执行策略
     */
    @Resource
    private Map<String, IExecuteStrategy> executeStrategyMap;

    @Resource
    private IAgentRepository repository;

    @Resource
    private ThreadPoolExecutor threadPoolExecutor;

    @Override
    public void dispatch(ExecuteCommandEntity requestParameter, ResponseBodyEmitter emitter) throws Exception {
        if (requestParameter == null) {
            throw new BizException("ILLEGAL_PARAMETER", "执行参数不能为空");
        }
        if (requestParameter.getAiAgentId() == null || requestParameter.getAiAgentId().isBlank()) {
            throw new BizException("ILLEGAL_PARAMETER", "aiAgentId 不能为空");
        }
        if (requestParameter.getMessage() == null || requestParameter.getMessage().isBlank()) {
            throw new BizException("ILLEGAL_PARAMETER", "message 不能为空");
        }
        if (emitter == null) {
            throw new BizException("ILLEGAL_PARAMETER", "流式响应对象不能为空");
        }
        if (requestParameter.getSessionId() == null || requestParameter.getSessionId().isBlank()) {
            requestParameter.setSessionId(UUID.randomUUID().toString());
        }

        AiAgentVO aiAgentVO = repository.queryAiAgentByAgentId(requestParameter.getAiAgentId());
        if (aiAgentVO == null) {
            throw new BizException("AGENT_NOT_FOUND", "未找到指定智能体");
        }

        String strategy = aiAgentVO.getStrategy();
        if (strategy == null || strategy.isBlank()) {
            throw new BizException("AGENT_STRATEGY_NOT_CONFIGURED", "智能体未配置执行策略");
        }
        IExecuteStrategy executeStrategy = executeStrategyMap.get(strategy);
        if (null == executeStrategy) {
            throw new BizException("AGENT_STRATEGY_NOT_FOUND", "不存在的执行策略类型 strategy: " + strategy);
        }

        // 3. 异步执行
        threadPoolExecutor.execute(() -> {
            try {
                executeStrategy.execute(requestParameter, emitter);
            } catch (Exception e) {
                log.error("Agent 执行异常，agentId: {}, sessionId: {}",
                        requestParameter.getAiAgentId(), requestParameter.getSessionId(), e);
                try {
                    AutoAgentExecuteResultEntity errorResult = AutoAgentExecuteResultEntity.createErrorResult(
                            "Agent 执行失败，请稍后重试",
                            requestParameter.getSessionId());
                    emitter.send("data: " + JSON.toJSONString(errorResult) + "\n\n");
                } catch (Exception ex) {
                    log.error("发送异常信息失败：{}", ex.getMessage(), ex);
                }
            } finally {
                try {
                    emitter.complete();
                } catch (Exception e) {
                    log.error("完成流式输出失败：{}", e.getMessage(), e);
                }
            }
        });
    }
}
