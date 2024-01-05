FROM rust:latest as builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM alpine

# GLIBC and LIBGCC needed for executable to run. 
# If container gives "bramble-git-proxy command not found, usually it's because of a linker issue"
RUN apk upgrade --no-cache && \
    apk add --no-cache libgcc gcompat

COPY --from=builder /app/target/release/bramble-git-proxy /bramble-git-proxy

ENTRYPOINT /bramble-git-proxy
