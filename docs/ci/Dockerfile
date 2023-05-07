FROM nginx:1.23.4

COPY ./ci/nginx.conf /etc/nginx/nginx.conf
COPY ./ci/server.template /etc/nginx/templates/server.conf.template
# COPY ./ci/server.conf /etc/nginx/conf.d/server.conf
# COPY ./ci/ga.conf.template /etc/nginx/templates/ga.conf.template
COPY dist /usr/share/nginx/html
