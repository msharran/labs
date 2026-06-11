from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import AutoScaling, EC2
from diagrams.aws.network import ELB, Endpoint, InternetGateway, NATGateway, Nacl, PrivateSubnet, PublicSubnet, RouteTable, VPC
from diagrams.aws.storage import S3

# Run from this directory with:
#   uv run python network_architecture.py
# Requires Graphviz installed locally.

with Diagram(
    "AWS Multi-AZ Network Baseline",
    filename="aws_network_baseline",
    outformat="png",
    show=False,
    direction="TB",
):
    s3 = S3("S3")
    internet = InternetGateway("Internet Gateway")

    with Cluster("Region"):
        with Cluster("VPC: 10.0.0.0/16"):
            vpc = VPC("VPC")
            nacl = Nacl("Network ACLs")
            s3_endpoint = Endpoint("S3 Gateway Endpoint")

            with Cluster("Public Route Tables"):
                public_rt_a = RouteTable("public-rt-a\n0.0.0.0/0 -> IGW")
                public_rt_b = RouteTable("public-rt-b\n0.0.0.0/0 -> IGW")

            with Cluster("Private Route Tables"):
                private_rt_a = RouteTable("private-rt-a\n0.0.0.0/0 -> NAT-A\nS3 -> Endpoint")
                private_rt_b = RouteTable("private-rt-b\n0.0.0.0/0 -> NAT-B\nS3 -> Endpoint")

            with Cluster("Availability Zone A"):
                with Cluster("Public Subnet A\n10.0.1.0/24"):
                    public_a = PublicSubnet("public-a")
                    nat_a = NATGateway("NAT Gateway A\n+ Elastic IP")
                    alb_a = ELB("ALB node A")

                with Cluster("Private Subnet A\n10.0.11.0/24"):
                    private_a = PrivateSubnet("private-a")
                    app_a = EC2("server A")

            with Cluster("Availability Zone B"):
                with Cluster("Public Subnet B\n10.0.2.0/24"):
                    public_b = PublicSubnet("public-b")
                    nat_b = NATGateway("NAT Gateway B\n+ Elastic IP")
                    alb_b = ELB("ALB node B")

                with Cluster("Private Subnet B\n10.0.12.0/24"):
                    private_b = PrivateSubnet("private-b")
                    app_b = EC2("server B")

            with Cluster("Application Tier"):
                alb = ELB("Application Load Balancer\nspans public subnets")
                asg = AutoScaling("Auto Scaling Group\nspans private subnets")

            vpc >> nacl

            public_a >> public_rt_a
            public_b >> public_rt_b
            private_a >> private_rt_a
            private_b >> private_rt_b

            public_rt_a >> internet
            public_rt_b >> internet
            public_a >> nat_a
            public_b >> nat_b

            private_rt_a >> nat_a
            private_rt_b >> nat_b
            private_rt_a >> s3_endpoint
            private_rt_b >> s3_endpoint

            alb >> [alb_a, alb_b]
            alb_a >> app_a
            alb_b >> app_b
            asg >> [app_a, app_b]

    internet >> Edge(label="0.0.0.0/0") >> alb
    s3_endpoint >> s3
