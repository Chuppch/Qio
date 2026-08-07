package com.chuppch.domain.agent.service.excutor.auto.vo;


import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;

/**
 * 执行历史记录管理器
 * 
 * 核心职责：
 * 1. 管理执行历史记录（详细层 + 摘要层）
 * 2. 实现滑动窗口机制，控制历史记录长度
 * 3. 自动压缩历史记录，生成摘要
 * 4. 根据需求提供不同粒度的历史信息
 */
public class ExecutionHistoryManager {
    
    // ========== 字段定义 ==========
    
    /** 最近详细记录（滑动窗口的详细层） */
    private final Deque<StepRecord> recentDetailedHistory;
    
    /** 摘要历史（压缩后的历史） */
    private final List<SummaryRecord> summaryHistory;
    
    /** 窗口大小（控制详细层保留多少步的完整记录） */
    private int windowSize;
    
    /** 压缩批次大小（每次压缩时处理多少步的记录） */
    private int compressionBatchSize;
    
    /** 总步数 */
    private int totalSteps;
    
    /** 最大摘要长度（避免摘要也无限增长） */
    private final int maxSummaryLength;
    
    // ========== 常量定义 ==========
    
    /** 默认窗口大小 */
    private static final int DEFAULT_WINDOW_SIZE = 5;
    
    /** 默认压缩批次大小 */
    private static final int DEFAULT_COMPRESSION_BATCH_SIZE = 3;
    
    /** 默认最大摘要长度 */
    private static final int DEFAULT_MAX_SUMMARY_LENGTH = 5000;
    
    /** 字段提取时的最大长度（字符数） */
    private static final int MAX_FIELD_LENGTH = 200;
    
    // ========== 构造函数 ==========
    
    /**
     * 默认构造函数
     */
    public ExecutionHistoryManager() {
        this.recentDetailedHistory = new ArrayDeque<>();
        this.summaryHistory = new ArrayList<>();
        this.windowSize = DEFAULT_WINDOW_SIZE;
        this.compressionBatchSize = DEFAULT_COMPRESSION_BATCH_SIZE;
        this.totalSteps = 0;
        this.maxSummaryLength = DEFAULT_MAX_SUMMARY_LENGTH;
    }
    
    /**
     * 带参数的构造函数
     * 
     * @param windowSize 窗口大小
     * @param compressionBatchSize 压缩批次大小
     * @param maxSummaryLength 最大摘要长度
     */
    public ExecutionHistoryManager(int windowSize, int compressionBatchSize, int maxSummaryLength) {
        this.recentDetailedHistory = new ArrayDeque<>();
        this.summaryHistory = new ArrayList<>();
        this.windowSize = windowSize;
        this.compressionBatchSize = compressionBatchSize;
        this.totalSteps = 0;
        this.maxSummaryLength = maxSummaryLength;
    }
    
    /**
     * 追加新步骤
     * 
     * @param stepNumber 步数
     * @param analysisResult 分析阶段的结果
     * @param executionResult 执行阶段的结果
     * @param supervisionResult 监督阶段的结果
     */
    public void appendStep(int stepNumber, String analysisResult, String executionResult, String supervisionResult) {
        // ========== 步骤 a: 创建 StepRecord 对象 ==========
        StepRecord record = new StepRecord();
        record.setStepNumber(stepNumber);
        record.setAnalysisResult(analysisResult);
        record.setExecutionResult(executionResult);
        record.setSupervisionResult(supervisionResult);
        record.setTimestamp(System.currentTimeMillis());
        
        // 格式化完整记录（用于 getHistory() 方法）
        String stepSummary = String.format("""
            === 第 %d 步完整记录 ===
            【分析阶段】%s
            【执行阶段】%s
            【监督阶段】%s
            """, stepNumber, analysisResult, executionResult, supervisionResult);
        record.setStepSummary(stepSummary);
        
        // ========== 步骤 b: 追加到 recentDetailedHistory ==========
        recentDetailedHistory.addLast(record);
        
        // ========== 步骤 c: 递增 totalSteps ==========
        totalSteps++;
        
        // ========== 步骤 d: 动态调整 windowSize 和 compressionBatchSize ==========
        adjustWindowSizeAndBatchSize();
        
        // ========== 步骤 e: 检查是否需要压缩 ==========
        if (shouldCompress()) {
            // ========== 步骤 f: 自动触发压缩 ==========
            compressHistory();
        }
    }
    
