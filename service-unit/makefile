TAG ?= v0.0.3
SUBDIRS := stateless stateful/mongo load-generator

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ prod TAG=$(TAG)
