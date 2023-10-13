package main

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/mum4k/termdash/widgets/text"
)

type Display struct {
	report_chan chan *Report
	alert_chan  chan *string
}

func (d *Display) displayReport(report *Report) {
	fmt.Println("Print report")
	fmt.Println("section", report.section)
	fmt.Println("hits", report.number_of_hits)
	fmt.Println("ip from", report.ip_from)
	fmt.Println("hits from this ip", report.hits_from_max_ip)

}

// quotations are used as text that is rolled up in a text widget.
var quotations = []string{
	"When some see coincidence, I see consequence. When others see chance, I see cost.",
	"You cannot pass....I am a servant of the Secret Fire, wielder of the flame of Anor. You cannot pass. The dark fire will not avail you, flame of Udûn. Go back to the Shadow! You cannot pass.",
	"I'm going to make him an offer he can't refuse.",
	"May the Force be with you.",
	"The stuff that dreams are made of.",
	"There's no place like home.",
	"Show me the money!",
	"I want to be alone.",
	"I'll be back.",
}

func (d *Display) listen(ctx context.Context, report_text *text.Text, alert_text *text.Text) {
	for {
		select {
		case report := <-d.report_chan:
			if err := report_text.Write(fmt.Sprintf("%s\n", "New Report")); err != nil {
				panic(err)
			}

			if err := report_text.Write(fmt.Sprintf("%s\n", "section with the most hits - "+report.section)); err != nil {
				panic(err)
			}
			if err := report_text.Write(fmt.Sprintf("%s%d\n", "number of hits  - ", report.number_of_hits)); err != nil {
				panic(err)
			}
		case alert := <-d.alert_chan:
			if err := alert_text.Write(fmt.Sprintf("%s\n", *alert)); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (d *Display) debug_listen() {
	for {
		select {
		case report := <-d.report_chan:
			fmt.Printf("%s\n", "New Report")
			fmt.Printf("%s\n", "section with the most hits - "+report.section)
			fmt.Printf("%s%d\n", "number of hits  - ", report.number_of_hits)
		case m := <-d.alert_chan:
			fmt.Println(*m)

		}
	}
}
func (d *Display) debug_display() {
	go d.debug_listen()
	for {
	}
}

func (d *Display) display() {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	rolled_report, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	rolled_alert, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	go d.listen(ctx, rolled_report, rolled_alert)

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			container.Left(
				container.Border(linestyle.Light),
				container.BorderTitle("Statistic Reports"),
				container.PlaceWidget(rolled_report),
			),
			container.Right(
				container.Border(linestyle.Light),
				container.BorderTitle("Alerts"),
				container.PlaceWidget(rolled_alert),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}

}