    /**
     * 获取完整的历史记录（摘要 + 详细记录）
     * 
     * @return 格式化的历史记录字符串
     */
    public String getHistory() {
        // 如果没有任何历史记录，返回首次执行标识
        if (summaryHistory.isEmpty() && recentDetailedHistory.isEmpty()) {
            return "[首次执行]";
        }
        
        StringBuilder historyBuilder = new StringBuilder();
        
        // ========== 格式化摘要历史 ==========
        if (!summaryHistory.isEmpty()) {
            for (SummaryRecord summaryRecord : summaryHistory) {
                historyBuilder.append(summaryRecord.getSummary()).append("\n\n");
            }
        }
        
        // ========== 格式化详细历史 ==========
        if (!recentDetailedHistory.isEmpty()) {
            for (StepRecord record : recentDetailedHistory) {
                historyBuilder.append(record.getStepSummary()).append("\n");
            }
        }
        
        return historyBuilder.toString().trim();
    }
    
    /**
     * 压缩历史记录
     * 将详细层中最旧的记录压缩为摘要，移到摘要层
     */
    private void compressHistory() {
        // ========== 步骤1: 取出需要压缩的记录 ==========
        // 确定实际要压缩的记录数量（可能不足 compressionBatchSize）
        int numRecordsToCompress = Math.min(compressionBatchSize, recentDetailedHistory.size());
        
        if (numRecordsToCompress == 0) {
            return;  // 没有记录需要压缩
        }
        
        // 从头部取出最旧的 N 条记录
        List<StepRecord> recordsToCompress = new ArrayList<>();
        for (int i = 0; i < numRecordsToCompress; i++) {
            StepRecord record = recentDetailedHistory.removeFirst();  // 从头部移除
            recordsToCompress.add(record);
        }
        
        // 记录起始和结束步数
        int startStep = recordsToCompress.get(0).getStepNumber();
        int endStep = recordsToCompress.get(recordsToCompress.size() - 1).getStepNumber();
        
        // ========== 步骤2: 提取关键信息 ==========
        List<String> keyInfoList = new ArrayList<>();
        for (StepRecord record : recordsToCompress) {
            String keyInfo = extractKeyInfo(record);
            keyInfoList.add(keyInfo);
        }
        
        // ========== 步骤3: 生成摘要 ==========
        // 格式化多条记录的关键信息
        StringBuilder summaryBuilder = new StringBuilder();
        summaryBuilder.append(String.format("步骤%d-%d执行摘要:\n", startStep, endStep));
        for (String keyInfo : keyInfoList) {
            summaryBuilder.append(keyInfo).append("\n\n");
        }
        String summary = summaryBuilder.toString().trim();
        
        // 创建 SummaryRecord 对象
        SummaryRecord summaryRecord = new SummaryRecord();
        summaryRecord.setStartStep(startStep);
        summaryRecord.setEndStep(endStep);
        summaryRecord.setSummary(summary);
        summaryRecord.setTimestamp(System.currentTimeMillis());
        
        // ========== 步骤4: 追加到摘要层 ==========
        summaryHistory.add(summaryRecord);
        
        // ========== 步骤5: 从详细层移除记录 ==========
        // 已经在步骤1中移除了，这里不需要再次移除
        
        // ========== 步骤6: 检查摘要层长度 ==========
        checkAndCompressSummaryHistory();
    }
    
    /**
     * 从步骤记录中提取关键信息
     * 
     * @param record 步骤记录
     * @return 格式化的关键信息字符串
     */
    private String extractKeyInfo(StepRecord record) {
        StringBuilder keyInfo = new StringBuilder();
        keyInfo.append(String.format("步骤%d:\n", record.getStepNumber()));
        
        // ========== 从 analysisResult 中提取 ==========
        // 提取任务状态
        String taskStatus = extractField(record.getAnalysisResult(), "**任务状态:**");
        if (taskStatus != null && !taskStatus.isEmpty()) {
            keyInfo.append("  任务状态: ").append(taskStatus.trim()).append("\n");
        } else {
            keyInfo.append("  任务状态: 未提供\n");
        }
        
        // 提取完成度评估
        String progress = extractField(record.getAnalysisResult(), "**完成度评估:**");
        if (progress != null && !progress.isEmpty()) {
            keyInfo.append("  完成度: ").append(progress.trim()).append("\n");
        } else {
            keyInfo.append("  完成度: 未提供\n");
        }
        
        // ========== 从 executionResult 中提取 ==========
        // 提取执行目标（摘要，可能截断）
        String executionTarget = extractField(record.getExecutionResult(), "**执行目标:**");
        if (executionTarget != null && !executionTarget.isEmpty()) {
            // 截断到前 MAX_FIELD_LENGTH 字符
            String targetSummary = truncate(executionTarget, MAX_FIELD_LENGTH);
            keyInfo.append("  执行目标: ").append(targetSummary).append("\n");
        } else {
            keyInfo.append("  执行目标: 未提供\n");
        }
        
        // 提取关键结果（摘要，可能截断）
        String executionResult = extractField(record.getExecutionResult(), "**执行结果:**");
        if (executionResult != null && !executionResult.isEmpty()) {
            // 截断到前 MAX_FIELD_LENGTH 字符
            String resultSummary = truncate(executionResult, MAX_FIELD_LENGTH);
            keyInfo.append("  关键结果: ").append(resultSummary).append("\n");
        } else {
            keyInfo.append("  关键结果: 未提供\n");
        }
        
        // ========== 从 supervisionResult 中提取 ==========
        // 提取质量评分
        String qualityScore = extractField(record.getSupervisionResult(), "**质量评分:**");
        if (qualityScore != null && !qualityScore.isEmpty()) {
            keyInfo.append("  质量评分: ").append(qualityScore.trim()).append("\n");
        } else {
            keyInfo.append("  质量评分: 未提供\n");
        }
        
        // 提取是否通过
        String passStatus = extractField(record.getSupervisionResult(), "**是否通过:**");
        if (passStatus != null && !passStatus.isEmpty()) {
            keyInfo.append("  是否通过: ").append(passStatus.trim()).append("\n");
        } else {
            keyInfo.append("  是否通过: 未提供\n");
        }
        
        return keyInfo.toString();
    }
    
