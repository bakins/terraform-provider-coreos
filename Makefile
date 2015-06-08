PREFIX := /usr/local

all: bin/terraform-provider-coreos

bin/terraform-provider-coreos:
	gb build all

clean:
	rm bin/*

distclean: clean
	rm -rf pkg

install: bin/terraform-provider-coreos
	install -m 755 -d $(PREFIX)/bin
	install -m 755 $< $(PREFIX)/bin/terraform-provider-coreos
