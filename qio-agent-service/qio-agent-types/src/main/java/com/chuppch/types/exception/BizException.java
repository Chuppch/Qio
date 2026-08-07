package com.chuppch.types.exception;

import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serial;

/**
 * @author chuppch
 * @description 业务异常
 * @create 2026/1/5
 */
@EqualsAndHashCode(callSuper = true)
@Data
public class BizException extends RuntimeException {

    @Serial
    private static final long serialVersionUID = 5317680961212299217L;

    /** 异常码 */
    private String code;

    /** 异常信息 */
    private String info;

    public BizException(String code) {
        this.code = code;
    }

    public BizException(String code, Throwable cause) {
        this.code = code;
        super.initCause(cause);
    }

    public BizException(String code, String message,Throwable  cause) {
        this.code = code;
        this.info = message;
        super.initCause(cause);
    }

    @Override
    public String toString() {
        return getClass().getName() + "{" +
                "code='" + code + '\'' +
                ", info='" + info + '\'' +
                '}';
    }
}
