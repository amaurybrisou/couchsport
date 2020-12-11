FROM node:14-alpine AS setup
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package*.json /usr/src/app/
ENV PATH /usr/src/app/node_modules/.bin:$PATH
RUN npm install --loglevel verbose
RUN npm install -g @vue/cli

FROM setup as run-dev
CMD npm run serve

FROM setup as build
COPY . .
RUN npm run build

FROM nginx:1.17.9-alpine as release
WORKDIR /usr/share/nginx/html
COPY --from=build /usr/src/app/dist/ .
RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/conf.d

ENTRYPOINT ["nginx", "-g", "daemon off;"]