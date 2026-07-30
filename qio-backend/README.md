<p align="center">
	<!-- <img alt="logo" src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/logo.png"> -->
	<img alt="logo" src="dcos/image/主页面照片.png">
</p>
<h4 align="center">跨越四海，侨缘线牵——侨缘信使，让世界没有距离。</h4>
<p align="center">
	<a href="https://gitee.com/trashwbin/qiaopi"><img src="https://img.shields.io/badge/%E4%BE%A8%E7%BC%98%E4%BF%A1%E4%BD%BF-github?logo=gitee&label=gitee&labelColor=%23C71D23&color=%23000"></a>
    <a href="https://github.com/trashwbin/Qiaopi_vue"><img src="https://img.shields.io/badge/%E5%89%8D%E7%AB%AF%E5%B7%A5%E7%A8%8B-gitee?logo=github&label=github&color=%23181717"></a>
	<a href="https://github.com/trashwbin/Qiaopi"><img src="https://img.shields.io/badge/Qiaopi-v1.0.1-brightgreen.svg"></a>
	<a href="https://github.com/trashwbin/Qiaopi?tab=MIT-1-ov-file"><img src="https://img.shields.io/github/license/mashape/apistatus.svg"></a>
</p>



## 项目介绍

《**侨缘信使**》是一个旨在宣传和传承侨批文化的互动网站。文化内容的数字化呈现、互动性与教育性的结合、情感共鸣的构建以及用户参与的文化共创。我们致力于通过网站，让更多人了解并参与到侨批文化的保护与传承中来，同时为文化带来新的活力。学习侨批文化、体验写侨批、收侨批和漂流瓶等功能，感受慢信文化的魅力，"跨越四海，侨缘线牵——侨缘信使，让世界没有距离。"

**技术栈**：***SpringBoot、MybatisPlus、MySQL、Redis、Graphics2D、JWT、Minio、RabbitMQ、ChatGLM***

## **在线体验**

**[侨缘信使🎉](http://110.41.58.26)**

## 演示图

<table>
    <tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/home.png"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/home-receive.png"/></td>
    </tr>
    <tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/home-introduce.gif"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/home-history.gif"/></td>
    </tr>
    <tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/write-letter.gif"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/nav-ai.png"/></td>
    </tr>
	<tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/send-letter.gif"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/drifting.png"/></td>
    </tr>	 
    <tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/game-explore.gif"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/game-question.gif"/></td>
    </tr>
	<tr>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/shop.gif"/></td>
        <td><img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/marketing.gif"/></td>
    </tr>
</table>



## 快速部署

本项目所需的初始化数据已经打包在`init_qiaopi`文件夹下，部分三方密钥需自行获取，账号密码记得修改为自己的，具体需要修改的**application.yml**如下: 

<img src="https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/init.png"  />

### MySQL数据、Redis数据

已打包，自行导入即可

![](https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/mysql-redis.png)

### 邮箱授权码

本项目使用QQ邮箱做发件功能，可自行替换，[SMTP/IMAP服务](https://wx.mail.qq.com/list/readtemplate?name=app_intro.html#/agreement/authorizationCode)

![](https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/QQ-email.png)

### Minio密钥

本项目需要使用对象存储服务，同时配合[X File Storage](https://x-file-storage.xuyanwu.cn)进行文件上传，推荐使用docker搭建minio，速率也是够的，当然也可以使用其它平台的对象存储服务

![](https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/minio-key.png)

### 智谱ChatGLM密钥

本项目的AI交互功能基于智谱ChatGLM实现，通过预设系统提示词使其更具备专业性，**系统提示词**也已上传，可自行修改，部署需要自行申请API Key，[智谱AI开放平台](https://bigmodel.cn/usercenter/apikeys)。

部分功能使用了其封装的 [Java SDK](https://github.com/MetaGLM/zhipuai-sdk-java-v4)，Maven镜像仓可能找不到该依赖，需自行下载 [Zhipu SDK For GLM Open API](https://mvnrepository.com/artifact/cn.bigmodel.openapi/oapi-java-sdk)

![](https://raw.githubusercontent.com/trashwbin/Qiaopi/refs/heads/master/init_qiaopi/images/zhipu-keys.png)

### RabbitMQ、JWT密钥、AES密钥

需要搭建RabbitMQ，服务于AI的互动消息，推荐docker搭建

自定义JWT密钥即可，32字符以上，用于用户登录后的令牌生成

自定义AES密钥即可，32字符，用于答题功能的题目加密