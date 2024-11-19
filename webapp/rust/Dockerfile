FROM rust:1.82-bookworm

WORKDIR /home/isucon/webapp/rust

RUN apt-get update && apt-get install --no-install-recommends -y \
  default-mysql-client-core=1.1.0 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY ./Cargo.toml ./Cargo.lock ./
RUN mkdir src && echo 'fn main() {}' > ./src/main.rs && cargo build --release --locked && rm src/main.rs target/release/deps/isuride-*
COPY ./src/ ./src/
RUN cargo build --release --locked --frozen

EXPOSE 8080
CMD ["/home/isucon/webapp/rust/target/release/isuride"]
