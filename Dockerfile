#
# Base
#
FROM node:20 AS base

WORKDIR /app

# Copy root package files and install Lerna
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

RUN lerna run --scope backend build 

RUN ncc build ./packages/backend/dist/index.js -o compiled-backend.js

#
# Frontend build stage
#
FROM base AS frontend-build

WORKDIR /app

RUN lerna run --scope frontend build 

RUN ncc build ./packages/frontend/dist/main.js -o compiled-frontend.js

#
# Backend
#
FROM node:20-slim AS backend

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
FROM node:20-slim AS frontend

WORKDIR /app

COPY --from=frontend-build /app/compiled-frontend.js .

COPY --from=frontend-build /app/packages/frontend/package.json .

COPY --from=frontend-build /app/packages/frontend/dist ./dist

RUN yarn install --production

ENV NODE_ENV=production

EXPOSE 5173

CMD ["node", "compiled-frontend.js"]
