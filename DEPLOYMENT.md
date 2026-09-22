# SpotSync Backend Deployment Guide

This guide details the prerequisites, environment variables, build instructions, and storage considerations for deploying the **SpotSync Go Backend**.

---

## 1. Environment Variables

The backend requires the following environment variables configured in production:

| Variable | Required | Description | Example / Default |
| :--- | :--- | :--- | :--- |
| `PORT` | Optional | Port on which the Echo HTTP server listens | `8080` (or dynamically assigned by PaaS) |
| `DSN` | **Required** | PostgreSQL connection URI | `postgresql://user:password@host:5432/spotsync?sslmode=require` |
| `JWT_SECRET` | **Required** | Secret key for signing and validating JWT tokens | `your_strong_random_jwt_secret_key` |
| `ALLOWED_ORIGIN` | **Required in Prod** | Allowed CORS origins (single URL or comma-separated) | `https://spotsync-client.vercel.app,http://localhost:3000` *(Default: `http://localhost:3000`)* |

---

## 2. Build & Start Commands

### Standard Linux / Docker / Cloud PaaS (Render, Railway, Fly.io)

- **Build Command**:
  ```bash
  go build -o server cmd/main.go
  ```

- **Start Command**:
  ```bash
  ./server
  ```

### Windows Local / VM

- **Build Command**:
  ```powershell
  go build -o server.exe cmd/main.go
  ```

- **Start Command**:
  ```powershell
  .\server.exe
  ```

---

## 3. Storage & Persistence Note (Important)

> [!WARNING]
> **Local Disk Storage for Uploads:**
> Uploaded files (such as driver profile avatars) are stored directly on the local filesystem under the `./uploads` directory and served statically via `/uploads/*`.
>
> On containerized or serverless hosting platforms with ephemeral filesystems (e.g., Render standard instances, Heroku dynos, Fly.io ephemeral containers), **any uploaded images will be lost on container restart or redeployment** unless:
> 1. A persistent storage volume is attached and mounted to the `./uploads` directory, OR
> 2. The upload domain service is updated in the future to stream files directly to cloud object storage (e.g., AWS S3, Cloudflare R2, or Cloudinary).
