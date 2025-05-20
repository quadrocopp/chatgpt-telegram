package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	TelegramID      []int64 `mapstructure:"TELEGRAM_ID"`
	TelegramToken   string  `mapstructure:"TELEGRAM_TOKEN"`
	EditWaitSeconds int     `mapstructure:"EDIT_WAIT_SECONDS"`

	// FreeKassa
	FKMerchantID string `mapstructure:"FK_MERCHANT_ID"`
	FKSecret1    string `mapstructure:"FK_SECRET_1"`
	FKSecret2    string `mapstructure:"FK_SECRET_2"`

	// Продукт
	ProductPrice  float64 `mapstructure:"PRODUCT_PRICE"`
	AccessDays    int     `mapstructure:"ACCESS_DAYS"`
	PaidChannelID int64   `mapstructure:"PAID_CHANNEL_ID"`
}

// HasTelegramID true, если id находится в списке админов
func (c *EnvConfig) HasTelegramID(id int64) bool {
	for _, v := range c.TelegramID {
		if v == id {
			return true
		}
	}
	return false
}

// LoadEnvConfig читает .env и переменные окружения
func LoadEnvConfig(path string) (*EnvConfig, error) {
	_ = godotenv.Load(path) // .env → в ENV

	v := viper.New()
	v.AutomaticEnv() // будет брать из ENV

	// явно привязываем каждую переменную к ключу Viper
	v.BindEnv("TELEGRAM_ID")
	v.BindEnv("TELEGRAM_TOKEN")
	v.BindEnv("EDIT_WAIT_SECONDS")

	v.BindEnv("FK_MERCHANT_ID")
	v.BindEnv("FK_SECRET_1")
	v.BindEnv("FK_SECRET_2")

	v.BindEnv("PRODUCT_PRICE")
	v.BindEnv("ACCESS_DAYS")
	v.BindEnv("PAID_CHANNEL_ID")

	// дефолты (можно оставить, можно убрать)
	v.SetDefault("EDIT_WAIT_SECONDS", 3)
	v.SetDefault("PRODUCT_PRICE", 100.0)
	v.SetDefault("ACCESS_DAYS", 1)

	var cfg EnvConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// (опционально) проверка обязательных полей
func (c *EnvConfig) ValidateWithDefaults() error {
	if c.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_TOKEN required")
	}
	if c.FKMerchantID == "" || c.FKSecret1 == "" || c.FKSecret2 == "" {
		return fmt.Errorf("FreeKassa credentials required")
	}
	return nil
}
