# CRM Monitoring Thresholds

ACC-017 production monitoring must cover the runtime host, Docker Compose services,
PostgreSQL, Nginx ingress, backup jobs, and TLS certificate renewal.

| Area | Warning | Critical | Evidence |
|---|---:|---:|---|
| Root disk free | below 20% free | below 10% free | `df -h /`, container volume usage |
| Memory available | below 20% available for 15 minutes | below 10% available or OOM event | `free -m`, kernel logs |
| CPU load | above 80% for 15 minutes | above 95% for 5 minutes | host metrics |
| Docker service health | any CRM container unhealthy | gateway, postgres, or identity-authz unhealthy | `docker compose -f docker-compose.prod.yml ps` |
| Nginx ingress | 5xx above normal baseline | sustained 5xx or TLS handshake failures | Nginx access/error logs |
| PostgreSQL | healthcheck failure | unavailable or storage pressure | container health + logs |
| Backup job | missed daily run | failed backup or failed off-server copy | backup logs and restore rehearsal |
| TLS certificate | under 48 hours remaining | under 24 hours remaining or expired | `openssl x509 -checkend` |

Certificate monitoring is mandatory because the approved IP certificate path uses
Let’s Encrypt short-lived certificates. Renewal failure must alert before the
certificate drops below 48 hours of validity.
