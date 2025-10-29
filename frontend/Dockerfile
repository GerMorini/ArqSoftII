FROM node:22-alpine AS builder
WORKDIR /app

ARG VITE_USERS_URL
ARG VITE_ACTIVITIES_URL
ARG VITE_SEARCH_URL
ARG VITE_LOG_LEVEL
ENV VITE_USERS_URL=$VITE_USERS_URL
ENV VITE_ACTIVITIES_URL=$VITE_ACTIVITIES_URL
ENV VITE_SEARCH_URL=$VITE_SEARCH_URL
ENV VITE_LOG_LEVEL=$VITE_LOG_LEVEL

COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build


FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template
EXPOSE 80
