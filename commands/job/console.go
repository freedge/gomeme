package job

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/freedge/gomeme/commands"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	tb "github.com/nsf/termbox-go"
)

type jobConsole struct {
	jobsStatusCommonCommand
}

func (cmd *jobConsole) Data() interface{} {
	return nil
}

func getColoredStatus(s string) string {
	switch s {
	case "Ended Not OK":
		return "red"
	case "Ended OK":
		return "green"
	case "Executing":
		return "blue"
	}
	return "white"
}

var p *widgets.Paragraph

// read a line, put it as title of p
func readStuff(events <-chan ui.Event) string {
	defer func() { p.Title = ">" }()
	s := ""
	for e := range events {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<Enter>":
				return s
			case "<C-c>", "<Escape>":
				return ""

			case "<C-<Backspace>>", "<Backspace>":
				if len(s) > 0 {
					s = s[:len(s)-1]
					p.Title = p.Title[:len(p.Title)-1]
					ui.Render(p)
				}
			case "<Space>":
				s = s + " "
				p.Title += " "
				ui.Render(p)
			default:
				if len(e.ID) == 1 {
					s = s + e.ID
					p.Title += e.ID
					ui.Render(p)
				}
			}
		}
	}
	return s
}

// display file, return if next file should be displayed
// this is broken under windows, but resizing the screen seems to work
func displayFile(events <-chan ui.Event, s string) bool {
	defer ui.Clear()
	defer ui.TerminalDimensions() // hopefully might clear the scren properly

	fmt.Println(s) // breaks the terminal

	for e := range events {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<Enter>",
				"<C-c>", "<Escape>", "<Space>", "q":
				return false
			}
		}
	}
	return false
}

func (cmd *jobConsole) Execute([]string) (err error) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	tb.SetInputMode(tb.InputEsc) // do not capture mouse events

	if _, err := cmd.GetJobs(); err != nil {
		return err
	}

	l := widgets.NewList()
	l.Title = "Statuses"
	l.Rows = make([]string, 0, len(cmd.reply.Statuses))

	for _, job := range cmd.reply.Statuses {
		s := fmt.Sprintf("[%-40.40s %5.5s %-14.14s %8.8s %14.14s](fg:%s) %17.17s %5.5s %12.12s %21.21s %8.8s % 4d",
			job.Folder+"/"+job.Name,
			strconv.FormatBool(job.Held),
			job.JobId, job.OrderDate, job.Status, getColoredStatus(job.Status), job.Host, strconv.FormatBool(job.Deleted), job.StartTime, job.Description, GetDurationAsString(job), job.NumberOfRuns)
		l.Rows = append(l.Rows, s)
	}

	l.SelectedRowStyle = ui.NewStyle(ui.ColorClear, ui.ColorClear, ui.ModifierBold)
	l.WrapText = false
	l.BorderLeft = false
	l.BorderRight = false
	l.PaddingLeft = -1
	l.PaddingRight = -1
	x, y := ui.TerminalDimensions()
	l.SetRect(0, 0, x, y)
	p = widgets.NewParagraph()
	p.SetRect(30, 0, x, 1)
	p.Title = ""
	p.Border = false
	ui.Render(l)
	var previousKey string
	events := ui.PollEvents()

	for e := range events {
		if e.Type == ui.ResizeEvent {
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x, y)

			ui.Render(l)
		} else if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "q", "<C-c>", "<Escape>":
				return
			case "j", "<Down>":
				l.ScrollDown()
			case "k", "<Up>":
				l.ScrollUp()
			case "<C-d>":
				l.ScrollHalfPageDown()
			case "<C-u>":
				l.ScrollHalfPageUp()
			case "<C-f>", "<PageDown>", "<Space>":
				l.ScrollPageDown()
			case "<C-b>", "<PageUp>":
				l.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					l.ScrollTop()
				}
			case "<Home>":
				l.ScrollTop()
			case "G", "<End>":
				l.ScrollBottom()
			case "/":
				p.Title = "/"
				ui.Render(p)
				tosearch := readStuff(events)
			s:
				for i := l.SelectedRow; i < len(cmd.reply.Statuses); i++ {
					if strings.Contains(l.Rows[i], tosearch) {
						l.SelectedRow = i
						break s
					}
				}
				p.Title = ""
			case "l", "o":
				jlc := jobLogCommand{Jobid: cmd.reply.Statuses[l.SelectedRow].JobId, Output: false}
				err := jlc.Execute([]string{})
				if err == nil {
					ui.Clear()
					displayFile(events, jlc.Data().(string))
				}

			case "<Enter>":

				jlc := jobLogCommand{Jobid: cmd.reply.Statuses[l.SelectedRow].JobId, Output: true}
				err := jlc.Execute([]string{})
				if err == nil {
					ui.Clear()
					displayFile(events, jlc.Data().(string))
				}

			case "s":
				jlc := jobStatusCommand{Jobid: cmd.reply.Statuses[l.SelectedRow].JobId}
				err := jlc.Execute([]string{})
				if err == nil {
					ui.Clear()
					jlc.PrettyPrint()
					displayFile(events, "")
				}
			case "t":
				jlc := jobTreeCommand{jobsStatusCommand: jobsStatusCommand{
					jobsStatusCommonCommand: jobsStatusCommonCommand{
						Jobid:      cmd.reply.Statuses[l.SelectedRow].JobId,
						Neighbours: true,
					}}}

				err := jlc.Execute([]string{})
				if err == nil {
					ui.Clear()
					jlc.PrettyPrint()
					displayFile(events, "")
				}
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			ui.Render(l)
			if p.Title != "" {
				ui.Render(p)
			}
		}
	}

	return nil
}

func (cmd *jobConsole) PrettyPrint() error {
	return nil
}

func init() {
	commands.AddCommand("con", "console", "Renders jobs in a terminal", &jobConsole{})
}
