### minio-webhook

minio-webhook listens for webhook events from MinIO server and logs the event to a log file.

Usage:
```
minio-webhook <logfile>
```

Environment variables:
* MINIO_WEBHOOK_AUTH_TOKEN: Authorization token to be used by minio server for sending events
* MINIO_WEBHOOK_PORT: Listening port (Default 8080)

The minio-webhook service can be setup as a systemd service using the provided minio-webhook.service file

Logs can be rotated using the standard logrotate tool. You can provide the postrotate command such that
minio-webhook writes to a new log file after log rotation.
```
postrotate
	systemctl reload minio-webhook
endscript
```