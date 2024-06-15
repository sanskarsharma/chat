FROM nginx:1.15.8-alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY public /www/public
# not being used in docker-compose, using CMD and not ENTRYPOINT
# COPY run.sh /

CMD ["nginx", "-g", "daemon off;"]
