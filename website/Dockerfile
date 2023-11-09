FROM node:lts as dependencies
WORKDIR /website
COPY package.json ./
RUN npm install

FROM node:lts as builder
WORKDIR /website
COPY . .
COPY --from=dependencies /website/node_modules ./node_modules
RUN npm run build

FROM node:lts as runner
WORKDIR /website
ENV NODE_ENV production
COPY --from=builder /website/.next ./.next
COPY --from=builder /website/node_modules ./node_modules
COPY --from=builder /website/package.json ./package.json

CMD [ "npm", "start" ]