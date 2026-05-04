# Ubuntu 安装方法

```bash
# 删除可能残留的失败容器和镜像（可选）
podman-compose down -v --rmi all

# 重新构建（backend 会重新拉取正确的 go 镜像）
podman-compose build

# 后台启动所有服务
podman-compose up -d
```


### 可选：管理多个 Go 版本

如果需要在同一台机器上安装和管理多个 Go 版本，可以试试 `gvm` (Go Version Manager)。

1.  首先，确保系统安装了必要的依赖：
    ```bash
    sudo apt update
    sudo apt install -y curl git mercurial make binutils bison gcc build-essential
    ```
2.  使用以下命令安装 `gvm`：
    ```bash
    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
    source ~/.gvm/scripts/gvm
    ```
3.  安装并切换 Go 版本：
    ```bash
    gvm install go1.22.4 -B  # -B 表示二进制安装，加快速度
    gvm use go1.22.4 --default
    ```


## podman

```bash
sudo apt install podman
```


## podman-compose

```bash
sudo apt install podman-compose
```


```bash
# 1. 停止并彻底删除当前 compose 项目相关的所有资源
podman-compose down -v

# 2. 手动删除可能残留的容器（如果上面没删干净）
podman stop tiny-forum-postgres tiny-forum-redis tiny-forum-backend tiny-forum-frontend tiny-forum-nginx 2>/dev/null
podman rm tiny-forum-postgres tiny-forum-redis tiny-forum-backend tiny-forum-frontend tiny-forum-nginx 2>/dev/null

# 3. 清理所有未使用的容器、镜像、卷（可选但推荐）
podman system prune -a -f --volumes

# 4. 重新构建并启动
podman-compose build --no-cache
podman-compose up -d
```