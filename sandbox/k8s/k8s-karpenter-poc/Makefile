export KARPENTER_VERSION = v0.18.1
export CLUSTER_NAME = karpenter-poc
export AWS_DEFAULT_REGION = us-east-1
export AWS_ACCOUNT_ID = $(shell aws sts get-caller-identity --query Account --output text)

.PHONY: create_cluster

create_cluster:
	@command eksctl >/dev/null || ( echo "eksctl is not installed. Install it and try again"; exit 1 )
	@echo "=> AWS: Account ${AWS_ACCOUNT_ID}, Region ${AWS_DEFAULT_REGION}"
	@echo "=> Creating EKS cluster ${CLUSTER_NAME} with Karpenter ${KARPENTER_VERSION}"
	@set -e; \
	FILE=$$(mktemp); \
	cat cluster.yaml > $$FILE; \
	sed -i "" "s/CLUSTER_NAME/$(CLUSTER_NAME)/g" $$FILE; \
	sed -i "" "s/AWS_DEFAULT_REGION/$(AWS_DEFAULT_REGION)/g" $$FILE; \
	cat $$FILE | eksctl create cluster -f -



