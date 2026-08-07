package com.chuppch.domain.agent.adapter.repository;

import com.chuppch.domain.agent.model.valobj.*;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Map;

/**
 * @author chuppch
 * @description
 * @create 2025/12/17
 */
public interface IAgentRepository {

    // ============ AiClientApiLoadDataStrategy - 数据库操作 ============

    List<AiClientApiVO> queryAiClientApiVOListByApiIds(List<String> apiIdList);

    // ============ AiClientLoadDataStrategy - 数据库操作 ============

    /**
     * 根据客户端ID列表，查询对应的 Api 信息（只返回已启用的 API）
     *
     * @param clientIdList 客户端ID列表
     * @return Api信息VO列表
     */
    List<AiClientApiVO> queryAiClientApiVOListByClientIds(List<String> clientIdList);

    /**
     * 根据客户端ID列表，查询对应的模型信息（只返回已启用的模型）
     *
     * @param clientIdList 客户端ID列表
     * @return 模型信息VO列表
     */
    List<AiClientModelVO> AiClientModelVOByClientIds(List<String> clientIdList);

    /**
     * 根据客户端ID列表，查询对应的工具信息（只返回已启用的工具）
     *
     * @param clientIdList 客户端ID列表
     * @return 工具信息VO列表
     */
    List<AiClientToolMcpVO> AiClientToolMcpVOByClientIds(List<String> clientIdList);

    /**
     * 根据客户端ID列表，查询对应的系统提示信息（只返回已启用的系统提示）
     *
     * @param clientIdList 客户端ID列表
     * @return 系统提示信息VO列表
     */
    Map<String, AiClientSystemPromptVO> queryAiClientSystemPromptMapByClientIds(List<String> clientIdList);

    /**
     * 根据客户端ID列表，查询对应的顾问信息（只返回已启用的顾问）
     *
     * @param clientIdList 客户端ID列表
     * @return 顾问信息VO列表
     */
    List<AiClientAdvisorVO> AiClientAdvisorVOByClientIds(List<String> clientIdList);

    /**
     * 根据客户端ID列表，查询对应的客户端信息（只返回已启用的客户端）
     *
     * @param clientIdList 客户端ID列表
     * @return 客户端信息VO列表
     */
    List<AiClientVO> AiClientVOByClientIds(List<String> clientIdList);

    // ============ AiClientModelLoadDataStrategy - 数据库操作 ============

    /**
     * 根据模型ID列表，查询对应的 Api 信息（只返回已启用的 API）
     * @param modelIdList
     * @return
     */
    List<AiClientApiVO> queryAiClientApiVOListByModelIds(List<String> modelIdList);

    /**
     * 根据模型ID列表，查询对应的模型信息（只返回已启用的模型）
     * @param modelIdList
     * @return
     */
    List<AiClientModelVO> AiClientModelVOByModelIds(List<String> modelIdList);

    // ============ ArmoryService - 数据库操作 ============

    /**
      * 根据代理ID，查询对应的客户端信息（只返回已启用的客户端）
      * @param agentId
      * @return
      */
    List<AiAgentClientFlowConfigVO> queryAiAgentClientsByAgentId(String aiAgentId);

    /**
      * 查询所有可用的代理
      * @return
      */
    List<AiAgentVO> queryAvailableAgents();

    // ============ 额外 - 数据库操作 ============

    /**
      * 根据客户端ID列表，查询对应的系统提示信息（只返回已启用的系统提示）
      * @param clientIdList
      * @return
      */
    List<AiClientSystemPromptVO> AiClientSystemPromptVOByClientIds(List<String> clientIdList);

    // ============ execute - RootNode 数据库操作============

    /**
      * 根据代理ID，查询对应的客户端流程配置信息
      * @param aiAgentId
      * @return
      */
    Map<String, AiAgentClientFlowConfigVO> queryAiAgentClientFlowConfig(String aiAgentId);

    // =========== AgentDispatchDispatchService ============

    /**
     * 根据代理ID，查询对应的agent信息
     * @param aiAgentId
     * @return
     */
    AiAgentVO queryAiAgentByAgentId(String aiAgentId);

    // =========== RagService -  rag 数据库操作 ============

    void createTagOrder(AiRagOrderVO aiRagOrderVO);
}
