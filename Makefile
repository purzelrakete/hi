.PHONY: test

test:
	make -C percolate test
	make -C spectrogram test
