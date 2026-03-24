package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
)

var (
	text      string
	folder    = "audio"
	language  = "en"
	proxy     = ""
	output    string
	delimiter = "|||||"
)

func concat(file string, out *os.File) {
	filePath := fmt.Sprintf("%s/%s.mp3", folder, file)
	defer func() { _ = os.RemoveAll(filePath) }()
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = io.Copy(out, f)
}

func main() {
	for _, value := range os.Args {
		arg := strings.Split(value, "=")
		switch arg[0] {
		case "text":
			text = arg[1]
			break
		case "text-file":
			if data, err := os.ReadFile(arg[1]); err == nil {
				text = string(data)
			}
			break
		case "folder":
			folder = arg[1]
			break
		case "language":
			language = arg[1]
			break
		case "proxy":
			proxy = arg[1]
			break
		case "output":
			output = arg[1]
			break
		case "delimiter":
			delimiter = arg[1]
			break
		case "-h", "--help":
			fmt.Println()
			fmt.Println(`=============================================`)
			fmt.Println("Usage:")
			fmt.Println()
			fmt.Println(`(required): text-file="path_to/text_file.txt"`)
			fmt.Println(`        OR: text="Text to speak"`)
			fmt.Println()
			fmt.Println(`(required): folder="audio"`)
			fmt.Println(`(required): language="en"`)
			fmt.Println(`(required): output="path_to/audio_file.mp3"`)
			fmt.Println()
			fmt.Println(`(option): delimiter="|||||"`)
			fmt.Println(`(option): proxy="https://my-domain.com")`)
			fmt.Println(`=============================================`)
			fmt.Println("Example:")
			fmt.Println()
			fmt.Println(
				` go run main.go folder="audio" language="vi" text-file="content.txt" output="audio/content.mp3"`,
			)
			fmt.Println(
				` go run main.go folder="audio" language="vi" text="Xin chào" output="audio/content.mp3"`,
			)
			fmt.Println(
				`===============================================================================================`,
			)
			fmt.Println()
			return
		default:
			break
		}
	}

	if text == "" {
		return
	}
	if folder != "" {
		_ = os.MkdirAll(folder, os.ModePerm)
	}
	if output == "" {
		output = fmt.Sprintf("%s/output_%d.mp3", folder, time.Now().Unix())
	}
	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer func() { _ = out.Close() }()

	speech := htgotts.Speech{
		Folder:   folder,
		Language: language,
		Proxy:    proxy,
		Handler:  &handlers.Native{},
	}

	newText := text
	spaces := regexp.MustCompile(`(\.\s)|(\.\n)|(\.\t)|(\n+)|(\t+)`).FindAllString(newText, -1)
	for _, space := range spaces {
		newText = strings.Replace(newText, space, delimiter, -1)
	}
	lines := strings.Split(newText, delimiter)
	for _, line := range lines {
		play := strings.TrimSpace(line)
		play = strings.Join(
			regexp.MustCompile(`([\w-.]+@([\w-]+\.)+[\w-]{2,4})|([\p{L}\p{M}\p{N},!?+=\S]+)|([0-9]+)`).
				FindAllString(play, -1),
			" ",
		)
		play = strings.Join(
			regexp.MustCompile(`([\w-.]+@([\w-]+\.)+[\w-]{2,4})|([\p{L}\p{M}\p{N},!?+=\s]+)|([0-9]+)`).
				FindAllString(play, -1),
			"",
		)
		if len(play) == 0 {
			continue
		}
		if strings.TrimSpace(play) == "" {
			continue
		}
		if len(play) > 70 {
			for _, chunk := range regexp.MustCompile(`([\w-.]+@([\w-]+\.)+[\w-]{2,4})|([\p{L}\p{M}\p{N}+=\s]+)|([0-9]+)`).FindAllString(play, -1) {
				if len(chunk) > 70 {
					words := strings.Split(chunk, " ")
					i := 0
					for i < len(words) {
						j := i + 20
						if j > len(words) {
							j = len(words)
						}
						tmpChunk := strings.Join(words[i:j], " ")
						fileName := fmt.Sprintf("%s_%d", language, time.Now().UnixNano())
						_, _ = speech.CreateSpeechFile(tmpChunk, fileName)
						concat(fileName, out)
						i = j
					}
				} else {
					fileName := fmt.Sprintf("%s_%d", language, time.Now().UnixNano())
					_, _ = speech.CreateSpeechFile(chunk, fileName)
					concat(fileName, out)
				}
			}
		} else {
			fileName := fmt.Sprintf("%s_%d", language, time.Now().UnixNano())
			_, _ = speech.CreateSpeechFile(play, fileName)
			concat(fileName, out)
		}
	}
}

// go run main.go folder="audio" language="vi" text-file="content.txt" output="audio/content.mp3"
// go run main.go folder="audio" language="vi" text="Xin chào, đây là bản thử nghiệm chuyển văn bản thành giọng nói tiếng Việt." output="audio/content.mp3"
// git update-ref -d HEAD;git add .;git commit -m "init";git tag -d v0.0.1;git tag v0.0.1;git push -f;git push -f --tags
