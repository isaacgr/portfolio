.PHONY: build install tests test bench clean

BUILD_DIR:=build
MODULE_NAME:=portfolio
INSTALL_DIR:=/usr/local/bin
CONFIG_DIR:=/var/lib/portfolio
BLOG_DIR:=/var/lib/portfolio/posts
AUTH_DIR:=/var/lib/portfolio/authentication
VERSION:= $(shell scripts/version.sh)

build:
	mkdir -p $(BUILD_DIR)
	cp ./docs/example.config.ini $(BUILD_DIR)/config.ini
	cp -r ./web $(BUILD_DIR)/
	@bash scripts/update-file-hash.sh $(BUILD_DIR)
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -gcflags "-c=16" -mod=vendor -o $(BUILD_DIR)/$(MODULE_NAME)-$(VERSION) ./cmd/$(MODULE_NAME)

install:
	mkdir -p $(CONFIG_DIR) $(BLOG_DIR) $(AUTH_DIR)
	cp $(BUILD_DIR)/$(MODULE_NAME)-$(VERSION) $(INSTALL_DIR)/
	
	# We dont want to overwrite the config.ini if it already exists
	# since it has a lot of annoying parameters to fill in
	cp -n $(BUILD_DIR)/config.ini $(CONFIG_DIR)/

	cp -r $(BUILD_DIR)/web $(CONFIG_DIR)/
	cp -r posts $(BLOG_DIR)/
	ln -sf $(INSTALL_DIR)/$(MODULE_NAME)-$(VERSION) $(INSTALL_DIR)/$(MODULE_NAME)

	@bash scripts/install.sh

uninstall:
	systemctl stop portfolio.service
	rm /usr/local/bin/portfolio*
	rm -r /var/lib/portfolio

tests: test bench

test:
	go test -race -cover -bench=^$$ ./...

bench:
	go test -run=^$$ -bench=. -benchmem -count=5 -cpu=1,8,16 ./...

distclean:: clean
	rm -rf $(BUILD_DIR)
maintainer-clean:: distclean

