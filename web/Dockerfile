FROM guergeiro/pnpm:22-9-slim AS vite_builder

RUN apt update && apt install -y tzdata git && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

ARG GITHUB_SHA=dev

ENV VITE_COMMIT_SHA=${GITHUB_SHA}

WORKDIR /app

COPY . .

RUN rm -rf node_modules dist

RUN pnpm install

RUN pnpm run build

FROM httpd:alpine3.19

COPY --from=vite_builder /app/dist/ /usr/local/apache2/htdocs/

RUN sed -i \
    -e 's/^#\(LoadModule rewrite_module modules\/mod_rewrite.so\)/\1/' \
    -e '/<Directory "\/usr\/local\/apache2\/htdocs">/,/<\/Directory>/s/AllowOverride None/AllowOverride All/' \
    /usr/local/apache2/conf/httpd.conf
RUN echo 'RewriteEngine On' > /usr/local/apache2/htdocs/.htaccess && \
    echo 'RewriteCond %{REQUEST_FILENAME} !-f' >> /usr/local/apache2/htdocs/.htaccess && \
    echo 'RewriteCond %{REQUEST_FILENAME} !-d' >> /usr/local/apache2/htdocs/.htaccess && \
    echo 'RewriteRule . /index.html [L]' >> /usr/local/apache2/htdocs/.htaccess

EXPOSE 80