<p align="center">
    <a href="https://github.com/adrien3d/lumen-api">
        <img width="500px" src="https://raw.githubusercontent.com/adrien3d/lumen-api/master/lumen-logo.png" />
    </a>
</p>

<h1 align="center">Lumen Api</h1>

<p align="center">
    <a href="https://github.com/adrien3d/lumen-api/blob/master/LICENSE.md">
        <img alt="Go Report Card" src="https://img.shields.io/github/license/adrien3d/lumen-api.svg">
    </a>
    <a href="https://goreportcard.com/report/github.com/adrien3d/lumen-api">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/adrien3d/lumen-api">
    </a>
    <a href="https://godoc.org/github.com/adrien3d/lumen-api">
        <img alt="GoDoc" src="https://godoc.org/github.com/adrien3d/lumen-api?status.svg">
    </a>
    <a href="https://hub.docker.com/r/adrien3d/lumen-api">
        <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/adrien3d/lumen-api.svg">
    </a>
</p>


Lumen API is an open-source, fast and scalable solution that enables you to speed up your project development, by defining standard common features.

Lumen API rely on GoLang Gin Gonic web framework, MongoDB and AWS SES for mail management.

## Getting started
### Generate API Keys
If you want to send mails (for user account management) lumen-api uses AWS SES, so you should get an [API key](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/get-aws-keys.html).

Create a `.env.prod` file from the included `.env.example` file template, while customizing data such as domain name, API keys...

Install `docker` and `docker-compose`

Run `docker-compose up -d`

Watch `yourip:4000`, you should have a welcome message saying `Welcome on lumen API`.

Congratulations, you are all set !

## NGinx configuration
### With HTTPS using certbot
You can copy paste and customize the [nginx/conf-https-step-1](https://github.com/adrien3d/lumen-api/tree/master/nginx/conf-https-step-1) to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload
Run certbot ....

Copy paste and customize the [nginx/conf-https-step-2](https://github.com/adrien3d/lumen-api/tree/master/nginx/conf-https-step-2)
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

## Built With

* [gin-gonic/gin](github.com/gin-gonic/gin) - Gin is a HTTP web framework written in Go (Golang).
* [mgo](github.com/globalsign/mgo) - The MongoDB driver for Go.
* [hermes](github.com/matcornic/hermes) - Golang package that generates clean, responsive HTML e-mails for sending transactional mail.
* [viper](github.com/spf13/viper) - Go configuration with fangs.
* [logrus](github.com/sirupsen/logrus) - Structured, pluggable logging for Go.
* [aws/aws-sdk-go](github.com/aws/aws-sdk-go) - AWS SDK for the Go programming language.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
* **Romain Braems** - *Initial work* - [IoThings](https://github.com/rb62680)
* **Adrien Chapelet** - *Initial work & Updates* - [IoThings](https://github.com/adrien3d)
* **Maxence Henneron** - *Initial work* - [IoThings](https://github.com/maxencehenneron)

