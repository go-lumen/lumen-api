# Base API

[![Go Report Card](https://goreportcard.com/badge/github.com/IoThingsDev/api)](https://goreportcard.com/report/github.com/adrien3d/base-api)
[![GoDoc](https://godoc.org/github.com/adrien3d/base-api?status.svg)](https://godoc.org/github.com/adrien3d/base-api)


Base API is an open-source, fast and scalable solution that enables you to speed up your project development, by defining standard common features.

Base API rely on GoLang Gin Gonic web framework, MongoDB and AWS SES for mail management.

## Getting started
### Generate API Keys
If you want to send mails (for user account management) base-api uses AWS SES, so you should get an [API key](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/get-aws-keys.html).

Create a `.env.prod` file from the included `.env.example` file template, while customizing data such as domain name, API keys...

Install `docker` and `docker-compose`

Run `docker-compose up -d`

Watch `yourip:4000`, you should have a welcome message saying `Welcome on base API`.

Congratulations, you are all set !

## NGinx configuration
### With HTTPS using certbot
You can copy paste and customize the [nginx/conf-https-step-1](https://github.com/adrien3d/base-api/tree/master/nginx/conf-https-step-1) to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload
Run certbot ....

Copy paste and customize the [nginx/conf-https-step-2](https://github.com/adrien3d/base-api/tree/master/nginx/conf-https-step-2)
to your `etc/nginx/sites-enabled/yourdomain`
`sudo service nginx reload`

### Without HTTPS

Copy paste and customize the conf-http
to your `etc/nginx/sites-enabled/yourdomain`
`sudo service nginx reload`


## Roadmap
Some features would be nice to have, such as user roles management, Stripe billing management, Twilio SMS alerts.... And may be implemented in the future.

## Miscellaneous
If you want something you consider relevant to be implemented, feel free to fork the repo, and create a PR.