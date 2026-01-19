package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/golang-web/models"
	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/joho/godotenv"
	"github.com/plutov/paypal/v4"
)

func Payments_home(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/payments/home.html", utilities.Frontend))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{

		"css":     css_sesion,
		"message": css_mensaje,
	}
	template.Execute(response, data)
}

func returnOrderPaypal(token string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	order := time.Now().Format("150405")
	data := `{
  "intent": "CAPTURE",
  "payment_source": {
    "paypal": {
      "experience_context": {
        "payment_method_preference": "IMMEDIATE_PAYMENT_REQUIRED",
        "landing_page": "LOGIN",
        "shipping_preference": "GET_FROM_FILE",
        "user_action": "PAY_NOW",
        "return_url": "https://example.com/returnUrl",
        "cancel_url": "https://example.com/cancelUrl"
      }
    }
  },
  "purchase_units": [
    {
      "invoice_id": "90210",
      "amount": {
        "currency_code": "USD",
        "value": "230.00",
        "breakdown": {
          "item_total": {
            "currency_code": "USD",
            "value": "220.00"
          },
          "shipping": {
            "currency_code": "USD",
            "value": "10.00"
          }
        }
      },
      "items": [
        {
          "name": "T-Shirt",
          "description": "Super Fresh Shirt",
          "unit_amount": {
            "currency_code": "USD",
            "value": "20.00"
          },
          "quantity": "1",
          "category": "PHYSICAL_GOODS",
          "sku": "sku01",
          "image_url": "https://example.com/static/images/items/1/tshirt_green.jpg",
          "url": "https://example.com/url-to-the-item-being-purchased-1",
          "upc": {
            "type": "UPC-A",
            "code": "123456789012"
          }
        },
        {
          "name": "Shoes",
          "description": "Running, Size 10.5",
          "sku": "sku02",
          "unit_amount": {
            "currency_code": "USD",
            "value": "100.00"
          },
          "quantity": "2",
          "category": "PHYSICAL_GOODS",
          "image_url": "https://example.com/static/images/items/1/shoes_running.jpg",
          "url": "https://example.com/url-to-the-item-being-purchased-2",
          "upc": {
            "type": "UPC-A",
            "code": "987654321012"
          }
        }
      ]
    }
  ]
}`

	byte_arr := []byte(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("PAYPAL_BASE_URI")+"/v2/checkout/orders", bytes.NewBuffer(byte_arr))
	if err != nil {
		fmt.Println(err, "err 1")
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("PayPal-Request-Id", "order_"+order)
	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err, "err 2")
		return ""
	}
	defer reg.Body.Close()

	body, err := io.ReadAll(reg.Body)
	paypal := models.PaypalOrderResponseModel{}
	err = json.Unmarshal(body, &paypal)
	if err != nil {
		fmt.Println(err, "err 3")
		return ""
	}

	fmt.Println(paypal, "PAYPALL")
	return paypal.Id
}

func returnScreenshotPaypal(token string, param string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("godotenv.Load:", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("PAYPAL_BASE_URI")+"/v2/checkout/orders/"+param+"/capture", nil)
	if err != nil {
		fmt.Println(err, "err 1")
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err, "err 2")
		return ""
	}
	defer reg.Body.Close()
	body, _ := io.ReadAll(reg.Body)
	fmt.Println("capture response:", reg.StatusCode, string(body))

	if reg.StatusCode == 422 {
		return "bad"
	}

	return "good"

}

func returnTokenPaypal() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("godotenv.Load:", err)
	}

	c, err := paypal.NewClient(os.Getenv("PAYPAL_CLIENT_ID"), os.Getenv("PAYPAL_CLIENT_SECRET"), paypal.APIBaseSandBox) // or paypal.APIBaseLive
	if err != nil {
		fmt.Println("paypal.NewClient:", err)
		return ""
	}
	c.SetLog(os.Stdout)

	accessToken, err := c.GetAccessToken(context.Background())
	if err != nil {
		fmt.Println("GetAccessToken:", err)
		return ""
	}
	return accessToken.Token
}

func Payments_paypal(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/payments/paypal.html", utilities.Frontend))
	token := returnTokenPaypal()
	paypalId := returnOrderPaypal(token)
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":       css_sesion,
		"message":   css_mensaje,
		"token":     token,
		"paypal_id": paypalId,
	}
	template.Execute(response, data)
}

func Payments_paypal_response(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/payments/paypal_response.html", utilities.Frontend))
	token := returnTokenPaypal()
	state := returnScreenshotPaypal(token, request.URL.Query().Get("token"))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_sesion,
		"message": css_mensaje,
		"token":   request.URL.Query().Get("token"),
		"state":   state,
	}
	template.Execute(response, data)
}
