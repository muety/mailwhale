<p align="center">
    <img src="assets/logo.svg" width="275px">
</p>

<h1 align="center">MailWhale</h1>
<h3 align="center">A <i>bring-your-own-SMTP-server</i> mail relay</h3>

<p align="center">
<img src="https://badges.fw-web.space/github/license/muety/mailwhale">
<a href="mailto:ferdinand@muetsch.io?subject=Thanks for creating MailWhale" target="_blank"><img src="https://badges.fw-web.space/badge/say%20thanks-%F0%9F%99%8F-%23159CE4"></a>
<a href="https://wakapi.dev" target="_blank"><img src="https://badges.fw-web.space/endpoint?url=https://wakapi.dev/api/compat/shields/v1/n1try/interval:any/project:mailwhale&color=blue"></a>
</p>

## 📄 Description

As a developer, chances are high that at some point you need to teach your application how to send mails.
Essentially, there are two options. Either you use a **professional mail sending service** or
you **include an SMTP client library** to your software and **plug your own mail server**.

Think of MailWhale like [Mailgun](https://mailgun.com), [SendGrid](https://sendgrid.com) or [SMTPeter](https://smtpeter.com), but open source and self-hosted. Or like [Postal](https://github.com/postalhq/postal) or [Cuttlefish](https://cuttlefish.io/), but less bloated and without running its own, internal SMTP server. 

However, if you want the best of both worlds – that is, send mails via simple HTTP calls and with no extra complexity, 
but still use your own infrastructure – you may want to go with ✉️🐳.

You get a simple **REST API**, which you can call to send out e-mail. You can plug your self-hosted SMTP server, as well as Google Mail or **literally any other e-mail provider**.

Stay tuned, there is a lot more to come.

![](assets/screenshot01.png)

## 🚧 Project State
The project is in a very early stage and breaking changes are likely to happen. We'd recommend to not yet use this in production or at least expect non-trivial effort required to upgrade to a new version.

## 📦 Installation
### Compile from source
```bash
# 1. Clone repo
$ git clone https://github.com/muety/mailwhale.git

# 2. Adapt config to your needs, i.e. set your SMTP server and credentials, etc.
$ cp config.default.yml config.yml
$ vi config.yml

# 3. Compile Web UI
$ cd webui/
$ yarn && yarn build
$ cd ..
# 4. Compile API
$ GO111MODULE=on go build

# 5. Run it
$ ./mailwhale
```

### With Docker
```bash
# 1. Clone repo
$ git clone https://github.com/muety/mailwhale.git

# 2. Adapt config to your needs, i.e. set your SMTP server and credentials, etc.
$ cp config.default.yml config.yml
$ vi config.yml

# 3. Build image
$ docker build -t mailwhale .

# 4. Create persistent volume
$ docker volume create mailwhale_data

# 5. Run
$ docker run -d \
  -p 127.0.0.1:3000:3000 \
  -v "$(pwd)/config.yml":/app/config.yml:ro \
  -v mailwhale_data:/data \
  --name mailwhale \
  mailwhale
```

**Note:** An official Docker image is about to come. Also, there will be no need to mount your config file into the container, as everything will be configurable using environment variables eventually. 

## ⌨️ Usage
First of all, you can get most tasks done through the web UI, available at http://localhost:3000.

### 1. Define a user
To get started with MailWhale, you need to create a **user** first. To do so, register a new user API or web UI. `security.allow_signup` needs to be set to `true`. 

### 2. Create an API client
It is good practice to not authenticate against the API as a user directly. Instead, create an **API client** with limited privileges, that could easily be revoked in the future. A client is identified by a **client ID** and a **client secret** (or token), very similar to what you might already be familiar with from AWS APIs. Usually, such a client corresponds to an individual client application of yours, which wants to access MailWhale's API. 

#### Request
```bash
$ curl -XPOST \
     -u 'admin@local.host:admin' \
     -H 'Content-Type: application/json' \
     --data-raw '{
         "description": "My juicy web app",
         "sender": "Epic Juice Store <noreply@epicjuicestore.org>",
         "permissions": ["send_mail"]
     }' \
     'http://localhost:3000/api/client'
```

#### Response
```
{
    "id": "SVNORFBUWGhxWGZSUUl0eA==",
    "description": "My juicy web app",
    "permissions": [
        "send_mail"
    ],
    "sender": "Epic Juice Store <noreply@epicjuicestore.org>",
    "api_key": "75c74447-c4af-453b-ad06-3a8ae969ed16"
}
```

The response contains your new client's ID (`id`) and secret (`api_key`). Remember these credentials, as they are needed for subsequent requests from your application.

Client authentication happens through HTTP [basic auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme). Most HTTP clients support basic auth out of the box (including cURL with its `-u` parameter). If your's doesn't, you can hash create the hash like so:

```bash
$ echo "Authorization: Basic $(echo '<client_id>:<client_secret>' | base64)"

# Result:
# Authorization: Basic U1ZOT1JGQlVXR2h4V0daU1VVbDBlQT09Ojc1Yzc0NDQ3LWM0YWYtNDUzYi1hZDA2LTNhOGFlOTY5ZWQxNg==
```

### 3. Send E-Mails

#### Plain text or HTML
```bash
$ curl -XPOST \
  -u '<client_id>:<client_secret>' \
  -H 'content-type: application/json' \
  --data '{
      "to": ["Jane Doe <jane@doe.com>"],
      "subject": "Dinner tonight?",
      "html": "<h1>Hey you!</h1><p>Wanna have dinner tonight?</p>"
  }' \
  'http://localhost:3000/api/mail'
```

You can also a `text` field instead, to send a plain text message.

#### Using a template
In case you have created a template using the web UI, you can reference it in a new mail like so:
```bash
$ curl -XPOST \
  -u '<client_id>:<client_secret>' \
  -H 'content-type: application/json' \
  --data '{
      "to": ["Jane Doe <jane@doe.com>"],
      "subject": "Dinner tonight?",
      "template_id": "8033ea08-2630-408b-82f9-d38b403243d0",
      "template_vars: {
        "text.greeting": "Hello new user!",
    }
  }' \
  'http://localhost:3000/api/mail'
```

## 🔧 Configuration Options
You can specify configuration options either via a config file (`config.yml`) or via environment variables. Here is an overview of all options.

| YAML Key                  | Environment Variable      | Default      | Description                                                         |
|---------------------------|---------------------------|--------------|---------------------------------------------------------------------|
| `env`                           | `MW_ENV`            | `dev`             | Whether to use development- or production settings |
| `mail.domain`                   | `MW_MAIL_DOMAIN`    | -                 | Default domain for sending mails |
| `web.listen_addr`               | `MW_WEB_LISTEN_ADDR`| `127.0.0.1:3000`  | IP and port for the web server to listen on (can be IPv4 or IPv6) |
| `web.cors_origin`               | -                   | [`http://localhost:5000`] | List of URLs which to accept CORS requests for |
| `web.public_url`                | `MW_PUBLIC_URL`     | `http://localhost:3000` | The URL under which your MailWhale server is available from the public internet |
| `smtp.host`                     | `MW_SMTP_HOST`      | -                 | SMTP relay host name or IP |
| `smtp.port`                     | `MW_SMTP_PORT`      | -                 | SMTP relay port |
| `smtp.username`                 | `MW_SMTP_USER`      | -                 | SMTP relay authentication user name |
| `smtp.password`                 | `MW_SMTP_PASS`      | -                 | SMTP relay authentication password |
| `smtp.tls`                      | `MW_SMTP_TLS`       | `false`           | Whether to require full TLS (not to be confused with STARTTLS) for the SMTP relay |
| `store.path`                    | `MW_STORE_PATH`     | `./data.json.db`   | Target location of the database file |
| `security.pepper`               | `MW_SECURITY_PEPPER`| -                 | Pepper to use for hashing user passwords |
| `security.allow_signup`         | `MW_SECURITY_ALLOW_SIGNUP` | `true`    | Whether to allow the registration of new users |
| `security.verify_users`         | `MW_SECURITY_VERIFY_USERS` | `true`    | Whether to require new users to activate their account using a confirmation mail |
| `security.verify_senders`       | `MW_SECURITY_VERIFY_SENDERS` | `true`  | Whether to validate sender addresses and their domains' SPF records |

### Sender verification & SPF Check
By default, mails are sent using a randomly generated address in the `From` header, which belongs to the domain configured via `mail.domain` (i.e. `user+abcdefgh@wakapi.dev`). Optionally, custom sender addresses can be configured on a per-API-client basis. However, it is recommended to properly configure [SPF](https://en.wikipedia.org/wiki/Sender_Policy_Framework) on that custom domain and instruct MailWhale to verify that configuration.

**As a user**, you need to configure your domain, which you want to use as part of your senders address (e.g. `example.org` for sending mails from `User Server <noreply@example.org>`), to publish an SPF record that delegates to the domain under which MailWhale is running (e.g. mailwhale.dev).
```
example.org.  IN  TXT v=spf1 include:mailwhale.dev
```

**As a server operator** of a MailWhale instance, you need to enable `mail.verify_senders` and set your `mail.domain` and `web.public_url`. For that domain, you need to configure an SPF record that allows your SMTP relay provider's (e.g. Mailbox.org, GMail, SendGrid, etc.) mail servers to be senders. Refer to your provider's documentation, e.g. [this](https://kb.mailbox.org/display/MBOKBEN/How+to+integrate+external+e-mail+accounts). 

## 🚀 Features (planned)

Right now, this app is very basic. However, there are several cool features on our roadmap.

* **Bounce handling:** Ultimately, we want to offer the ability to plug an IMAP server in addition, to get notified about
  bounced / undelivered mails.
* **Statistics:** There will be basic statistics about when which client has sent how many mails, how many were
  successful or were bounced, etc.
* **Client libraries:** To make the developer experience even smoother, client SDKs for different programming languages will we added some time.
* **Minor enhancements:** IPv6- and TLS support, API documentation, ...

## 📓 License

MIT
