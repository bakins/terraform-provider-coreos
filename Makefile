all: bin/terraform-provider-coreos

bin/terraform-provider-coreos:
	gb build all

clean:
	rm bin/*

distclean: clean
	rm -rf pkg
