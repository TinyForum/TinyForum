# Quick Start

This guide covers Docker-based deployment.

---

## Prerequisites

- Docker Engine 24+ with Docker Compose
- Git

---

## Installation

```bash
git clone https://github.com/caoyang2002/TinyForum.git
cd TinyForum
cp .env.example .env
```

---

## Configuration

Edit `backend/config/private.yml` to set:

- Database credentials
- JWT secret (at least 32 characters)
- Email configuration (optional)
- Admin account credentials

---

## Start

```bash
docker compose up -d
```

Access the forum at `http://localhost:8080`.

---

## Stop

```bash
docker compose down
```

---

## Logs

```bash
docker compose logs -f
```

---

For more detailed guides (including Podman and local development), see the [Chinese Quick Start](/zh-CN/guide/quickstart).
