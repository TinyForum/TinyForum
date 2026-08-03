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

访问：

```bash
http://localhost:8080
```

## 停止
```bash
sudo docker compose down
```

# 重启
```bash
docker compose up --build -d
```

## 查看日志

```bash
sudo docker logs tiny-forum-nginx --tail 20
sudo docker logs tiny-forum-frontend --tail 50
sudo docker logs tiny-forum-backend --tail 50
```

