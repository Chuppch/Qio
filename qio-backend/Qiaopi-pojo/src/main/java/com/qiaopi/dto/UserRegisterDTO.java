package com.qiaopi.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.*;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Schema(name = "UserRegisterDTO", description = "用户注册对象")
public class UserRegisterDTO {

    /**
     * 用户名
     */
    //必填项
    @Schema(description = "用户名" , required = true,example = "admin")
    private String username;

    /**
     * 用户密码
     */
    @Schema(description = "用户密码" , required = true , example = "123456")
    private String password;

    /**
     * 验证码
     */
    @Schema(description = "验证码", required = true, example = "1234")
    private String code;
    /**
     * 再次确认密码
     */
    @Schema(description = "再次确认密码" , required = true,example = "123456")
    private String confirmPassword;


}
