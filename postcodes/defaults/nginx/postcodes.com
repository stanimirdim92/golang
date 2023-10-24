server {
  listen 80;
  listen [::]:80;
  server_name postcodes.local;
  root /var/www/html/postcodes;

location / {
    proxy_pass http://postcodes.local:8089;
}
access_log off;
error_log /var/log/nginx/postcodes.local.log error;
}


