
CC := gcc

CFLAGS := -msse4.2 -fno-strict-aliasing
LDFLAGS:= -lrte_eal -lrte_mbuf -lrte_mempool -lrte_ethdev
TARGET := dns
SRCS := main.c
OBJS := $(SRCS:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJS) $(TARGET)