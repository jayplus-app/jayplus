FROM node:lts as dependencies
WORKDIR /frontend
COPY package.json ./
RUN npm install

FROM node:lts as builder
WORKDIR /frontend
COPY . .
COPY --from=dependencies /frontend/node_modules ./node_modules
RUN npm run build

FROM node:lts as runner
WORKDIR /frontend
ENV NODE_ENV production
COPY --from=builder /frontend/.next ./.next
COPY --from=builder /frontend/node_modules ./node_modules
COPY --from=builder /frontend/package.json ./package.json

EXPOSE 3000
CMD [ "npm", "start" ]