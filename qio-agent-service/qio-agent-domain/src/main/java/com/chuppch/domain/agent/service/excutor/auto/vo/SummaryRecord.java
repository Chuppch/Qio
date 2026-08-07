package com.chuppch.domain.agent.service.excutor.auto.vo;

import lombok.Data;

import java.io.Serializable;

/**
 * 摘要记录
 * 表示压缩后的摘要信息
 */
@Data
public class SummaryRecord implements Serializable {
    
    /** 起始步数 */
    private int startStep;
    
    /** 结束步数 */
    private int endStep;
    
    /** 摘要内容 */
    private String summary;
    
    /** 生成时间戳 */
    private long timestamp;
}

