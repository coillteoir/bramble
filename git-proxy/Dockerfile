# Build the manager binary
FROM golang:1.22.1 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY cmd cmd
COPY util util
COPY v1alpha1 v1alpha1

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o bramble-git-proxy main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details

FROM alpine as debug
WORKDIR /
COPY --from=builder /workspace/bramble-git-proxy .
USER 65532:65532

ENTRYPOINT ["/bramble-git-proxy"]


FROM gcr.io/distroless/static:nonroot as release
WORKDIR /
COPY --from=builder /workspace/bramble-git-proxy .
USER 65532:65532

ENTRYPOINT ["/bramble-git-proxy"]
