BENCH_SRCS := $(shell ls benchmark/*)

benchmark: benchmark/*
	make -C benchmark


.PHONY: rm-enc
rm-enc:
	rm -f data/*.enc
