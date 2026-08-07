#!/bin/bash

# https://cr.console.aliyun.com/cn-hangzhou/instance/credentials

# Ensure the script exits if any command fails
set -e

# Define variables for the registry and image
ALIYUN_REGISTRY="registry.cn-hangzhou.aliyuncs.com"
IMAGE_NAME="qio-agent-console"
IMAGE_TAG="1.0"

# 读取本地配置文件
if [ -f ".local-config" ]; then
  source .local-config
else
  echo ".local-config 文件不存在，请创建并填写 ALIYUN_USERNAME 和 ALIYUN_PASSWORD"
  exit 1
fi

NAMESPACE="${ALIYUN_NAMESPACE:-qio}"
LOCAL_IMAGE="${IMAGE_NAME}:${IMAGE_TAG}"
REMOTE_IMAGE="${ALIYUN_REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${IMAGE_TAG}"

# Login to Aliyun Docker Registry
echo "Logging into Aliyun Docker Registry..."
docker login --username="${ALIYUN_USERNAME}" --password="${ALIYUN_PASSWORD}" $ALIYUN_REGISTRY

# Tag the Docker image
echo "Tagging the Docker image..."
docker tag "${LOCAL_IMAGE}" "${REMOTE_IMAGE}"

# Push the Docker image to Aliyun
echo "Pushing the Docker image to Aliyun..."
docker push "${REMOTE_IMAGE}"

echo "Docker image pushed successfully! "

echo "检出地址：docker pull ${REMOTE_IMAGE}"
echo "标签设置：docker tag ${REMOTE_IMAGE} ${LOCAL_IMAGE}"

# Logout from Aliyun Docker Registry
echo "Logging out from Aliyun Docker Registry..."
docker logout $ALIYUN_REGISTRY
