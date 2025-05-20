package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/m1guelpf/chatgpt-telegram/internal/payment"
	"github.com/m1guelpf/chatgpt-telegram/src/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.LoadEnvConfig(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := cfg.ValidateWithDefaults(); err != nil {
		log.Fatalf("validation: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}

	fk := payment.NewFreeKassa(cfg.FKMerchantID, cfg.FKSecret1, cfg.FKSecret2)
	store := payment.NewStore()

	http.HandleFunc("/freekassa/callback", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if !fk.Verify(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		orderID := r.FormValue("MERCHANT_ORDER_ID")
		tgID, ok := store.Get(orderID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		store.Delete(orderID)

		// ---------- выдаём доступ пользователю ----------
		params := map[string]string{
			"chat_id": strconv.FormatInt(cfg.PaidChannelID, 10),
		}

		resp, err := bot.MakeRequest("exportChatInviteLink", params)
		if err != nil {
			log.Printf("exportChatInviteLink error: %v", err)
		} else {
			// resp.Result — это json.RawMessage, достаём строку-ссылку
			var link string
			if err := json.Unmarshal(resp.Result, &link); err != nil {
				log.Printf("unmarshal invite link: %v", err)
			} else {
				days := strconv.Itoa(cfg.AccessDays)
				text := fmt.Sprintf(
					"✅ Оплата прошла успешно. Доступ выдан на %s дн.\n🔗 Ваша ссылка: %s",
					days, link,
				)
				bot.Send(tgbotapi.NewMessage(tgID, text))
			}
		}

		w.Write([]byte("YES"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Printf("webhook listening on :%s", port)
		http.ListenAndServe(":" + port, nil)
	}()

	// === Bot updates ===
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for upd := range updates {
		if upd.Message == nil {
			continue
		}
		switch upd.Message.Command() {
		case "start":
			handleStart(bot, upd.Message, fk, cfg, store)
		case "buy":
			handleBuy(bot, upd.Message, fk, cfg, store)
		case "help":
			bot.Send(tgbotapi.NewMessage(upd.Message.Chat.ID, "Команды: /buy – купить доступ"))
		}
	}
}

func handleStart(bot *tgbotapi.BotAPI, m *tgbotapi.Message, fk *payment.FreeKassa, cfg *config.EnvConfig, store *payment.Store) {
	text := "👋 Добро пожаловать!\nЧтобы получить доступ к закрытому каналу нажмите /buy"
	bot.Send(tgbotapi.NewMessage(m.Chat.ID, text))
}

func handleBuy(bot *tgbotapi.BotAPI, m *tgbotapi.Message, fk *payment.FreeKassa, cfg *config.EnvConfig, store *payment.Store) {
	orderID := payment.NewOrderID(m.From.ID)
	store.Put(orderID, m.From.ID)
	url := fk.GenerateURL(cfg.ProductPrice, orderID, m.From.UserName+"@t.me")
	msg := tgbotapi.NewMessage(m.Chat.ID, "💳 Для оплаты перейдите по ссылке:\n"+url)
	bot.Send(msg)
}
