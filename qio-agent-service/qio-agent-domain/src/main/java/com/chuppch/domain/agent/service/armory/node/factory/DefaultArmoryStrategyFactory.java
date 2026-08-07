package com.chuppch.domain.agent.service.armory.node.factory;

import com.chuppch.domain.agent.model.entity.ArmoryCommandEntity;
import com.chuppch.domain.agent.service.armory.node.RootNode;
import cn.bugstack.wrench.design.framework.tree.StrategyHandler;
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
 * @create 2025/12/16
 */
@Service
public class DefaultArmoryStrategyFactory {

    private final RootNode rootNode;

    public DefaultArmoryStrategyFactory(RootNode rootNode) {
        this.rootNode = rootNode;
    }

    public StrategyHandler <ArmoryCommandEntity,DefaultArmoryStrategyFactory.DynamicContext,String> armoryStrategyHandler(){
        return rootNode;
    }

    /**
     * 动态上下文对象
     * <p>
     * 用于在策略处理器执行过程中传递和共享上下文数据。
     * 通过键值对的方式存储任意类型的数据，支持策略链中不同节点之间的数据传递。
     * </p>
     *
     */
    @Data
    @Builder
    @AllArgsConstructor
    @NoArgsConstructor
    public static class DynamicContext {

        // 统一上下文传递对象的数据格式
        private Map<String, Object> dataObjects = new HashMap<>();

        public <T> void setValue(String key, T value) {
            dataObjects.put(key, value);
        }

        public <T> T getValue(String key) {
            return (T) dataObjects.get(key);
        }
    }

}