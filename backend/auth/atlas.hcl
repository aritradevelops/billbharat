env "local" {
  # The database that Atlas will diff against
  url = "postgresql://postgres:admin@localhost:5432/billbharat-auth-service?sslmode=disable"

  # Folder where migrations will be generated
  migration {
    dir = "file://internal/persistence/migrations"
  }

  # The schema.sql file to compare against
  dev = "docker://postgres/15/dev" # ephemeral dev db atlas uses for analysis

  # atlas migrate diff business_and_users --env local --to "file://internal/persistence/schemas"
  # atlas migrate apply --env local
}
