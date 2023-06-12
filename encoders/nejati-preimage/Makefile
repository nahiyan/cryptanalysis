SRCS = formula.cpp util.cpp hash.cpp sha1.cpp sha256.cpp md4.cpp
OBJS = $(SRCS:.cpp=.o)
TARGET = nejati_preimage_encoder
CFLAGS = -std=c++11 -O2
MACROS = #-D_DOTMATRIX_ #-D_COMPRESSOR_TEST

CFLAGS += $(MACROS)

.cpp.o: $(SRCS)
	g++ -c $< -o $@ $(CFLAGS)

all: $(OBJS) main.cpp
	g++ $(OBJS) $(CFLAGS) main.cpp -o $(TARGET)

clean:
	rm -f $(OBJS) $(TARGET)
