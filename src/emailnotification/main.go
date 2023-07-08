package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/emailnotification/model"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/emailnotification/service"
)

func main() {
	nc, err := service.ConnectToNats()
	if err != nil {
		log.Fatalf("unable to connect to nats %v", err)
	}
	nc.Subscribe("confirmOrder.*", func(m *nats.Msg) {
		var (
			req model.EmialContent
			res model.Response
		)
		if err := json.Unmarshal(m.Data, &req); err != nil {
			panic(err)
		}

		emailadre, emailcontent, emailtype := service.CreateEmail(req, "confirm")
		err := service.SendEmail(emailadre, emailcontent, emailtype)
		if err != nil {
			res.Send = false
		} else {
			res.Send = true
		}
		e, errMarshal := json.Marshal(res)
		if errMarshal != nil {
			fmt.Println(errMarshal)
			return
		}
		nc.Publish(m.Reply, []byte(e))
	})

	nc.Subscribe("confirmCancel.*", func(m *nats.Msg) {
		var (
			req model.EmialContent
			res model.Response
		)
		if err := json.Unmarshal(m.Data, &req); err != nil {
			panic(err)
		}
		emailadre, emailcontent, emailtype := service.CreateEmail(req, "cancel")
		if err := service.SendEmail(emailadre, emailcontent, emailtype); err != nil {
			res.Send = false
		} else {
			res.Send = true
		}
		e, errMarshal := json.Marshal(res)
		if errMarshal != nil {
			fmt.Println(errMarshal)
			return
		}
		nc.Publish(m.Reply, []byte(e))
	})

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}

}
