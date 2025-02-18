import pulumi
import pulumi_kubernetes as kubernetes
from typing import Optional, List

config = pulumi.Config()

k8s_provider = kubernetes.Provider(
    "k8s-provider",
    enable_server_side_apply=True,
    kubeconfig=config.require_object("k8s_provider")["kube_config"],
    cluster_identifier=config.require_object("k8s_provider")["cluster_identifier"]
)

for definition in config.require_object("definitions"):
    crd_definition = kubernetes.yaml.ConfigFile(f"{definition['name']}", file=definition['yaml_file'])
