FROM node:current-alpine as frontend
WORKDIR /app

COPY frontend/*.json ./
RUN npm install  

COPY frontend/index.html index.html
COPY frontend/*.js ./
COPY frontend/public public
COPY frontend/vite.config.ts vite.config.ts
COPY frontend/src src
RUN npm run build

FROM node:current-alpine as backend
WORKDIR /api

COPY backend/*.json ./
RUN npm install 

COPY backend/src src/
RUN npm run build
COPY --from=frontend /app/dist /api/public

ENTRYPOINT ["node", "dist/index.js"]
