package com.chuppch.domain.agent.service.excutor.auto.vo;

import lombok.Data;

import java.io.Serializable;

/**
 * 步骤记录
 * 表示单步执行的完整记录
 */
@Data
public class StepRecord implements Serializable {
    
    /** 步数 */
    private int stepNumber;
    
    /** 分析阶段的结果 */
    private String analysisResult;
    
    /** 执行阶段的结果 */
    private String executionResult;
    
    /** 监督阶段的结果 */
    private String supervisionResult;
    
    /** 执行时间戳 */
    private long timestamp;
    
    /** 格式化后的完整记录 */
    private String stepSummary;
}

