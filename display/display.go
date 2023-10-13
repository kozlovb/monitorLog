package display

import (
	"monitorLog/stats"

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
	Report_chan chan *stats.Report
	Alert_chan  chan *string
}

func (d *Display) listen(ctx context.Context, report_text *text.Text, alert_text *text.Text) {
	for {
		select {
		case report := <-d.Report_chan:
			if err := report_text.Write(fmt.Sprintf("%s\n", "New Report")); err != nil {
				panic(err)
			}

			if err := report_text.Write(fmt.Sprintf("%s\n", "section with the most hits - "+report.Section)); err != nil {
				panic(err)
			}
			if err := report_text.Write(fmt.Sprintf("%s%d\n", "number of hits  - ", report.Number_of_hits)); err != nil {
				panic(err)
			}
		case alert := <-d.Alert_chan:
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
		case report := <-d.Report_chan:
			fmt.Printf("%s\n", "New Report")
			fmt.Printf("%s\n", "section with the most hits - "+report.Section)
			fmt.Printf("%s%d\n", "number of hits  - ", report.Number_of_hits)
		case m := <-d.Alert_chan:
			fmt.Println(*m)

		}
	}
}
func (d *Display) Debug_display() {
	go d.debug_listen()
	for {
	}
}

func (d *Display) Display() {
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
