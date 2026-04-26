# mac 运行



## 1. 确保已安装 Podman
```bash
brew install podman  # 如果还没安装
```

## 2. 初始化并启动 Podman machine
```bash
podman machine init      # 创建虚拟机（默认分配 2GB 内存，可根据需要调整）

```bash
podman machine rm   # 删除现有 machine（会丢失已有的容器和镜像，谨慎操作）
podman machine init --memory 4096 --cpus 2   # 分配 4GB 内存
podman machine start
```
## 启动

```bash
podman machine start     # 启动虚拟机
```

- 首次启动会自动配置 socket 连接。完成后，你可以用 `podman info` 验证是否正常。

根据提示输入，我的输出是：

```bash
export DOCKER_HOST='unix:///var/folders/22/jpg085553gv6xck8gdsnlhcr0000gn/T/podman/podman-machine-default-api.sock'

# 也可以加到 ~/.zshrc 中，然后重新启动：
podman machine stop
podman machine start
```
## 3. 验证连接
```bash
podman version           # 应显示 Client 和 Server 版本
podman ps                # 应正常返回（空列表）
```

## 4. 安装 podman-compose（如果还没有）
```bash
brew install podman-compose   # 或使用 brew install podman-compose
```

## 5. 重新运行 podman-compose
```bash
cd ~/Documents/Github/TinyForum
podman-compose up -d
```

此时应该可以正常启动。

## 查看日志
```bash
podman-compose logs -f
podman logs tiny-forum-backend                            
podman logs --tail 50 tiny-forum-backend

```

## 停止
```bash
podman-compose down
podman machine stop
```

## 补充说明

- **资源占用**：Podman machine 默认会分配 2GB 内存和 2 个 CPU，如果你的 Mac 资源紧张，可以在 `init` 时指定 `--cpus` 和 `--memory`。
- **与 Docker Desktop 的区别**：Podman machine 在 macOS 上本质上也是一个虚拟机（类似 Docker Desktop 的 HyperKit/VM），但无守护进程、无 root 权限要求。
- **推荐方案**：如果你已经安装了 **Docker Desktop for Mac**，可以直接使用 `docker-compose up -d`（因为你的原配置本身就是 Docker Compose 格式），会更简单。Docker Desktop 在 macOS 上也是通过虚拟机运行，但用户体验更成熟。如果你希望使用 Podman 作为后端，则按上述步骤操作。

## 常见问题

**Q: 可以不使用 `podman machine` 直接在 macOS 运行 Podman 吗？**  
A: 不可以，macOS 没有 Linux 内核，必须依赖虚拟机。这也是为什么 Docker Desktop 同样需要虚拟化后端。

如果你只是想快速运行 TinyForum 项目，**建议直接使用 Docker Desktop**（安装后执行 `docker compose up -d` 即可），省去 Podman machine 的额外配置步骤。