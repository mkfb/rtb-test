# rtb-test

To start env

```make up```

To stop env

```make down```

To send request

```curl -X POST -H "Content-Type: application/json" -H "X-DevOps-Test: true" -d @bid_request.json http://localhost/bid```

Monitoring

```http://localhost/nginx_status```

To configure log rotate 
add to /etc/logrotate.d/nginx

```
  ~/rtb/logs/*.log {
  size 1M
  missingok
  rotate 7
  dateext
  compress
  delaycompress
  notifempty
  sharedscripts
  postrotate
    cd rtb && sudo docker compose kill -s USR1 nginx
  endscript
}
```
