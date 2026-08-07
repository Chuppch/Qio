package com.chuppch.config;

import org.springframework.ai.openai.OpenAiEmbeddingModel;
import org.springframework.ai.openai.OpenAiEmbeddingOptions;
import org.springframework.ai.openai.api.OpenAiApi;
import org.springframework.ai.transformer.splitter.TokenTextSplitter;
import org.springframework.ai.vectorstore.pgvector.PgVectorStore;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.ai.document.MetadataMode;

import org.springframework.jdbc.core.JdbcTemplate;

/**
 * @author chuppch
 * @description
 * @create 2025/12/16
 */
@Configuration
public class AiAgentConfig {

    // 初始化向量存储
    @Bean("vectorStore")
    @Primary // 优先使用这个bean，避免覆盖问题
    public PgVectorStore pgVectorStore(@Value("${spring.ai.openai.base-url}") String baseUrl,
                                       @Value("${spring.ai.openai.api-key}") String apiKey,
                                       @Value("${spring.ai.openai.embedding.options.model:text-embedding-v4}") String embeddingModelName,
                                       @Qualifier("pgVectorJdbcTemplate") JdbcTemplate jdbcTemplate) {

        OpenAiApi openAiApi = OpenAiApi.builder()
                .baseUrl(baseUrl)
                .apiKey(apiKey)
                .build();

        OpenAiEmbeddingOptions embeddingOptions = OpenAiEmbeddingOptions.builder()
                .model(embeddingModelName)  // 指定模型名称：text-embedding-v4
                .dimensions(1536)  // 显式指定维度为 1536，与数据库表 vector_store_openai 的 VECTOR(1536) 定义一致
                .build();

        OpenAiEmbeddingModel embeddingModel =
            new OpenAiEmbeddingModel(
                    openAiApi,
                    MetadataMode.EMBED,
                    embeddingOptions
            ); 
        
        
        return PgVectorStore.builder(jdbcTemplate, embeddingModel)
                .vectorTableName("vector_store_openai")
                .build();
    }

    // 初始化文本分割器
    @Bean("tokenTextSplitter")
    public TokenTextSplitter tokenTextSplitter() {
        // todo 后续需要补充更详细的参数
        return new TokenTextSplitter();
    }

}