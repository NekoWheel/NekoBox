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

RUN echo "ErrorDocument 404 /index.html" >> /usr/local/apache2/conf/httpd.conf

EXPOSE 80