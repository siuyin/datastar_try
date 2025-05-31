# Datastart experiments

## Nginx confing 
URL: https://app.beyondbroadcast.com/dstar/

```
location /dstar/ {
  proxy_pass http://localhost:8082;
  proxy_buffering off;
}
```
