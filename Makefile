.PHONY: cryptanalysis

cryptanalysis:
	cd runners/go-driver && go build -o ../../cryptanalysis
	$(MAKE) -C encoders/nejati-collision
	$(MAKE) -C encoders/nejati-preimage
