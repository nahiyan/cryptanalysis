TARGET = nejati_preimage_encoder
LEGACY_TARGET = nejati_preimage_encoder_cpp
LEGACY_SRCS = formula.cpp util.cpp hash.cpp sha1.cpp sha256.cpp md4.cpp
LEGACY_OBJS = $(LEGACY_SRCS:.cpp=.o)
CXXFLAGS = -std=c++11 -O2

.PHONY: all rust legacy clean

all: rust

rust:
	cargo build --release
	cp target/release/$(TARGET) $(TARGET)

legacy: $(LEGACY_OBJS) main.cpp
	$(CXX) $(LEGACY_OBJS) $(CXXFLAGS) main.cpp -o $(LEGACY_TARGET)

%.o: %.cpp
	$(CXX) -c $< -o $@ $(CXXFLAGS)

clean:
	cargo clean
	rm -f $(LEGACY_OBJS) $(TARGET) $(LEGACY_TARGET)
