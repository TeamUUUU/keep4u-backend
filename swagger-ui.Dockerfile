FROM swaggerapi/swagger-ui:latest

RUN apk add --update nodejs
RUN npm install -g swagger-cli

COPY /api /foo
WORKDIR /foo

RUN swagger-cli validate api.yml && \
 swagger-cli bundle -t yaml -o api.bundle.yaml api.yml

ENV SWAGGER_JSON "/foo/api.bundle.yaml"

EXPOSE 8080
