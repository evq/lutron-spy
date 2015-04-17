.PHONY: clean install

all: lutron-spy

lutron-spy: lutron-spy.go
	goxc

install: lutron-spy
	mkdir -p $(DESTDIR)/usr/local/bin
	mkdir -p $(DESTDIR)/usr/local/share/lutron-spy
	mkdir -p $(DESTDIR)/etc/init.d
	cp lutron-spy $(DESTDIR)/usr/local/bin/lutron-spy
	cp example-config.json $(DESTDIR)/usr/local/share/lutron-spy/example-config.json
	cp S59lutron-spy $(DESTDIR)/etc/init.d/S59lutron-spy

clean:
	rm -f lutron-spy
