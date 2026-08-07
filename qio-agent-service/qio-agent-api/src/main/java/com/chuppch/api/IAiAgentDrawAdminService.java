package com.chuppch.api;

import com.chuppch.api.dto.AiAgentDrawConfigQueryRequestDTO;
import com.chuppch.api.dto.AiAgentDrawConfigRequestDTO;
import com.chuppch.api.dto.AiAgentDrawConfigResponseDTO;
import com.chuppch.api.response.Response;

import java.util.List;

/**
 * @author chuppch
 * @description
 * @create 2026/1/6
 */
public interface IAiAgentDrawAdminService {

    /**
     * 保存拖拉拽流程图配置
     *
     * @param request 配置请求参数
     * @return 保存结果
     */
    Response<String> saveDrawConfig(AiAgentDrawConfigRequestDTO request);

    /**
     * 获取拖拉拽流程图配置
     *
     * @param configId 配置ID
     * @return 配置数据
     */
    Response<AiAgentDrawConfigResponseDTO> getDrawConfig(String configId);

    /**
     * 分页查询拖拉拽流程图配置列表
     *
     * @param request 查询条件与分页参数
     * @return 配置列表
     */
    Response<List<AiAgentDrawConfigResponseDTO>> queryDrawConfigList(AiAgentDrawConfigQueryRequestDTO request);

    /**
     * 删除拖拉拽流程图配置
     *
     * @param configId 配置ID
     * @return 删除结果
     */
    Response<String> deleteDrawConfig(String configId);

}
