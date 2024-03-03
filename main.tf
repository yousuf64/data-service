terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.1"
    }
  }
}

provider "docker" {}

# Pull image
resource "docker_image" "envoy" {
  name         = "envoyproxy/envoy:dev"
  keep_locally = true
}

# Create container
resource "docker_container" "envoy" {
  image = docker_image.envoy.image_id
  name  = "envoy-load-balancer"
  ports {
    internal = 80
    external = 80
  }
  mounts {
    source = "${path.cwd}/envoy.yaml"
    target = "/envoy.yaml"
    type   = "bind"
  }
  command = ["-c", "envoy.yaml"]
}

