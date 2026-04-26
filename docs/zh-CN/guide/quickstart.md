# 快速开始

> 这里使用容器启动改项目
>
> 已通过 MacOS 和 Ubuntu 的 Docker / Podman 测试
>
> 这里仅简单介绍本地通过 Docker 环境启动

## Docker 启动（Ubuntu）

!> 前置条件

本项目使用 Docker Compose 启动，请先安装 Docker（Docker 已内置 Compose 命令），参考 [官方文档](https://docs.docker.com/engine/install/)

### 1. 安装 Docker

#### 1.1 如果是测试环境可以使用快捷脚本：

```bash
curl -fsSL https://get.docker.com | sudo sh
```

#### 1.2 如果在生产环境中，推荐以下方式：

（可选）先移除系统可能自带的旧版本，

```bash
sudo apt update
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
```

配置 Docker 官方 APT 仓库

```bash
# 1. 安装依赖包，用于通过 HTTPS 添加仓库
sudo apt install -y ca-certificates curl
# 2. 创建密钥存储目录并下载 Docker 官方 GPG 密钥
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
# 3. 将 Docker 仓库添加到 APT 源列表中
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
```

安装 Docker 引擎

```bash
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

测试安装是否成功

```bash
sudo docker run hello-world
```

### 2. 下载本项目

```bash
git clone https://github.com/TinyForum/TinyForum.git
```

### 3. 设置环境变量文件

```bash
cp .env.example .env
```

### 4. 修改后端配置

由于 Docker 使用 Nginx 作为反向代理，所以需要修改后端配置文件。

打开 `backend/config/basic.yaml.compose`，修改以下内容，改为你自己的前端地址：

```yaml
allow_origins:
  - http://localhost:3000 # 这是开发环境的地址，如果你是生产环境，请删除这一行
  - http://127.0.0.1:3000 # 这是开发环境的地址，如果你是生产环境，请删除这一行
  - http://192.168.5.243:8080 # 修改这个，这是 Docker 内的前端地址
```

### 5. 启动

```bash
sudo docker compose up -d
``` 

### 查看日志
```bash
sudo docker compose logs -f
```
## 敏感内容检测

本项目配置了敏感内容检测（DFA），但由于 DFA 误判严重，因而添加了 LLM 进行敏感内容复核。
> 注意：
> 
> Docker 容器没有内置 Ollama，需要在宿主机安装 Ollam 和模型，然后修改配置 `backend/config/private.yaml.compose` 中的 `ollama.base_url` 改为你自己的 Ollama 地址。例如 192.168.5.243:11434

可以使用 ollama 启动 LLM 后端，安装方法如下：

```bash
sudo apt update && sudo apt upgrade
curl -fsSL https://ollama.com/install.sh | sh
```

然后下载一个轻量的 LLM，比如 Qwen3-0.6B-GGUF

可以在 `backend/llm` 执行脚本下载：
```bash
cd backend/llm
./download.sh
```

>  通过官方脚本安装 Ollama 时，它会自动创建一个 systemd 服务并启动，默认监听 127.0.0.1:11434

然后修改 `backend/config/llm.yaml`，将 `ollama.model_path` 修改为地址。

创建模型

```bash
ollama create qwen3-0.6b -f Modelfile
```


## 测试

打开浏览器，访问 `http://localhost:8080`，即可看到项目。

## 停止

```bash
docker-compose down
```