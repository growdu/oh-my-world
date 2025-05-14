# front/Dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY . .
RUN yarn install && yarn run build-only

# 使用 nginx 作为最终镜像
FROM nginx:latest
COPY --from=build /app/dist /usr/share/nginx/html
