package com.chuppch.api.admin;

import com.chuppch.api.dto.DataStatisticsResponseDTO;
import com.chuppch.api.response.Response;

/**
 * @author chuppch
 * @description
 * @create 2026/1/6
 */
public interface IAiAgentDataStatisticsAdminService {

    /**
     * 获取系统数据统计
     * @return 统计数据响应
     */
    Response<DataStatisticsResponseDTO> getDataStatistics();

}