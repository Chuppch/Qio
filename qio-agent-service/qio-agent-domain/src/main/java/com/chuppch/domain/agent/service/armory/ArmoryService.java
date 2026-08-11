package com.chuppch.domain.agent.service.armory;

import cn.bugstack.wrench.design.framework.tree.StrategyHandler;
import com.chuppch.domain.agent.adapter.repository.IAgentRepository;
import com.chuppch.domain.agent.model.entity.ArmoryCommandEntity;
import com.chuppch.domain.agent.model.valobj.AiAgentClientFlowConfigVO;
import com.chuppch.domain.agent.model.valobj.AiAgentVO;
import com.chuppch.domain.agent.model.valobj.enums.AiAgentEnumVO;
import com.chuppch.domain.agent.service.IArmoryService;
import com.chuppch.domain.agent.service.armory.node.factory.DefaultArmoryStrategyFactory;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

/**
 * @author chuppch
 * @description
 * @create 2025/12/16
 */
@Slf4j
@Service
public class ArmoryService implements IArmoryService {

    @Resource
    private IAgentRepository repository;

    @Resource
    private DefaultArmoryStrategyFactory defaultArmoryStrategyFactory;

    @Override
    public List<AiAgentVO> acceptArmoryAllAvailableAgents() {
        // 获取所有可用的智能体
        List<AiAgentVO> aiAgentVOS = repository.queryAvailableAgents();
        if (aiAgentVOS == null || aiAgentVOS.isEmpty()) {
            log.info("当前没有可装配的智能体");
            return Collections.emptyList();
        }

        // 循环装配智能体
        for (AiAgentVO aiAgentVO : aiAgentVOS) {
            if (aiAgentVO == null || aiAgentVO.getAgentId() == null || aiAgentVO.getAgentId().isBlank()) {
                log.warn("跳过缺少 agentId 的智能体配置");
                continue;
            }
            String agentId = aiAgentVO.getAgentId();
            try {
                acceptArmoryAgent(agentId);
            } catch (Exception e) {
                // 单个智能体装配失败不应阻断其余可用智能体的初始化。
                log.error("智能体装配失败，agentId: {}", agentId, e);
            }
        }
        return aiAgentVOS;
    }

    @Override
    public void acceptArmoryAgent(String agentId) {
        if (agentId == null || agentId.isBlank()) {
            throw new IllegalArgumentException("agentId 不能为空");
        }

        List<AiAgentClientFlowConfigVO> aiAgentClientFlowConfigVOS = repository.queryAiAgentClientsByAgentId(agentId);
        if (aiAgentClientFlowConfigVOS == null || aiAgentClientFlowConfigVOS.isEmpty()) {
            throw new IllegalStateException("智能体未配置可装配的客户端，agentId: " + agentId);
        }

        // 获取命令ID列表
        List<String> commandIdList = aiAgentClientFlowConfigVOS.stream()
                .filter(config -> config != null && config.getClientId() != null && !config.getClientId().isBlank())
                .map(AiAgentClientFlowConfigVO::getClientId)
                .collect(Collectors.toList());
        if (commandIdList.isEmpty()) {
            throw new IllegalStateException("智能体客户端配置缺少 clientId，agentId: " + agentId);
        }

        // 装配智能体
        try {
            StrategyHandler<ArmoryCommandEntity, DefaultArmoryStrategyFactory.DynamicContext, String> armoryStrategyHandler =
                    defaultArmoryStrategyFactory.armoryStrategyHandler();

            armoryStrategyHandler.apply(
                    ArmoryCommandEntity.builder()
                            .commandType(AiAgentEnumVO.AI_CLIENT.getCode()) // 区分不同策略的关键 - 执行 AI_CLIENT 策略
                            .commandIdList(commandIdList)
                            .build(),
                    new DefaultArmoryStrategyFactory.DynamicContext());
        } catch (Exception e) {
            throw new IllegalStateException("装配智能体失败，agentId: " + agentId, e);
        }
    }

    @Override
    public List<AiAgentVO> queryAvailableAgents() {
        return repository.queryAvailableAgents();
    }

    @Override
    public void acceptArmoryAgentClientModelApi(String apiId) {
        if (apiId == null || apiId.isBlank()) {
            throw new IllegalArgumentException("apiId 不能为空");
        }

        try {
            StrategyHandler<ArmoryCommandEntity, DefaultArmoryStrategyFactory.DynamicContext, String> armoryStrategyHandler =
                    defaultArmoryStrategyFactory.armoryStrategyHandler();

            armoryStrategyHandler.apply(
                    ArmoryCommandEntity.builder()
                            .commandType(AiAgentEnumVO.AI_CLIENT_API.getCode()) // 区分不同策略的关键 - 执行 AI_CLIENT_API 策略
                            .commandIdList(Collections.singletonList(apiId)) //  这里使用 singletonList 可以节约内存
                            .build(),
                    new DefaultArmoryStrategyFactory.DynamicContext()
            );
        } catch (Exception e) {
            throw new IllegalStateException("装配智能体 API 失败，apiId: " + apiId, e);
        }
    }
}
