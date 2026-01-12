### Xendit Basic Usage

Comprehensive suite of APIs and no code tools, you can seamlessly accept payments, send payouts, and manage your
business finances with ease.

### Card Processing Usage

- [Simulation for xendit.js](https://js.xendit.co/test_collect_card_data.html)
- [Card collection with .js documentation](https://docs.xendit.co/docs/cards-collecting-card-information)
- [Card Save](https://docs.xendit.co/docs/cards-save-a-card)

1. Create a session
    - [Refer to official documentation](https://docs.xendit.co/apidocs/create-session)
    - Endpoint/URL: https://api.xendit.co/sessions
   ```JSON
   {
      "reference_id": "your-reference-id",
      "session_type": "PAY",
      "mode": "CARDS_SESSION_JS",
      "amount": 100,
      "currency": "PHP",
      "country": "PH",
      "customer": {
        "reference_id": "your-reference-id",
        "type": "INDIVIDUAL",
        "email": "test@yourdomain.com",
        "mobile_number": "+63000000000",
        "individual_detail": {
          "given_names": "Lorem",
          "surname": "Ipsum"
        }
      },
      "cards_session_js": {
        "success_return_url": "https://yourcompany.com/success",
        "failure_return_url": "https://yourcompany.com/failure"
      }
   }
   ```
2. Create a custom form for xendit.js
    - [xendit.js](https://js.xendit.co/cards-session.min.js)
    - [Guest checkout](https://docs.xendit.co/docs/cards-guest-checkout-one-off-payment)
    - Form Details:
        - Xendit.setPublishableKey("your-public-key");
        - card_number
        - expiry_month
        - expiry_year
        - cvn
        - cardholder_first_name
        - cardholder_last_name
        - cardholder_email
        - cardholder_phone_number
        - payment_session_id {returned from create session api}
    - Usage:
    ```JAVASCRIPT
    Xendit.payment.collectCardData(formData, function (err, response) {
        console.log(err, response);
    });
    ```
3. Payment API Webhook
    - [Payments API Webhook](https://docs.xendit.co/docs/payments-api-webhooks)