FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM public.ecr.aws/lambda/go:1

ENV TZ=Asia/Tokyo

COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

CMD [ "main" ]