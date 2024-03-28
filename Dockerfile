FROM node:alpine AS tailwind
WORKDIR /app
COPY ./tailwind.config.js /app/tailwind.config.js
COPY ./templates /app/templates
COPY ./static /app/static
RUN npm install -g tailwindcss
RUN tailwindcss -i ./static/style.css -o ./static/output.css --minify


FROM golang:1.22-bullseye AS builder
WORKDIR /app

RUN apt-get update
RUN apt-get install -y libvips-dev

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate
RUN go mod download

COPY --from=tailwind /app/static /app/static
RUN go build -o server .

COPY awsconfig /root/.aws/config

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./server"]
