FROM node:18-slim

WORKDIR /app

ENV PATH=/app/node_modules/.bin:$PATH

COPY ./package.json ./yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . .

EXPOSE 5173

CMD ["yarn", "dev", "--host"]
