FROM node:20-slim

WORKDIR /app

ENV PATH="$PATH:/app/node_modules/.bin"

COPY package.json yarn.lock .

RUN yarn install --frozen-lockfile

COPY lerna.json .

RUN yarn add lerna@8.2.1 -W

COPY . .

ENTRYPOINT [ "lerna" ]

CMD [ "run", "dev" ]
