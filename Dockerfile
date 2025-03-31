# # Этап 1: Сборка
# FROM golang:1.23 AS builder  
# WORKDIR /app

# COPY go.mod .
# COPY go.sum .       
# RUN go mod download  

# COPY . .                    

# # Сборка проекта
# RUN go build -o /app/todo-go-app ./cmd/app/main.go

# # Этап 2: Минимальный рантайм (final)
# FROM alpine:latest 
# WORKDIR /app     
# COPY --from=builder /app/todo-go-app .  
# CMD ["./todo-go-app"]


# тяжелый образ
FROM golang:1.23-alpine 

RUN go version 
ENV GOPATH=/
COPY ./ ./ 

RUN go mod download
RUN go build -o todo-go-app ./cmd/app/main.go

CMD ["./todo-go-app"]
