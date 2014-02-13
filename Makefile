.PHONY: test

test:
	make -C percolate test
	make -C spectrogram test
	make -C words test
	make -C fftw test
