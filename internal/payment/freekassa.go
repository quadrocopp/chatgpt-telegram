package payment

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "net/http"

    "github.com/google/uuid"
)

type FreeKassa struct {
    MerchantID string
    Secret1    string
    Secret2    string
}

func NewFreeKassa(id, s1, s2 string) *FreeKassa {
    return &FreeKassa{id, s1, s2}
}

func (fk *FreeKassa) GenerateURL(amount float64, orderID, email string) string {
    amountStr := fmt.Sprintf("%.2f", amount) // 10.00
    currency  := "RUB"                      // если хотите менять – передавайте аргументом

    raw := fmt.Sprintf("%s:%s:%s:%s:%s",
        fk.MerchantID, amountStr, fk.Secret1, currency, orderID)
    sign := md5.Sum([]byte(raw))

    return fmt.Sprintf(
        "https://pay.freekassa.ru/?m=%s&oa=%s&currency=%s&o=%s&s=%x&email=%s",
        fk.MerchantID, amountStr, currency, orderID, sign, email)
}

func (fk *FreeKassa) Verify(r *http.Request) bool {
    r.ParseForm()
    amount   := r.FormValue("AMOUNT")
    merchant := r.FormValue("MERCHANT_ID")
    orderID  := r.FormValue("MERCHANT_ORDER_ID")
    sign     := r.FormValue("SIGN")

    // формула для callback: MID:AMOUNT:ORDER_ID:SECRET2
    raw := fmt.Sprintf("%s:%s:%s:%s",
        fk.MerchantID, amount, fk.Secret2, orderID)
    sum := md5.Sum([]byte(raw))
    expect := hex.EncodeToString(sum[:])

    return merchant == fk.MerchantID && sign == expect
}

func (fk *FreeKassa) calcSign(amount interface{}, orderID string, secret string) string {
	/*
		Подпись: md5(MERCHANT_ID:AMOUNT:SECRET:ORDER_ID)
	*/
	raw := fmt.Sprintf("%s:%v:%s:%s", fk.MerchantID, amount, secret, orderID)
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// Генерация уникального заказа (привяжем к Telegram-ID)
func NewOrderID(tgID int64) string {
    return fmt.Sprintf("%d-%s", tgID, uuid.NewString())
}
