import pulumi
import pulumi_kubernetes as kubernetes
import crds.gatewayapi.pulumi_crds as gatewayapi
from typing import Optional, List

config = pulumi.Config()

service_name = config.require("name")
gateway_class_name = config.require("gateway_class_name")

k8s_provider_config = config.require_object("k8s_provider")
is_canary = config.get("deployment_strategy", "") == "canary"

k8s_namespace = service_name
if is_canary:
    k8s_namespace = f"{k8s_namespace}-canary"

k8s_provider = kubernetes.Provider(
    "k8s-provider",
    namespace=k8s_namespace,
    enable_server_side_apply=True,
    kubeconfig=k8s_provider_config["kube_config"],
    cluster_identifier=k8s_provider_config["cluster_identifier"]
)

# Register Gateway API CRDs
gateway_api_crd = kubernetes.yaml.ConfigFile(f"gw-api-crds-{service_name}",
    file="gateway-api.crd.yaml",
    opts=pulumi.ResourceOptions(provider=k8s_provider)
)

service_labels = {
    "service": service_name,
}

# Create a namespace (user supplies the name of the namespace)
namespace = kubernetes.core.v1.Namespace(
    f"ns-{service_name}",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels=service_labels,
        name=k8s_namespace
    ),
    opts=pulumi.ResourceOptions(provider=k8s_provider)
)

gateway = gatewayapi.gateway.v1.Gateway(
    f"gateway-{service_name}",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels=service_labels,
        name=service_name,
    ),
    spec=gatewayapi.gateway.v1.GatewaySpecArgs(
        gateway_class_name=gateway_class_name,
        listeners=[
            gatewayapi.gateway.v1.GatewaySpecListenersArgs(
                name="http",
                protocol="HTTP",
                port=80,
                allowed_routes=gatewayapi.gateway.v1.GatewaySpecListenersAllowedRoutesArgs(
                    namespaces=gatewayapi.gateway.v1.GatewaySpecListenersAllowedRoutesNamespacesArgs(
                        from_="All"
                    )
                )
            )
        ],
    ),
    opts=pulumi.ResourceOptions(provider=k8s_provider)
)

# export useful outputs
pulumi.export("gateway_name", gateway.metadata.name)
pulumi.export("gateway_namespace", gateway.metadata.namespace)

# TODO: Expose the gateway using AWS LB Controller annotations
# Follow this doc: https://istio.io/latest/docs/tasks/traffic-management/ingress/gateway-api/#resource-attachment-and-scaling

# TODO: Resource attachment and scaling
# https://istio.io/latest/docs/tasks/traffic-management/ingress/gateway-api/#resource-attachment-and-scaling
