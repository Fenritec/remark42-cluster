TAG := $(shell bash getTag.sh)

docker:
	@echo $(shell docker build -t fenritec/remark42-cluster:$(TAG) .)
