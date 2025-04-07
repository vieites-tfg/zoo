FROM node:18-slim as builder

WORKDIR /app

ENV PATH=/app/node_modules/.bin:$PATH

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn run build

FROM node:18-alpine

WORKDIR /app

COPY --from=builder /app/build ./build

EXPOSE 5173

CMD ["npx", "serve", "build"]

