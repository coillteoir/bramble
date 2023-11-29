FROM node:latest as frontend

WORKDIR /app
COPY ./frontend/ .
RUN npm install
RUN npm run build

FROM node:latest as backend
WORKDIR /api
COPY ./backend/ .


RUN npm install
RUN npm run build
COPY --from=frontend /app/dist /api/public

ENV IN_CLUSTER 1

CMD ["node", "dist/index.js"]
