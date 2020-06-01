all: target/release/libtxauxdecoder.a
	go build -o chainindex ./cmd/chainindex

txauxdecoder: target/release/libtxauxdecoder.a

target/release/libtxauxdecoder.a: adapter/txauxdecoder/src/lib.rs adapter/txauxdecoder/Cargo.toml
	cargo build --release --manifest-path=adapter/txauxdecoder/Cargo.toml

clean:
	rm -rf adapter/txauxdecoder/target