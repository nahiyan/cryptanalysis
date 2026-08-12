TARGET = nejati_collision_encoder
LEGACY_TARGET = nejati_collision_encoder_cpp
LEGACY_SRCS = diff-parser.cpp formula.cpp main.cpp sha256x.cpp util.cpp xformula.cpp

.PHONY: all rust legacy clean

all: rust

rust:
	cargo build --release
	cp target/release/$(TARGET) $(TARGET)

legacy:
	$(CXX) $(CXXFLAGS) $(LEGACY_SRCS) -o $(LEGACY_TARGET)

clean:
	cargo clean
	rm -f $(TARGET) $(LEGACY_TARGET) encoder
