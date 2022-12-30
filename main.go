package main

import (
	"QiwiSteamPay/database"
	"QiwiSteamPay/qiwi"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PayRequest struct {
	Sum   string `json:"sum"`
	Login string `json:"login"`
	Mail  string `json:"mail"`
}

func main() {
	bills := make(chan string)
	go checkPayments(bills)

	http.HandleFunc("/pay", pay(bills))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pay(ch chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payReq PayRequest
		err := json.NewDecoder(r.Body).Decode(&payReq)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("[ERROR] " + err.Error())
			return
		}
		payLink := qiwi.GetPayLink()
		w.Write([]byte("https://oplata.qiwi.com/create?publicKey=" + payLink.PublicKey +
			"&amount=" + payReq.Sum +
			"&successUrl=" + payLink.SuccessUrl +
			"&billId=" + payLink.BillId +
			"&comment=" + payReq.Login)) // Формирование ссылки оплаты
		ch <- payLink.BillId
	}
}

func checkPayments(bills chan string) {
	for m := range bills {
		billId := m
		go func() {
			log.Printf("[INFO] Запуск проверки платежа: %s", billId)
			status := qiwi.CheckPaymentStatus(billId, time.Now())
			if status.Status.Value == "PAID" {
				transactId := database.AddNewProduct(status.Comment, status.Amount.Value)
				floatSum, _ := strconv.ParseFloat(status.Amount.Value, 32)
				amount := qiwi.GetCurrencySum(floatSum)
				p2pResp, err := qiwi.P2P(amount, transactId)
				if err != nil {
					log.Printf("[ERROR] Ошибка отправки запроса на перевод между счетам для аккаунта: %s", status.Comment)
				}
				log.Printf("[INFO] Статус перевода между валютой аккаунта %s: %s", status.Comment, p2pResp.Status)
				time.Sleep(30 * time.Second)
				resp, total := qiwi.SendMoneyToSteam(amount, transactId, status.Comment)
				if total != nil {
					log.Println(total.Error())
				}
				log.Printf("[INFO] Статус перевода на steam кошелек аккаунта %s: %s", status.Comment, resp.Status)
			}
			log.Printf("[INFO] Платеж %s имеет статус %s", billId, status.Status.Value)
		}()
	}
}
