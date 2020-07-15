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

To send access/audit logs from MinIO server, please configure MinIO using the command:
```
mc admin config set myminio audit_webhook endpoint=http://webhookendpoint:8080 auth_token=webhooksecret
```

To send MinIO's error logs to minio-webhook, please configure MinIO using the command:
```
mc admin config set myminio logger_webhook endpoint=http://webhookendpoint:8081 auth_token=webhooksecret
```

Note: audit_webhook and logger_webhook should *not* be configured to send events to the same minio-webhook instance.

Logs can be rotated using the standard logrotate tool. You can provide the postrotate command such that
minio-webhook writes to a new log file after log rotation.
```
postrotate
	systemctl reload minio-webhook
endscript
```