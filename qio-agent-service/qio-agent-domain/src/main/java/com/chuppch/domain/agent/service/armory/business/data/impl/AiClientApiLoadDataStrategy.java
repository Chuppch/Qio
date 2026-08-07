package com.chuppch.domain.agent.service.armory.business.data.impl;

import com.chuppch.domain.agent.adapter.repository.IAgentRepository;
import com.chuppch.domain.agent.model.entity.ArmoryCommandEntity;
import com.chuppch.domain.agent.model.valobj.AiClientApiVO;
import com.chuppch.domain.agent.model.valobj.enums.AiAgentEnumVO;
import com.chuppch.domain.agent.service.armory.business.data.ILoadDataStrategy;
import com.chuppch.domain.agent.service.armory.node.factory.DefaultArmoryStrategyFactory;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ThreadPoolExecutor;

/**
 * @author chuppch
 * @description AI客户端API配置数据加载策略
 * @create 2025/12/17
 */
@Slf4j
@Service("aiClientApiLoadDataStrategy")
public class AiClientApiLoadDataStrategy implements ILoadDataStrategy {

    // 仓储
    @Resource
    private IAgentRepository repository;

    // 线程池
    @Resource
    protected ThreadPoolExecutor threadPoolExecutor;

    // 数据加载
    @Override
    public void loadData(ArmoryCommandEntity armoryCommandEntity, DefaultArmoryStrategyFactory.DynamicContext dynamicContext) {
        List<String> apiIdList = armoryCommandEntity.getCommandIdList();

        CompletableFuture<List<AiClientApiVO>> aiClientApiListFuture = CompletableFuture.supplyAsync(() -> {
            log.info("查询配置数据(ai_client_api) {}", apiIdList);
            return repository.queryAiClientApiVOListByApiIds(apiIdList);
        },threadPoolExecutor);

        CompletableFuture.allOf(aiClientApiListFuture).thenRun(() -> {
            dynamicContext.setValue(AiAgentEnumVO.AI_CLIENT_API.getDataName(), aiClientApiListFuture.join());
        }).join();
    }

}