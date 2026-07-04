# Configuration

Hardcoding values like `PORT=8080` or `DB_PASS=secret` directly into your source code is a major security and architectural violation. 

Following the **Twelve-Factor App methodology**, all configuration that varies between deployments (Development vs Staging vs Production) must be passed in via **Environment Variables**.

## 1. Native `os.Getenv`

The standard library `os` package is perfectly capable of reading environment variables.

```go
func LoadConfig() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback default
    }
    
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        panic("DATABASE_URL must be provided!")
    }
}
```
While simple, this becomes tedious when you have 50 configuration variables. You have to write `if "" { ... }` 50 times, and `os.Getenv` only returns strings, meaning you have to manually use `strconv.Atoi` for numbers.

## 2. Using `.env` Files (`godotenv`)

When running locally, manually exporting 20 environment variables in your terminal before running `go run main.go` is annoying. 

The industry standard is to use a `.env` file in the root of your project:
```text
PORT=9090
DATABASE_URL=postgres://user:pass@localhost:5432/db
```
*(Warning: NEVER commit the `.env` file to Git! Add it to `.gitignore`!)*

You can load this file into Go's environment using the popular `github.com/joho/godotenv` package.

```go
import "github.com/joho/godotenv"

func main() {
    // Loads the .env file into the system environment
    err := godotenv.Load()
    if err != nil {
        fmt.Println("No .env file found. Relying on system environment.")
    }

    fmt.Println(os.Getenv("PORT")) // 9090
}
```

## 3. Enterprise Configuration (`viper`)

For massive enterprise applications, you need configuration that can read from `.env` files, JSON files, CLI flags, and remote servers (like AWS Secrets Manager), while automatically converting types (strings to ints).

The undisputed king of configuration in the Go ecosystem is **Viper** (`github.com/spf13/viper`).

```go
import "github.com/spf13/viper"

type Config struct {
    Port        int    `mapstructure:"PORT"`
    DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadEnterpriseConfig() (*Config, error) {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv() // Also read from system environment variables
    
    // Set defaults
    viper.SetDefault("PORT", 8080)
    
    // Read the file
    viper.ReadInConfig()

    var config Config
    // Magically unmarshals the env vars directly into the struct, 
    // automatically converting strings to integers!
    err := viper.Unmarshal(&config)
    
    return &config, err
}
```
Viper is so powerful that it powers the Kubernetes CLI and Docker. It is the gold standard for Go configuration.
