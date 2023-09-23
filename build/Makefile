TAG ?= test
SUBDIRS := service-unit stateful-unit/mongo load-generator tbctl

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ prod TAG=$(TAG)
