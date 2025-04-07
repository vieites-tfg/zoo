FROM node:20-alpine AS builder

LABEL "org.opencontainers.image.source"="https://github.com/vieites-tfg/zoo"

WORKDIR /app

COPY yarn.lock package.json ./

COPY packages/backend/package.json ./packages/backend/

COPY ./packages/backend/tsconfig.json ./packages/backend/

RUN yarn install --frozen-lockfile

COPY packages/backend/src ./packages/backend/src

WORKDIR /app/packages/backend

RUN yarn build

FROM node:20-alpine AS production

LABEL org.opencontainers.image.source="https://github.com/vieites-tfg/zoo"

LABEL org.opencontainers.image.description="Backend image for zoo app"

WORKDIR /app

ARG PORT=3000

ENV PORT=${PORT}

EXPOSE ${PORT}

COPY yarn.lock package.json ./

COPY packages/backend/package.json ./packages/backend/

RUN yarn install --frozen-lockfile --production && yarn cache clean

COPY --from=builder /app/packages/backend/dist ./dist

CMD ["node", "dist/index.js"]
