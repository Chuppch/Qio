# 注意修改 src/config/api.ts 为你的ip，无论是云服务器IP还是本地IP，都要修改为你的
docker build --load -t qio-agent-console:1.0 -f ./Dockerfile .
