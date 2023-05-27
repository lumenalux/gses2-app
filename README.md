# **gses2-app BTC to UAH exchange API**

![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

This is an API that provides the current exchange rate between Bitcoin and the Ukrainian Hryvnia (UAH). It allows users to subscribe to rate updates and receive those updates via email.

## Installation

1. Clone the repository to your desired location:

   ```bash
   git clone https://github.com/lumenalux/gses2-app.git gses2-api
   cd gses2-api
   ```

2. Build the Docker image:

   ```bash
   docker build --tag gses2-app .
   ```

## Usage

1. Run the Docker container:

   ```bash
   docker run -p 8080:8080 gses2-app
   ```

2. Use the API:

   - Get the current BTC to UAH rate:

     ```bash
     curl localhost:8080/api/rate
     ```

   - Subscribe to rate updates:

     ```bash
     curl -X POST -d "email=subscriber@email.com" localhost:8080/api/subscribe
     ```

   - Send rate updates to all subscribers:

     ```bash
     curl -X POST localhost:8080/api/sendEmails
     ```

## Detailed API Usage

For detailed examples of how the API works including screenshots, please see [API_USAGE.md](./API_USAGE.md).

## App settings

The application uses a `config.yaml` file for SMTP server and email configuration. Update the settings in `config.yaml` to use your own SMTP server for sending email updates.

## Configuring the Email Template

The `config.yaml` file contains configuration for the SMTP server as well as the content of the email notifications sent to subscribers. This includes a template for the body of the email, which uses Go's text/template syntax. The application replaces `{{.Rate}}` with the current BTC to UAH exchange rate before sending the email.

Here's what the `config.yaml` looks like:

```yaml
smtp:
  host: smpt-server.example.com
  port: 465
  user: <user>
  password: <password>
email:
  from: no.reply@currency.info.api
  subject: BTC to UAH exchange rate
  body: The BTC to UAH exchange rate is {{.Rate}} UAH per BTC
```

In the `email` section:

- `from`: This field specifies the email address that will appear as the sender of the email.
- `subject`: This field contains the subject line of the email.
- `body`: This field contains the body of the email. Any instance of `{{.Rate}}` in this field will be replaced with the current BTC to UAH exchange rate when the email is sent.

If you wish to modify the content of the email, simply edit the `subject` and/or `body` fields as desired. Remember to rebuild your Docker image to apply the new settings after making changes to `config.yaml`.

!Note: It's important to keep the `{{.Rate}}` placeholder in the `body` field if you want to include the current exchange rate in the email.
