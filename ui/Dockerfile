FROM node:current-alpine as frontend
WORKDIR /app
COPY ./frontend/ .

RUN npm install
RUN npm run build

FROM node:current-alpine as backend
WORKDIR /api
COPY ./backend/ .

RUN npm install && npm run build
COPY --from=frontend /app/dist /api/public

ENTRYPOINT ["node", "dist/index.js"]
