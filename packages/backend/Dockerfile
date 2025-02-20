FROM node:18-slim

WORKDIR /app

ENV PATH=/app/node_modules/.bin:$PATH

COPY ./packages/backend/package.json ./yarn.lock ./

RUN yarn install --frozen-lockfile

COPY ./packages/backend .

EXPOSE 3000

CMD ["yarn", "dev"]
