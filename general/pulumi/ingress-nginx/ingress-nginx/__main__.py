import pulumi
import pulumi_kubernetes as kubernetes

config = pulumi.Config()
k8s_namespace = config.get("k8sNamespace", "default")

# Create a namespace (user supplies the name of the namespace)
sre_sys_ns = kubernetes.core.v1.Namespace(
    "sre-system-ns",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        # labels=app_labels,
        name=k8s_namespace,
    )
)

# Use Helm to install the Nginx ingress controller
external_ingresscontroller = kubernetes.helm.v3.Release(
    "external-ingress-nginx-chart",
    chart="ingress-nginx",
    namespace=sre_sys_ns.metadata.name,
    repository_opts={
        "repo": "https://kubernetes.github.io/ingress-nginx",
    },
    skip_crds=True,
    values={
        "controller": {
            "ingressClass": "external-ingress",
            "ingressClassResource": {
                "controllerValue": "k8s.io/external-ingress-nginx",
                "name": "external-ingress",
            },
            "service": {
                "type": "NodePort",
                "extraLabels": {
                        "app": "external-ingress",
                },
            },
        },
    },
    version="4.12.0"
)

# Export some values for use elsewhere
pulumi.export("name", external_ingresscontroller.name)
