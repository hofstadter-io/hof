FROM nginx:1.21.6

COPY ./ci/nginx.conf /etc/nginx/nginx.conf
COPY ./ci/server.conf /etc/nginx/conf.d/default.conf
COPY dist /usr/share/nginx/html
