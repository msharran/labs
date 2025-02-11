import pulumi
import pulumi_kubernetes as kubernetes
from typing import Optional, List

config = pulumi.Config()

service_name = config.require("name")
image = config.get("image", None)
port = config.get_int("port", 80)

is_canary = config.get("deployment_strategy", "") == "canary"
canary_weight = config.get_int("canary_weight", 0)

provider_type = config.get("provider_type", None)
helm_values = config.get_object("helm_values", {})
disable_ingress = config.get_bool("ingress_disabled", False)
ingress_class_name = config.get("ingress_class_name","")


k8s_namespace = service_name
if is_canary:
    k8s_namespace = f"canary-{service_name}"

service_labels = {
    "service": service_name,
    "is_canary": "true" if is_canary else "false",
}

# Create a namespace (user supplies the name of the namespace)
namespace = kubernetes.core.v1.Namespace(
    f"ns-{service_name}",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels=service_labels,
        name=k8s_namespace,
    )
)

def ingress_rule(name=service_name, suffix="", port=80) -> kubernetes.networking.v1.IngressRuleArgs:
    return kubernetes.networking.v1.IngressRuleArgs(
        host=f"{name}.example.com",
        http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
            paths=[kubernetes.networking.v1.HTTPIngressPathArgs(
                path="/",
                path_type="Prefix",
                backend=kubernetes.networking.v1.IngressBackendArgs(
                    service=kubernetes.networking.v1.IngressServiceBackendArgs(
                        name=f"{name}{suffix}",
                        port=kubernetes.networking.v1.ServiceBackendPortArgs(
                            number=port
                        )
                    )
                )
            )]
        )
    )

def create_ingress(rules: List[kubernetes.networking.v1.IngressRuleArgs], opts: Optional[pulumi.ResourceOptions] = None):
    if ingress_class_name == "":
        pulumi.error("missing ingress_class_name")
        raise Exception("Ingress error")

    canary_annotations = {}
    if is_canary:
        canary_annotations = {
            "nginx.ingress.kubernetes.io/canary": "true",
            "nginx.ingress.kubernetes.io/canary-weight": f"{canary_weight}"
        }
    ingress = kubernetes.networking.v1.Ingress(
            f"ingress-{service_name}",
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                namespace=k8s_namespace,
                name=f"{service_name}",
                annotations={
                    **canary_annotations,
                }
            ),
            spec=kubernetes.networking.v1.IngressSpecArgs(
                ingress_class_name=ingress_class_name,
                rules=rules
            )
        )

if provider_type == "helm":
    helm = config.require_object("helm")
    helm_release = kubernetes.helm.v3.Release(
        service_name,
        chart=helm["chart"],
        namespace=namespace.metadata.name,
        repository_opts=helm["repository_opts"],
        skip_await=False, # required by ingress which depends on this
        values=helm["values"],
        version=helm["version"],
    )


    pulumi.export("name", helm_release.name)
    pulumi.export("status", helm_release.status)

    if not disable_ingress:
        rules = [
            kubernetes.networking.v1.IngressRuleArgs(
                host=f"{service_name}{rule["name_suffix"]}.example.com",
                http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                    paths=[
                        kubernetes.networking.v1.HTTPIngressPathArgs(
                        path="/",
                        path_type="Prefix",
                        backend=kubernetes.networking.v1.IngressBackendArgs(
                            service=kubernetes.networking.v1.IngressServiceBackendArgs(
                                name= pulumi.Output.concat(
                                    helm_release.status.name,
                                    rule["name_suffix"],
                                ),
                                port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                    number=rule["port"]
                                )
                            )
                        )
                    )]
                )
            )
            for rule in helm["ingress_rules"]
        ]

        create_ingress(
            rules=rules,
            opts=pulumi.ResourceOptions(depends_on=helm_release)
        )
    else:
        pulumi.info("ingress disabled, skipping creation")
else:
    if not image:
        pulumi.error("failed! no image supplied")
        raise Exception("missing image")

    deployment = kubernetes.apps.v1.Deployment(
        f"deployment-{service_name}",
        metadata=kubernetes.meta.v1.ObjectMetaArgs(
            namespace=k8s_namespace,
            name=f"{service_name}"
        ),
        spec=kubernetes.apps.v1.DeploymentSpecArgs(
            replicas=1,
            selector=kubernetes.meta.v1.LabelSelectorArgs(
                match_labels={
                    **service_labels
                }
            ),
            template=kubernetes.core.v1.PodTemplateSpecArgs(
                metadata=kubernetes.meta.v1.ObjectMetaArgs(
                    labels={
                        **service_labels
                    }
                ),
                spec=kubernetes.core.v1.PodSpecArgs(
                    containers=[kubernetes.core.v1.ContainerArgs(
                        name=service_name,
                        image=image,
                        ports=[kubernetes.core.v1.ContainerPortArgs(
                            container_port=port
                        )]
                    )]
                )
            )
        )
    )

    # Create the Nginx service
    service = kubernetes.core.v1.Service(
        f"service-{service_name}",
        metadata=kubernetes.meta.v1.ObjectMetaArgs(
            namespace=k8s_namespace,
            name=f"{service_name}"
        ),
        spec=kubernetes.core.v1.ServiceSpecArgs(
            selector={
                **service_labels
            },
            ports=[kubernetes.core.v1.ServicePortArgs(
                port=80,
                target_port=port
            )]
        )
    )
    if not disable_ingress:
        create_ingress(
            rules=[
                kubernetes.networking.v1.IngressRuleArgs(
                    host=f"{service_name}.example.com",
                    http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                        paths=[
                            kubernetes.networking.v1.HTTPIngressPathArgs(
                                path="/",
                                path_type="Prefix",
                                backend=kubernetes.networking.v1.IngressBackendArgs(
                                    service=kubernetes.networking.v1.IngressServiceBackendArgs(
                                        name=service_name,
                                        port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                            number=port
                                        )
                                    )
                                )
                            )]
                    )
                )
            ],
            opts=pulumi.ResourceOptions(depends_on=service)
        )
    else:
        pulumi.info("ingress disabled, skipping creation")

    # Create the Nginx ingress
