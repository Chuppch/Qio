package com.qiaopi.vo;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Schema(description = "字体颜色")
public class FontColorVO {
    @Schema(description = "id")
    private Long id;

    @Schema(description = "颜色名称")
    private String description;

    @Schema(description = "预览图片")
    private String previewImage;


}
