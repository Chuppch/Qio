package com.chuppch.domain.agent.service.armory.business.data;

import com.chuppch.domain.agent.model.entity.ArmoryCommandEntity;
import com.chuppch.domain.agent.service.armory.node.factory.DefaultArmoryStrategyFactory;

/**
 * @author chuppch
 * @description 数据加载策略
 * @create 2025/12/16
 */
public interface ILoadDataStrategy {

    void loadData(ArmoryCommandEntity armoryCommandEntity, DefaultArmoryStrategyFactory.DynamicContext dynamicContext);

}