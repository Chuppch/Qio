package com.chuppch.config;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;

/**
 *  TODO 后续需要完善
 * Spring AI 专用Jackson全局配置类
 * 彻底解决：code10(换行符)、code32(空格符)等所有JSON解析异常
 * 优先级最高，Spring AI底层的ModelOptionsUtils会自动使用该配置
 * 
 * @author chuppch
 * @create 2026/01/04
 */
@Configuration
public class SpringAiJacksonGlobalConfig {

    @Bean
    @Primary // 必须加这个注解，标记为全局默认的ObjectMapper，优先级最高
    @SuppressWarnings("deprecation") // 抑制deprecated警告，这些配置虽然标记为deprecated但仍然是解决JSON解析问题的关键
    public ObjectMapper objectMapper() {
        ObjectMapper objectMapper = new ObjectMapper();
        
        // ========== 注册JSR310模块：支持Java 8时间类型（LocalDateTime等）序列化/反序列化 ==========
        objectMapper.registerModule(new JavaTimeModule());
        // 禁用时间戳格式，将LocalDateTime序列化为ISO-8601字符串格式（如"2024-01-01T12:00:00"），前端可正常解析
        objectMapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
        
        // ========== 核心双开关：一次性解决本次+上次的所有JSON解析错误 ==========
        // 开关1：允许JSON字符串中包含【未转义的原生控制字符】(code10换行符、code13回车符等) → 解决上次的报错
        // 注意：虽然标记为deprecated，但这是解决code10错误的唯一有效方法
        objectMapper.configure(JsonParser.Feature.ALLOW_UNQUOTED_CONTROL_CHARS, true);
        // 开关2：允许JSON字符串中包含【所有合法/非法的转义符】(比如\空格、\xxx等) → 解决本次的 code32 核心报错
        // 注意：虽然标记为deprecated，但这是解决code32错误的唯一有效方法
        objectMapper.configure(JsonParser.Feature.ALLOW_BACKSLASH_ESCAPING_ANY_CHARACTER, true);

        // ========== 推荐附加容错配置（Spring AI必备，锦上添花，避免后续其他报错） ==========
        // 1. 解析JSON时，遇到未知字段（大模型返回的JSON字段可能多变），不抛异常
        objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        // 2. 解析JSON时，空字符串转成null，不抛异常
        objectMapper.configure(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT, true);
        // 3. 允许JSON数字以0开头（比如00123），大模型可能返回这种格式
        // 注意：虽然标记为deprecated，但这是处理大模型不标准JSON格式的必要配置
        objectMapper.configure(JsonParser.Feature.ALLOW_NUMERIC_LEADING_ZEROS, true);
        
        return objectMapper;
    }
}

