package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper

	OpenAISession string
}

// LoadOrCreatePersistentConfig uses the default config directory for the current OS
// to load or create a config file named "chatgpt.json"
func LoadOrCreatePersistentConfig() (*Config, error) {
    v := viper.New()

    // Получаем путь ~/.config
    configPath, err := os.UserConfigDir()
    if err != nil {
        return nil, fmt.Errorf("Couldn't find user config dir: %v", err)
    }

    // Настраиваем viper
    v.AddConfigPath(configPath)
    v.SetConfigName("chatgpt")       // файл chatgpt.json
    v.SetConfigType("json")

    // Пытаемся прочитать существующий конфиг
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Если файла нет — сначала создаём директорию
            if err := os.MkdirAll(configPath, 0o755); err != nil {
                return nil, fmt.Errorf("Couldn't create config dir: %v", err)
            }
            // А потом — сам файл
            if err := v.SafeWriteConfig(); err != nil {
                return nil, fmt.Errorf("Couldn't create config file: %v", err)
            }
        } else {
            return nil, fmt.Errorf("Couldn't load config: %v", err)
        }
    }

    // Теперь конфиг гарантированно существует — читаем его
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("Couldn't re-read config: %v", err)
    }

    // Привязываем значения из файла к нашей структуре Config
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("Couldn't parse config: %v", err)
    }

    return &cfg, nil
}

func (cfg *Config) SetSessionToken(token string) error {
	// key must match the struct field name
	cfg.v.Set("OpenAISession", token)
	cfg.OpenAISession = token
	return cfg.v.WriteConfig()
}
