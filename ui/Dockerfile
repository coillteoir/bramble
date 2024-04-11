FROM node:current-alpine as frontend
WORKDIR /app

COPY frontend/*.json .
RUN npm install  

COPY frontend/index.html index.html
COPY frontend/*.js .
COPY frontend/public public
COPY frontend/src src
COPY frontend/vite.config.ts vite.config.ts
RUN npm run build

FROM node:current-alpine as backend
WORKDIR /api

COPY ./backend .


COPY backend/*.json .
COPY backend/src src/

RUN npm install 
RUN npm run build
COPY --from=frontend /app/dist /api/public

ENTRYPOINT ["node", "dist/index.js"]
