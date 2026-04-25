# 安装 docker
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```
## 安装 docker-compose

```bash
sudo apt install -y docker-compose-plugin
```


sudo systemctl daemon-reload
sudo systemctl restart docker

## 构建并启动

```bash
docker compose up -d --build
```

