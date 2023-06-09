package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/jayplus-app/jayplus/internal/driver/models"
	"github.com/jayplus-app/jayplus/pkg/messaging"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/webhook"
)

type dateTimeListRequestPayload struct {
	StartDate   string `json:"dateTime"`
	ServiceType string `json:"serviceType"`
	VehicleType string `json:"vehicleType"`
}

func (app *application) response(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.errorLog.Println(err)
	}
}

type definitionPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

func (app *application) NewVehicleType(w http.ResponseWriter, r *http.Request) {
	app.newDefinition("vehicle-types", w, r)
}

func (app *application) NewServiceType(w http.ResponseWriter, r *http.Request) {
	app.newDefinition("service-types", w, r)
}

func (app *application) newDefinition(defType string, w http.ResponseWriter, r *http.Request) {
	var payload definitionPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	def, err := app.db.NewDefinition(defType, payload.Name, payload.Description, payload.Icon, payload.Order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.response(w, http.StatusOK, def)
}

func (app *application) getDefinitions(defType string, w http.ResponseWriter, r *http.Request) error {
	list, err := app.db.GetDefinitions(defType)
	if err != nil {
		return err
	}

	l := len(list)
	types := make([]*definitionPayload, l, l)
	for i, item := range list {
		types[i] = &definitionPayload{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Icon:        item.Icon,
			Order:       item.Order,
		}
	}

	out := map[string]interface{}{
		"name":  defType,
		"types": types,
	}

	app.response(w, http.StatusOK, out)

	return nil
}

type definitionResponse struct {
	models.Defintion
	Order int `json:""`
}

func (app *application) VehicleTypes(w http.ResponseWriter, r *http.Request) {
	if err := app.getDefinitions("vehicle-types", w, r); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) ServiceTypes(w http.ResponseWriter, r *http.Request) {
	if err := app.getDefinitions("service-types", w, r); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) DateTimeList(w http.ResponseWriter, r *http.Request) {
	var input dateTimeListRequestPayload

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// validate data

	d, err := time.Parse(time.DateOnly, input.StartDate)
	if err != nil {
		app.errorLog.Println(err)
		w.Write([]byte(`{"success": false, "message": "Invalid date format"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

	// Get timeframes for the specified date
	timeslots, err := app.db.GetTimeframes(input.ServiceType, input.VehicleType, d)
	if err != nil {
		app.errorLog.Println(err)
		w.Write([]byte(`{"success": false, "message": "Invalid input"}`))
		return
	}

	out, err := json.Marshal(map[string][]models.Timeslot{
		d.Format(time.DateOnly): timeslots,
	})
	if err != nil {
		app.errorLog.Println(err)
		w.Write([]byte(`{"success": false, "message": "Invalid request"}`))
		return
	}

	w.Write(out)
}

type priceRequestPayload struct {
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Time        time.Time `json:"time"`
}

func (app *application) Price(w http.ResponseWriter, r *http.Request) {
	var input priceRequestPayload

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// TODO: [THREAD:4] Validate input

	price, err := app.db.GetPrice(input.ServiceType, input.VehicleType, input.Time)
	if err != nil {
		app.errorLog.Println(err)
		w.Write([]byte(`{"success": false, "message": "Invalid input"}`))
		return
	}

	out := []byte(strconv.FormatFloat(float64(price), 'f', 2, 64))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

type invoiceRequestPayload struct {
	VehicleType string    `json:"vehicleType"`
	ServiceType string    `json:"serviceType"`
	Time        time.Time `json:"time"`
}

type invoiceResponsePayload struct {
	BillNumber        uint      `json:"billNumber"`
	Time              time.Time `json:"time"`
	CancelledAt       time.Time `json:"cancelledAt"`
	TransactionNumber int       `json:"TransactionNumber"`
	ServiceType       string    `json:"serviceType"`
	VehicleType       string    `json:"vehicleType"`
	ServiceCost       float32   `json:"serviceCost"`
	Discount          float32   `json:"discount"`
	Total             float32   `json:"total"`
	Deposit           float32   `json:"deposit"`
	Remaining         float32   `json:"remaining"`
	PayedAt           time.Time `json:"payedAt"`
}

func (app *application) Invoice(w http.ResponseWriter, r *http.Request) {
	var input invoiceRequestPayload

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorLog.Println(err)
		return
	}

	booking, err := app.db.MakeBooking(input.VehicleType, input.ServiceType, input.Time)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// validate data

	// out := []byte(`{
	// 	"Transaction Number": "13",
	// 	"Bill Number": "37",
	// 	"Type of Service": "Show Room",
	// 	"Vehicle Type": "Sedan",
	// 	"Date": "14 Mar 2023",
	// 	"Time": "15:00",
	// 	"Service Cost": "169.00 $",
	// 	"Discount": "Not Specified",
	// 	"Total": "169.00 $",
	// 	"Deposit": "30.00 $",
	// 	"Remaining": "139.00 $"
	// }`)

	bookingResponse := &invoiceResponsePayload{
		BillNumber:        booking.ID,
		Time:              booking.BookedAt,
		CancelledAt:       booking.CancelledAt,
		TransactionNumber: booking.TransactionNumber,
		ServiceType:       booking.ServiceType,
		VehicleType:       booking.VehicleType,
		ServiceCost:       booking.ServiceCost,
		Discount:          booking.Discount,
		Total:             booking.Total,
		Deposit:           booking.Deposit,
		Remaining:         booking.Remaining,
		PayedAt:           booking.PayedAt,
	}

	out, err := json.Marshal(bookingResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

type payInvoiceRequestPayload struct {
	BillNumber uint   `json:"billNumber"`
	Phone      string `json:"phone"`
}

func (app *application) PayInvoice(w http.ResponseWriter, r *http.Request) {
	var input payInvoiceRequestPayload

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorLog.Println(err)
		return
	}

	// TODO: Validate phone number

	booking, err := app.db.GetBooking(input.BillNumber)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := map[string]interface{}{
		"phone": input.Phone,
	}

	booker, err := app.db.GetOrMakeBooker(data)
	if err != nil {
		if err.Error() == "no rows affected" {
			app.errorLog.Printf("Failed to set booker's phone number on invoice payment request [booking_id: %v] [booker_id: %v] [err: %v]\n", booking.ID, booking.BookerID, err)
			return
		}
	}

	if err := app.db.AssociateBooking(booker, booking); err != nil {
		app.errorLog.Printf("Failed to associate booker and booking on invoice payment request [booking_id: %v] [booker_id: %v] [err: %v]\n", booking.ID, booking.BookerID, err)
	}

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(booking.Deposit * 100)),
		Currency: stripe.String(string(stripe.CurrencyCAD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	app.infoLog.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	result := struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	}

	app.response(w, http.StatusOK, result)
}

func (app *application) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		app.errorLog.Printf("⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, app.config.StripeWHSecret)
	if err != nil {
		app.errorLog.Printf("⚠️  Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			app.errorLog.Printf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		app.infoLog.Printf("Successful payment for %d.", paymentIntent.Amount)
		if err := app.handlePaymentIntentSucceeded(paymentIntent); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			app.errorLog.Printf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handlePaymentMethodAttached(paymentMethod)
	default:
		app.errorLog.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) handlePaymentIntentSucceeded(paymentIntent stripe.PaymentIntent) error {
	app.infoLog.Printf("💰 Payment received! %v\n", paymentIntent.ID)
	bookingID, err := strconv.Atoi(paymentIntent.Metadata["billNumber"])
	if err != nil {
		metadata, err := json.Marshal(paymentIntent.Metadata)
		if err != nil {
			metadata = []byte("invalid metadata")
		}

		app.errorLog.Printf("Failed to update booking status on payment success [metadata: %v] [err: %v]\n", metadata, err)
		return err
	}

	data := map[string]string{
		"payed_at": time.Now().Format(time.RFC3339),
	}

	if err := app.db.UpdateBooking(bookingID, data); err != nil {
		if err.Error() == "no rows affected" {
			app.errorLog.Printf("Failed to update booking status on payment success [booking_id: %v] [err: %v]\n", bookingID, err)
			return err
		}
	}

	b, err := app.db.GetBooking(uint(bookingID))
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	user, err := app.db.GetBooker(b.BookerID)

	if defs, err := app.db.GetDefinitions("sms.payment.success"); err == nil {
		tpl := defs[0].Description // "Successfull Payment\nBill No.: %d"

		if defs, err := app.db.GetDefinitions("sms.sender_number"); err == nil {
			sender := defs[0].Description
			rcp := []string{
				user.Phone,
			}

			sms := &messaging.Message{
				Body:       fmt.Sprintf(tpl, bookingID),
				Sender:     sender,
				Recipients: rcp,
			}

			app.msgGW.Send(sms)
		} else {
			app.errorLog.Printf("Failed to find sms template on payment success [booking_id: %v] [err: %v]\n", bookingID, err)
		}
	} else {
		app.errorLog.Printf("Failed to find sms template on payment success [booking_id: %v] [err: %v]\n", bookingID, err)
	}

	return nil
}
