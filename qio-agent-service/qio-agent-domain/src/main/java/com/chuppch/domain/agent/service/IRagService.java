package com.chuppch.domain.agent.service;

import org.springframework.web.multipart.MultipartFile;

import java.util.List;

/**
 * @author chuppch
 * @description
 * @create 2026/1/6
 */
public interface IRagService {

    void storeRagFile(String name, String tag, List<MultipartFile> files);

}
