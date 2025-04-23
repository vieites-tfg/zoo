#
# Base
#
FROM node:20 AS base

WORKDIR /app

COPY package.json lerna.json yarn.lock* ./

COPY packages packages/

RUN yarn install

RUN yarn global add lerna@8.2.1

RUN yarn global add @vercel/ncc

#
# Backend build
#
FROM base AS backend-build

WORKDIR /app

RUN lerna run --scope @vieites-tfg/zoo-backend build 

RUN ncc build ./packages/backend/dist/index.js -o compiled-backend.js

#
# Frontend build stage
#
FROM base AS frontend-build

WORKDIR /app

RUN lerna run --scope @vieites-tfg/zoo-frontend build 

#
# Backend
#
FROM node:20-alpine AS backend

WORKDIR /app

COPY --from=backend-build /app/compiled-backend.js .

COPY --from=backend-build /app/packages/backend/package.json .

ENV NODE_ENV=production

RUN yarn install --production

EXPOSE 3000

CMD ["node", "index.js"]

#
# Frontend
#
FROM nginx:alpine AS frontend

WORKDIR /usr/share/nginx/html

COPY --from=frontend-build /app/packages/frontend/dist .

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
