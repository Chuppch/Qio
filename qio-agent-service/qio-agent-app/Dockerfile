# 基础镜像，可以先执行 docker pull openjdk:17-jdk-slim
FROM openjdk:17-jdk-slim

LABEL org.opencontainers.image.authors="Qio Team"

# 配置
ENV PARAMS=""

# 时区
ENV TZ=PRC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 添加应用
ADD target/qio-agent-app.jar /qio-agent-app.jar

ENTRYPOINT ["sh","-c","java -jar $JAVA_OPTS /qio-agent-app.jar $PARAMS"]
