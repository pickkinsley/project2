# Infrastructure Setup

## Server Details

- **Provider:** AWS Lightsail
- **IP:** 44.219.227.2
- **Domain:** pickkinsley.online
- **SSL:** Let's Encrypt (auto-renews)
- **OS:** Ubuntu 24.04 LTS

---

## Database

- **MySQL:** 8.0
- **Database:** `packsmart`
- **User:** `packsmart_user`
- **Password:** see your password manager
- **Connection string:** `packsmart_user:[PASSWORD]@tcp(localhost)/packsmart`

Store the actual password in a `.env` file (not committed to git) or your password manager.

---

## Directory Structure

| Path | Contents |
|---|---|
| `/var/www/packsmart` | Frontend static files |
| `/home/ubuntu/packsmart-backend` | Go backend binary |
| `/etc/nginx/sites-enabled/packsmart` | nginx config |
| `/var/archives/` | Old Jones County XC project |

---

## Services

| Service | Details |
|---|---|
| nginx | Port 443 (HTTPS) + Port 80 → 443 redirect |
| Backend API | Port 8080, proxied through nginx at `/api/` |

---

## Deployment Process

### Frontend

```bash
# Build locally
npm run build

# Deploy to server
scp -i ~/.ssh/lightsail.pem -r dist/* ubuntu@44.219.227.2:/var/www/packsmart/
```

### Backend

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server main.go

# Deploy binary
scp -i ~/.ssh/lightsail.pem server ubuntu@44.219.227.2:/home/ubuntu/packsmart-backend/

# Restart (SSH in, then)
pkill server
nohup ./server &
```

---

## Verification

- **Frontend:** https://pickkinsley.online
- **API Health:** https://pickkinsley.online/api/health
- **Status:** ✅ Verified April 6, 2026
