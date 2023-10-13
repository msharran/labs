service {
  desired_count                 = 1
  path                          = "/"
  alb_scheme                    = "internal"
  health_check_path             = "/metrics"
  health_check_interval_seconds = 45
  unhealthy_threshold_count     = 3
  health_check_timeout_seconds  = 25
  healthy_threshold_count       = 3

  container {
    # Port on which the service listens
    # NOTE: EXPOSE 9900 in Dockerfile
    port = 9900

    # CPU millis required, 1 CPU core = 1024 millis
    cpu = 512

    # RAM memory required by service in MB
    memory = 1024

    environment "ENVIRONMENT" {
      value = "central"
    }

    environment "ENVIRONMENT" {
      value = "central"
    }
  }
}
