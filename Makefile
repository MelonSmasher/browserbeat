BEAT_NAME=browserbeat
BEAT_PATH=github.com/MelonSmasher/browserbeat
BEAT_GOPATH=$(firstword $(subst :, ,${GOPATH}))
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
LIBBEAT_MAKEFILE=$(ES_BEATS)/libbeat/scripts/Makefile
GOPACKAGES=$(shell govendor list -no-status +local)
GOBUILD_FLAGS=-i -ldflags "-X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.buildTime=$(NOW) -X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.commit=$(COMMIT_ID)"
MAGE_IMPORT_PATH=${BEAT_PATH}/vendor/github.com/magefile/mage
NO_COLLECT=true
CHECK_HEADERS_DISABLED=true

# Path to the libbeat Makefile
-include $(LIBBEAT_MAKEFILE)

# Initial beat setup
.PHONY: setup
setup: pre-setup git-add

pre-setup: copy-vendor git-init
	$(MAKE) -f $(LIBBEAT_MAKEFILE) mage ES_BEATS=$(ES_BEATS)
	$(MAKE) -f $(LIBBEAT_MAKEFILE) update BEAT_NAME=$(BEAT_NAME) ES_BEATS=$(ES_BEATS) NO_COLLECT=$(NO_COLLECT)

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/beats
	git archive --remote ${BEAT_GOPATH}/src/github.com/elastic/beats HEAD | tar -x --exclude=x-pack -C vendor/github.com/elastic/beats
	mkdir -p vendor/github.com/magefile
	cp -R vendor/github.com/elastic/beats/vendor/github.com/magefile/mage vendor/github.com/magefile

.PHONY: git-init
git-init:
	git init

.PHONY: git-add
git-add:
	git add -A
	git commit -m "Add generated browserbeat files"


BEAT_VERSION=0.0.3
VERSION_QUALIFIER=alpha1

VERSION=${BEAT_VERSION}-${VERSION_QUALIFIER}
VERSION_GOFILEPATH=${BEAT_PATH}/version/version.go

BEAT_URL=https://${BEAT_PATH}

## @packaging VERSION=x.y.z set the version of the beat to x.y.z
## Override the built in set_version cmd
set_version:
	dev-tools/set_version ${VERSION}
	$(shell cp ./pack.sh.example ./pack.sh)
	$(shell chmod +x ./pack.sh)
	$(shell sed -i 's/0.0.0/${BEAT_VERSION}/g' ./pack.sh)
	$(shell sed -i 's/alpha/${VERSION_QUALIFIER}/g' ./pack.sh)

tag:
	$(shell git tag v${VERSION})
