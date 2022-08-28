FROM node:18-alpine3.15

WORKDIR /app

# install extra tools
RUN apk update && apk add \
	ruby-dev

RUN gem install \
	bundler \
	prettier_print \
	syntax_tree \
	syntax_tree-haml

# install direct dependencies
COPY package.json yarn.lock /app/
RUN yarn install

# our server code
COPY prettier.js .

# runtime settings
EXPOSE 3000
CMD ["node", "prettier.js"]
