FROM node:22.17

WORKDIR /home/node/app

COPY . /home/node/app

ARG GRAPHQL_CODEGEN_SCHEMA_URL
ENV GRAPHQL_CODEGEN_SCHEMA_URL=${GRAPHQL_CODEGEN_SCHEMA_URL}

RUN npm install
RUN npm run codegen
RUN npm run build

EXPOSE 3000

CMD ["npm", "run", "preview"]