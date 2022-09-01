package docker

import (
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/morikuni/aec"

	"encoding/json"
	"fmt"
	"io"
)

// Taken and adapted from https://github.com/moby/moby/blob/master/pkg/jsonmessage/jsonmessage.go
func HandlePullOutput(
	pullOutputReader io.Reader,
	streamHandler func(logLine string) error,
) error {

	var (
		dec        = json.NewDecoder(pullOutputReader)
		ids        = make(map[string]uint)
		out        = newLogger(streamHandler)
		isTerminal = true
	)

	for {
		var diff uint
		var jm jsonmessage.JSONMessage

		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if jm.Aux != nil {
			continue
		}

		if jm.ID != "" && (jm.Progress != nil || jm.ProgressMessage != "") {
			line, ok := ids[jm.ID]

			if !ok {
				// NOTE: This approach of using len(id) to
				// figure out the number of lines of history
				// only works as long as we clear the history
				// when we output something that's not
				// accounted for in the map, such as a line
				// with no ID.
				line = uint(len(ids))
				ids[jm.ID] = line

				if isTerminal {
					fmt.Fprintf(out, "\n")
				}
			}

			diff = uint(len(ids)) - line

			if isTerminal {
				cursorUp(out, diff)
			}
		} else {
			// When outputting something that isn't progress
			// output, clear the history of previous lines. We
			// don't want progress entries from some previous
			// operation to be updated (for example, pull -a
			// with multiple tags).
			ids = make(map[string]uint)
		}

		err := jm.Display(out, isTerminal)

		if jm.ID != "" && isTerminal {
			cursorDown(out, diff)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type logger struct {
	streamHandler func(logLine string) error
}

func newLogger(streamHandler func(logLine string) error) logger {
	return logger{
		streamHandler: streamHandler,
	}
}

func (l logger) Write(p []byte) (n int, err error) {
	err = l.streamHandler(string(p))

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func cursorUp(out io.Writer, l uint) {
	fmt.Fprint(out, aec.Up(l))
}

func cursorDown(out io.Writer, l uint) {
	fmt.Fprint(out, aec.Down(l))
}
