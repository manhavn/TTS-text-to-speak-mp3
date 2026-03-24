# TTS text to speak `>>` mp3

```shell
 go mod tidy
 go mod vendor

 go run main.go folder="audio" language="vi" text-file="content.txt" output="audio/content.mp3"
 # OR
 go run main.go folder="audio" language="vi" text="Xin chào, đây là bản thử nghiệm chuyển văn bản thành giọng nói tiếng Việt." output="audio/content.mp3"
```
