package com.qiaopi.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.qiaopi.entity.Paper;
import com.qiaopi.entity.Signet;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface PaperMapper extends BaseMapper<Paper> {
}
