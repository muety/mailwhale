<p align="center">
    <img src="assets/logo.svg" width="275px">
</p>

<h1 align="center">MailWhale</h1>
<h3 align="center">A <i>bring-your-own-SMTP-server</i> mail relay</h3>

<p align="center">
<img src="https://badges.fw-web.space/github/license/muety/mailwhale">
<a href="https://saythanks.io/to/n1try" target="_blank"><img src="https://badges.fw-web.space/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg"></a>
<a href="https://wakapi.dev" target="_blank"><img src="https://badges.fw-web.space/endpoint?url=https://wakapi.dev/api/compat/shields/v1/n1try/interval:any/project:mailwhale&color=blue"></a>
</p>

## 📄 Description

Being a web developer, chances are high that at some point you need to teach your application how to send mails.
Essentially, there are two options. Either you use a **professional mail sending service**
like [Mailgun](https://mailgun.com), [SendGrid](https://sendgrid.com), [SMTPeter](https://smtpeter.com) and the like or
you **include an SMTP client library** to your software and **plug your own mail server**.

However, if you want the best of both worlds – that is, send mails via simple HTTP calls and with no extra complexity ,
but still use your own infrastructure, you may want to go with ✉️🐳.

You get a simple REST API, which you can call to send out e-mail. Stay tuned, there is a lot more to come.

## 📦 Installation

Assuming you have Go installed, all you need to do is:

```bash
# 1. Clone repo
$ git clone https://github.com/muety/mailwhale.git

# 2. Adapt config to your needs, i.e. set your SMTP server and credentials, etc.
$ cp config.default.yml config.yml
$ vi config.yml

# 3. Run it
$ GO111MODULE=on go build
$ ./mailwhale
```

## ⌨️ Usage

MailWhale has the notion of **clients**, which are applications allowed to access the API to send mails or manage other
clients. Once you start MailWhale for the first time, a default client is created and its credentials are printed to the
console (e.g. `f9142379-0d5e-4692-abb7-d42fba2a2fef`). Remember those, as you will need them to use the API.

### Create new client application

```bash
$ curl -XPOST \
  -u 'root:<your_api_key>' \
  -H 'content-type: application/json' \
  --data '{ "name": "my-cool-app", "permissions": [ "send_mail" ] }' \
  http://localhost:3000/api/client

# Response (201 Created):
# {"name":"my-cool-app-2","permissions":["send_mail"],"api_key":"ce67d653-92ef-46a1-98fd-50feb1c07495"}
```

### Send an HTML mail (synchronously)

```bash
curl -XPOST \
  -u 'root:<your_api_key>' \
  -H 'content-type: application/json' \
  --data '{
      "from": "John Doe <john.doe@example.org>",
      "to": ["Jane Doe <jane@doe.com>"],
      "subject": "Dinner tonight?",
      "html": "<h1>Hey you!</h1><p>Wanna have dinner tonight?</p>"
  }' \
  http://localhost:3000/api/mail
```

You can also a `text` field instead, to send a plain text message.

## 🚀 Features (planned)

Right now, this app is very basic. However, there are several cool features on our roadmap.

* **Mail Templates:** Users will be able to create complex (HTML) templates or presets for their mails, which can then
  be referenced in send requests.
* **Bound handling:** Ultimately, we want to offer the ability to plug an IMAP server in addition, to get notified about
  bounced / undelivered mails.
* **Statistics:** There will be basic statistics about when which client has sent how many mails, how many were
  successful or were bounced, etc.
* **Web UI:** A nice-looking web UI will make client- and template management easier.
* **Minor enhancements:** IPv6- and TLS support, API documentation, ...

## 📓 License

MIT