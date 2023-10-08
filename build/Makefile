TAG ?= test
SUBDIRS := service-unit datagen stateful-unit/mongo stateful/redis load-generator tbctl

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ prod TAG=$(TAG)
