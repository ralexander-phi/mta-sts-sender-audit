FROM nginx:1.27-alpine
WORKDIR /static/
COPY /mta-sts.txt /static/.well-known/mta-sts.txt
COPY /nginx.conf /etc/nginx/nginx.conf

