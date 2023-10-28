TAG ?= test
SUBDIRS := service-unit datagen stateful-unit/mongo stateful-unit/redis load-generator hexctl

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ prod TAG=$(TAG)
