package com.chuppch.domain.agent.service;

import com.chuppch.domain.agent.model.valobj.AiAgentVO;

import java.util.List;

/**
 * @author chuppch
 * @description
 * @create 2025/12/17
 */
public interface IArmoryService {

    // 装配所有可用的智能体
    List<AiAgentVO> acceptArmoryAllAvailableAgents();

    // 装配单个智能体
    void acceptArmoryAgent(String agentId);

    // 查询所有可用的智能体
    List<AiAgentVO> queryAvailableAgents();

    // 装配智能体客户端模型API
    void acceptArmoryAgentClientModelApi(String apiId);

}
