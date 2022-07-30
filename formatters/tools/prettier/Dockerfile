FROM node:18-alpine3.15

WORKDIR /app

COPY package.json yarn.lock /app/
RUN yarn install

# setup server
COPY prettier.js .

EXPOSE 3000
CMD ["node", "prettier.js"]
