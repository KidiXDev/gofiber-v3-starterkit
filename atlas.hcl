variable "envfile" {
  type    = string
  default = ".env"
}

locals {
  envfile = {
    for line in split("\n", file(var.envfile)) : split("=", line)[0] => trimspace(regex("=(.*)", line)[0])
    if !startswith(line, "#") && length(split("=", line)) > 1
  }
}

data "external_schema" "bun" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "bun" {
  src = data.external_schema.bun.url
  dev = local.envfile["DEV_DATABASE_URL"]
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  url = local.envfile["DATABASE_URL"]
  dev = local.envfile["DEV_DATABASE_URL"]
  migration {
    dir = "file://migrations"
  }
}
