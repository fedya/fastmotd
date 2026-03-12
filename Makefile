PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
PROFILE_DIR ?= /etc/profile.d
CONFDIR ?= /etc

BINARY = fastmotd

.PHONY: all build install clean

all: build

build:
	go build -o $(BINARY) .

install: build
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	install -Dm644 fastmotd.sh $(DESTDIR)$(PROFILE_DIR)/fastmotd.sh
	install -Dm644 fastmotd.toml $(DESTDIR)$(CONFDIR)/fastmotd.toml

clean:
	rm -f $(BINARY)
