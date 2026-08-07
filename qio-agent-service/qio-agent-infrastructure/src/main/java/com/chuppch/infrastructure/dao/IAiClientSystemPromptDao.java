package com.chuppch.infrastructure.dao;

import com.chuppch.infrastructure.dao.po.AiClientSystemPrompt;
import org.apache.ibatis.annotations.Mapper;

import java.util.List;

/**
 * 系统提示词配置表 DAO
 * @author chuppch
 * @description 系统提示词配置表数据访问对象
 */
@Mapper
public interface IAiClientSystemPromptDao {

    /**
     * 插入系统提示词配置
     * @param aiClientSystemPrompt 系统提示词配置对象
     * @return 影响行数
     */
    int insert(AiClientSystemPrompt aiClientSystemPrompt);

    /**
     * 根据ID更新系统提示词配置
     * @param aiClientSystemPrompt 系统提示词配置对象
     * @return 影响行数
     */
    int updateById(AiClientSystemPrompt aiClientSystemPrompt);

    /**
     * 根据提示词ID更新系统提示词配置
     * @param aiClientSystemPrompt 系统提示词配置对象
     * @return 影响行数
     */
    int updateByPromptId(AiClientSystemPrompt aiClientSystemPrompt);

    /**
     * 根据ID删除系统提示词配置
     * @param id 主键ID
     * @return 影响行数
     */
    int deleteById(Long id);

    /**
     * 根据提示词ID删除系统提示词配置
     * @param promptId 提示词ID
     * @return 影响行数
     */
    int deleteByPromptId(String promptId);

    /**
     * 根据ID查询系统提示词配置
     * @param id 主键ID
     * @return 系统提示词配置对象
     */
    AiClientSystemPrompt queryById(Long id);

    /**
     * 根据提示词ID查询系统提示词配置
     * @param promptId 提示词ID
     * @return 系统提示词配置对象
     */
    AiClientSystemPrompt queryByPromptId(String promptId);

    /**
     * 查询启用的系统提示词配置
     * @return 系统提示词配置列表
     */
    List<AiClientSystemPrompt> queryEnabledPrompts();

    /**
     * 根据提示词名称查询系统提示词配置
     * @param promptName 提示词名称
     * @return 系统提示词配置列表
     */
    List<AiClientSystemPrompt> queryByPromptName(String promptName);

    /**
     * 查询所有系统提示词配置
     * @return 系统提示词配置列表
     */
    List<AiClientSystemPrompt> queryAll();

}

