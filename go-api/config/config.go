package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Config struct {
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	Port               string
}

var AppConfig *Config
var SupabaseClient *supabase.Client

func Load() {
	//Loads environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		Port:               getEnv("PORT", "8080"),
	}

	if AppConfig.SupabaseURL == "" || AppConfig.SupabaseAnonKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_ANON_KEY are required")
	}

	SupabaseClient = InitSupabase()

	log.Println("Config loaded successfully")
}

func InitSupabase() *supabase.Client {
	client, err := supabase.NewClient(AppConfig.SupabaseURL, AppConfig.SupabaseServiceKey, nil)
	if err != nil {
		log.Fatal("Failed to initialize Supabase client:", err)
	}
	return client
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