    /**
     * 动态调整窗口大小和压缩批次大小
     * 根据 totalSteps 动态调整参数
     */
    private void adjustWindowSizeAndBatchSize() {
        // 根据 totalSteps 动态调整 windowSize
        if (totalSteps <= 10) {
            windowSize = 5;
        } else if (totalSteps <= 20) {
            windowSize = 8;
        } else if (totalSteps <= 50) {
            windowSize = 10;
        } else {
            windowSize = 12;  // 上限
        }
        
        // 根据 totalSteps 动态调整 compressionBatchSize
        if (totalSteps <= 10) {
            compressionBatchSize = 2;
        } else if (totalSteps <= 20) {
            compressionBatchSize = 3;
        } else if (totalSteps <= 50) {
            compressionBatchSize = 4;
        } else {
            compressionBatchSize = 5;  // 上限
        }
    }
    
    /**
     * 判断是否需要压缩历史记录
     * 
     * @return true 如果需要压缩，false 否则
     */
    private boolean shouldCompress() {
        // 检查详细记录数量是否超过窗口大小
        return recentDetailedHistory.size() > windowSize;
    }
    
    /**
     * 从文本中提取字段内容
     * 
     * @param text 文本内容
     * @param fieldMarker 字段标记（如 "**任务状态:**"）
     * @return 字段内容，如果未找到则返回 null
     */
    private String extractField(String text, String fieldMarker) {
        if (text == null || text.isEmpty()) {
            return null;
        }
        
        // 查找字段标记的位置
        int markerIndex = text.indexOf(fieldMarker);
        if (markerIndex == -1) {
            return null;  // 未找到字段标记
        }
        
        // 找到字段内容的起始位置（标记后）
        int contentStart = markerIndex + fieldMarker.length();
        
        // 查找下一个字段标记的位置（如果存在）
        int nextMarkerIndex = text.indexOf("**", contentStart);
        
        // 提取字段内容
        String fieldContent;
        if (nextMarkerIndex != -1) {
            // 有下一个字段，提取到下一个字段之前
            fieldContent = text.substring(contentStart, nextMarkerIndex);
        } else {
            // 没有下一个字段，提取到文本结束
            fieldContent = text.substring(contentStart);
        }
        
        // 去除首尾空白字符
        return fieldContent.trim();
    }
    
    /**
     * 截断文本到指定长度
     * 
     * @param text 文本内容
     * @param maxLength 最大长度
     * @return 截断后的文本
     */
    private String truncate(String text, int maxLength) {
        if (text == null || text.isEmpty()) {
            return "";
        }
        
        if (text.length() <= maxLength) {
            return text;
        }
        
        // 截断到指定长度，并添加省略号
        return text.substring(0, maxLength) + "...";
    }
    
    /**
     * 检查摘要层长度，如果超过最大长度则进一步压缩
     */
    private void checkAndCompressSummaryHistory() {
        // 计算摘要层的总长度
        int totalLength = 0;
        for (SummaryRecord record : summaryHistory) {
            totalLength += record.getSummary().length();
        }
        
        // 如果超过最大长度，移除最旧的摘要
        if (totalLength > maxSummaryLength) {
            while (totalLength > maxSummaryLength && !summaryHistory.isEmpty()) {
                SummaryRecord oldest = summaryHistory.remove(0);
                totalLength -= oldest.getSummary().length();
            }
        }
    }
}

