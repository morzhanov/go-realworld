# pull official base image
FROM node:16-alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH
ENV API_URI=http://127.0.0.1/3005

# install app dependencies
COPY ../web/package.json ./
COPY ../web/package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent

# add app
COPY ../web ./

# start app
CMD ["npm", "start"]
