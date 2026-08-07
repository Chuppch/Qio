package com.chuppch.domain.agent.service.excutor.auto.step.factory;

import cn.bugstack.wrench.design.framework.tree.StrategyHandler;
import com.chuppch.domain.agent.model.entity.ExecuteCommandEntity;
import com.chuppch.domain.agent.model.valobj.AiAgentClientFlowConfigVO;
import com.chuppch.domain.agent.service.excutor.auto.step.RootNode;
import com.chuppch.domain.agent.service.excutor.auto.vo.ExecutionHistoryManager;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;

/**
 * @author chuppch
 * @description
 * @create 2025/12/19
 */
@Service
public class DefaultAutoAgentExecuteStrategyFactory {

    private final RootNode executeRootNode;

    public DefaultAutoAgentExecuteStrategyFactory(RootNode executeRootNode) {
        this.executeRootNode = executeRootNode;
    }

    public StrategyHandler<ExecuteCommandEntity, DefaultAutoAgentExecuteStrategyFactory.DynamicContext, String> armoryStrategyHandler(){
        return executeRootNode;
    }

    @Data
    @Builder
    @AllArgsConstructor
    @NoArgsConstructor
    public static class DynamicContext {

        /** 当前执行步骤号 */
        private int step = 1;

        /** 最大任务步骤数（防止无限循环） */
        private int maxStep = 1;

        /** 执行历史记录管理器（实现滑动窗口和分层存储） */
        private ExecutionHistoryManager executionHistoryManager;

        /** 当前任务描述（根据质量监督反馈动态更新） */
        private String currentTask;

        /** 任务完成标志（质量检查通过后设置为 true） */
        boolean isCompleted = false;

        /** AI 客户端配置映射表（key=客户端类型，value=配置信息） */
        private Map<String, AiAgentClientFlowConfigVO> aiAgentClientFlowConfigVOMap;

        /** 动态数据存储（用于节点间传递数据，如 emitter、analysisResult 等） */
        private Map<String, Object> dataObjects = new HashMap<>();

        /** 设置键值对数据 */
        public <T> void setValue(String key, T value) {
            dataObjects.put(key, value);
        }

        /** 获取键值对数据 */
        public <T> T getValue(String key) {
            return (T) dataObjects.get(key);
        }
    }

}