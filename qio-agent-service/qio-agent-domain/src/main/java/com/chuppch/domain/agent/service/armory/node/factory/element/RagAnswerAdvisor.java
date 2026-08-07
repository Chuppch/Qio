package com.chuppch.domain.agent.service.armory.node.factory.element;

import com.alibaba.fastjson2.JSON;
import org.springframework.ai.chat.client.ChatClientRequest;
import org.springframework.ai.chat.client.ChatClientResponse;
import org.springframework.ai.chat.client.advisor.api.AdvisorChain;
import org.springframework.ai.chat.client.advisor.api.BaseAdvisor;
import org.springframework.ai.chat.client.advisor.api.CallAdvisorChain;
import org.springframework.ai.chat.client.advisor.api.StreamAdvisorChain;
import org.springframework.ai.chat.messages.AssistantMessage;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.ai.vectorstore.SearchRequest;
import org.springframework.ai.vectorstore.VectorStore;
import org.springframework.ai.vectorstore.filter.Filter;
import org.springframework.ai.vectorstore.filter.FilterExpressionTextParser;
import reactor.core.publisher.Flux;
import org.springframework.util.StringUtils;


import org.springframework.ai.document.Document;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

/**
 * @author chuppch
 * @description
 * @create 2025/12/18
 */
public class RagAnswerAdvisor implements BaseAdvisor {

    // 向量存储
    private final VectorStore vectorStore;
    // 搜索请求
    private final SearchRequest searchRequest;
    // 用户文本提示词
    private final String userTextAdvise;

    public RagAnswerAdvisor(VectorStore vectorStore, SearchRequest searchRequest) {
        this.vectorStore = vectorStore;
        this.searchRequest = searchRequest;
        this.userTextAdvise = "\nContext information is below, surrounded by ---------------------\n\n---------------------\n{question_answer_context}\n---------------------\n\nGiven the context and provided history information and not prior knowledge,\nreply to the user comment. If the answer is not in the context, inform\nthe user that you can't answer the question.\n";
    }

    @Override
    public ChatClientRequest before(ChatClientRequest chatClientRequest, AdvisorChain advisorChain) {
        HashMap<String, Object> context = new HashMap(chatClientRequest.context());

        // 提取用户原始查询文本
        String userText = chatClientRequest.prompt().getUserMessage().getText();
        // 拼接提示词模板
        String advisedUserText = userText + System.lineSeparator() + this.userTextAdvise;

        // 构建搜索请求 - 设置用户输入文本，设置过滤表达式，
        SearchRequest searchRequestToUse = SearchRequest.from(this.searchRequest).query(userText).filterExpression(this.doGetFilterExpression(context)).build();
        // 执行向量相似度搜索
        List<Document> documents = this.vectorStore.similaritySearch(searchRequestToUse);
        // 将检索结果写入上下文
        context.put("qa_retrieved_documents", documents);

        // 合并文档文本 - 将文档列表转为流，提取每个文档的文本内容，通过换行符连接所有文档文本
        String documentContext = documents.stream().map(Document::getText).collect(Collectors.joining(System.lineSeparator()));
        // 放入文档上下文
        Map<String, Object> advisedUserParams = new HashMap(chatClientRequest.context());
        advisedUserParams.put("question_answer_context", documentContext);

        return ChatClientRequest.builder()
                .prompt(Prompt.builder().messages(new UserMessage(advisedUserText), new AssistantMessage(JSON.toJSONString(advisedUserParams))).build())
                .context(advisedUserParams)
                .build();
    }


    @Override
    public ChatClientResponse after(ChatClientResponse chatClientResponse, AdvisorChain advisorChain) {
        // 构建响应 - 从响应中提取元数据，将检索结果写入响应元数据
        ChatResponse.Builder chatResponseBuilder = ChatResponse.builder().from(chatClientResponse.chatResponse());
        chatResponseBuilder.metadata("qa_retrieved_documents", chatClientResponse.context().get("qa_retrieved_documents"));
        ChatResponse chatResponse = chatResponseBuilder.build();

        return ChatClientResponse.builder()
                .chatResponse(chatResponse)
                .context(chatClientResponse.context())
                .build();
    }

    @Override
    public ChatClientResponse adviseCall(ChatClientRequest chatClientRequest, CallAdvisorChain callAdvisorChain) {
        // 直接调用before处理请求，将处理后的请求传给责任链，执行实际的AI调用
        ChatClientResponse chatClientResponse = callAdvisorChain.nextCall(this.before(chatClientRequest, callAdvisorChain));
        return this.after(chatClientResponse, callAdvisorChain);
    }

    @Override
    public Flux<ChatClientResponse> adviseStream(ChatClientRequest chatClientRequest, StreamAdvisorChain streamAdvisorChain) {
        return BaseAdvisor.super.adviseStream(chatClientRequest, streamAdvisorChain);
    }

    @Override
    public String getName() {
        return BaseAdvisor.super.getName();
    }

    @Override
    public int getOrder() {
        return 0;
    }

    // 获取过滤表达式
    private Filter.Expression doGetFilterExpression(HashMap<String, Object> context) {
        // 根据是否存在 qa_filter_expression 进行判断，如果存在则使用动态过滤表达式，否则使用默认过滤表达式
        return context.containsKey("qa_filter_expression") && StringUtils.hasText(context.get("qa_filter_expression").toString())
                ?
                // 动态过滤表达式
                (new FilterExpressionTextParser()).parse(context.get("qa_filter_expression").toString())
                :
                // 默认过滤表达式
                this.searchRequest.getFilterExpression();
    }

}