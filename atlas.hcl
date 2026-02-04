# Atlas configuration file
# See https://atlasgo.io for documentation

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/domain",
    "--dialect", "postgres",
  ]
}

# Environment variables
variable "db_host" {
  type    = string
  default = getenv("DB_HOST") != "" ? getenv("DB_HOST") : "localhost"
}

variable "db_port" {
  type    = string
  default = getenv("DB_PORT") != "" ? getenv("DB_PORT") : "5432"
}

variable "db_user" {
  type    = string
  default = getenv("DB_USER") != "" ? getenv("DB_USER") : "postgres"
}

variable "db_password" {
  type    = string
  default = getenv("DB_PASSWORD") != "" ? getenv("DB_PASSWORD") : ""
}

variable "db_name" {
  type    = string
  default = getenv("DB_NAME") != "" ? getenv("DB_NAME") : "auth_service"
}

variable "db_sslmode" {
  type    = string
  default = getenv("DB_SSLMODE") != "" ? getenv("DB_SSLMODE") : "disable"
}

# Environment definitions
env "local" {
  src = data.external_schema.gorm.url
  
  # Database URL with explicit sslmode
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"
  
  # Migration directory
  migration {
    dir = "file://migrations"
  }
  
  # Dev database for schema diffing
  dev = "docker://postgres/15/dev?search_path=public"
}

env "dev" {
  src = data.external_schema.gorm.url
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"
  
  migration {
    dir = "file://migrations"
  }
  
  dev = "docker://postgres/15/dev?search_path=public"
}

env "prod" {
  src = data.external_schema.gorm.url
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=${var.db_sslmode}"
  
  migration {
    dir = "file://migrations"
    revisions_schema = "atlas_schema_revisions"
  }
  
  # Lint policies for production
  lint {
    destructive {
      error = true
    }
  }
}
