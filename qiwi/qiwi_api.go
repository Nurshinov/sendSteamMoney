package qiwi

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type PayLink struct {
	PublicKey  string
	SuccessUrl string
	BillId     string
}

type BillStatus struct {
	SiteId string `json:"siteId"`
	BillId string `json:"billId"`
	Amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Status struct {
		Value           string    `json:"value"`
		ChangedDateTime time.Time `json:"changedDateTime"`
	} `json:"status"`
	Comment            string    `json:"comment"`
	CreationDateTime   time.Time `json:"creationDateTime"`
	ExpirationDateTime time.Time `json:"expirationDateTime"`
	PayUrl             string    `json:"payUrl"`
}

type Currency struct {
	Result []struct {
		Set  string  `json:"set"`
		From string  `json:"from"`
		To   string  `json:"to"`
		Rate float64 `json:"rate"`
	} `json:"result"`
}

// Проверка статуса платежа
func CheckPaymentStatus(bill string, evT time.Time) *BillStatus {
	var t BillStatus
	client := http.Client{}
	timeout := evT.Add(10 * time.Minute)
	req, err := http.NewRequest("GET", "https://api.qiwi.com/partner/bill/v1/bills/"+bill, nil)
	if err != nil {
		log.Println("[ERROR] Ошибка проверки статуса платежа ID = " + bill + "\n" + err.Error())
	}
	req.Header.Add("Authorization", os.Getenv("SECRET_KEY"))
	for time.Now().Before(timeout) {
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}
		err = json.NewDecoder(resp.Body).Decode(&t)
		if t.Status.Value == "PAID" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return &t
}

// Функция получения курса валюты, в данном случае курс тенге
func GetCurrencySum(sum float64) float64 {
	var t Currency
	var rate float64
	client := http.Client{}
	req, _ := http.NewRequest("GET", "https://edge.qiwi.com/sinap/crossRates", nil)
	req.Header.Add("Authorization", os.Getenv("QIWI_WALLET_API"))
	resp, err := client.Do(req)
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range t.Result {
		if v.From == "643" && v.To == "398" {
			rate = v.Rate
			break
		}
	}
	tingeAmount := sum / rate
	return tingeAmount
}

// Функция перевода между счетами счетом в рублях и счетом в тенге
func P2P(amount float64, id string) (*http.Response, error) {
	client := http.Client{}
	body := fmt.Sprintf(
		`{
			"id": "%s",
			"sum": {
				"amount": %f,
				"currency": "398"
			},  
			"paymentMethod": {
				"accountId": "643",
				"type": "Account"
			},
			"fields": {
				"account": "79170278608"
			}
		}`, id, amount)
	req, _ := http.NewRequest("POST", "https://edge.qiwi.com/sinap/api/v2/terms/99/payments", strings.NewReader(body))
	req.Header.Add("Authorization", os.Getenv("QIWI_WALLET_API"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}

func SendMoneyToSteam(amount float64, id string, account string) (*http.Response, error) {
	client := http.Client{}
	body := fmt.Sprintf(`{
		"id": "%s",
		"sum": {
			"amount": %f,
			"currency": "398"
		},
		"paymentMethod": {
			"accountId": "398",
			"type": "Account"
		},
		"comment": "",
		"fields": {
			"account": "%s"
			
		}
	}`, id, amount, account)
	req, _ := http.NewRequest("POST", "https://edge.qiwi.com/sinap/api/v2/terms/31212/payments", strings.NewReader(body))
	req.Header.Add("Authorization", os.Getenv("QIWI_WALLET_API"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}

func GetPayLink() *PayLink {
	var payLink PayLink
	payLink.PublicKey = os.Getenv("PUBLIC_KEY")
	payLink.SuccessUrl = "http://localhost:8080"
	payLink.BillId = uuid.New().String()
	return &payLink
}
