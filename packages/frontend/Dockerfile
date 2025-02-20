FROM node:18-slim

WORKDIR /app

ENV PATH=/app/node_modules/.bin:$PATH

COPY ./packages/frontend/package.json ./yarn.lock ./

RUN yarn install --frozen-lockfile

COPY ./packages/frontend .

EXPOSE 5173

CMD ["yarn", "dev", "--host"]
