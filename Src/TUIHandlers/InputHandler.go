package tuihandlers

import (
	logger "WillSmith/Logger"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const _INPUT_CHAN_LEN int = 5
const _ECHO_CHANNEL_LEN int = 128

func inputHandler(outpChan *chan string, echoChan *chan byte) {
	var builder = strings.Builder{}
	var isComand = false
	var reader = bufio.NewReader(os.Stdin)
	// Reenabling echoing
	defer exec.Command("stty", "-F", "/dev/tty", "sane").Run()
	for {
		var r,  err = reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				if isComand {
					*outpChan <- builder.String()
					builder.Reset()
				}
				continue
			}
			fmt.Println(err.Error())
		}

		// logger.SendInfo(fmt.Sprintf("Byte: %d", r));

		if !isComand {
			switch r {
			case ':':
				isComand = true
			case 27: // Escape codes like arrow keys
				var codeAlternative, codeErr = readEscapeCodes(reader)
				if codeErr != nil {
					logger.SendWarning(fmt.Sprintf("Failed to read escape code: %v", codeErr.Error()))
				}
				*outpChan <- codeAlternative
				continue
			default:
				*outpChan <- string(r)
				continue
			}
		} else if r == '\n' {
			*outpChan <- builder.String()
			builder.Reset()
			isComand = false
			continue
		}


		switch r {
		case 127:
			logger.SendWarning("flag1")
			if builder.Len() > 0 {
				var temp = builder.String()
				builder.Reset()
				builder.WriteString(temp[:len(temp)-1])
				*echoChan <- '\b'
				*echoChan <- ' '
				*echoChan <- '\b'
			}

		default:
			*echoChan <- r
			builder.WriteByte(r)
		}
	}
}

// By default this function reads 2 bytes from stdio
func readEscapeCodes(reader *bufio.Reader) (string, error) {
	// By this point we read 0d27 from stdio
	var builder = strings.Builder{}
	for range 2 {
		var r,  err = reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return builder.String(), nil
			}
			return "", err
		}
		builder.WriteByte(r)
	}
	return builder.String(), nil
}

func CreateInputHandler() (chan string, chan byte) {
	var outp = make(chan string, _INPUT_CHAN_LEN)
	var echo = make(chan byte, _ECHO_CHANNEL_LEN)
	// Disabling terminal echo
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	go inputHandler(&outp, &echo)
	return outp, echo
}
