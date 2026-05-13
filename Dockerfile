FROM golang:1.26.3-alpine AS build_deps

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM gcr.io/distroless/static:nonroot

COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["/webhook"]
