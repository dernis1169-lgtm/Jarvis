        package main
        import (
            "context"
            "net/http"
            "os"
            "github.com/gin-gonic/gin"
            "github.com/sashabaranov/go-openai"
        )
        func main() {
            r := gin.Default()
            apiKey := os.Getenv("OPENAI_API_KEY")
            r.POST("/query", func(c *gin.Context) {
                var req struct{ Query string `json:"query"` }
                if err := c.ShouldBindJSON(&req); err != nil {
                    c.JSON(400, gin.H{"error": err.Error()})
                    return
                }
                client := openai.NewClient(apiKey)
                resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
                    Model: openai.GPT3Dot5Turbo,
                    Messages: []openai.ChatCompletionMessage{
                        {Role: "system", Content: "Ты Джарвис, ИИ Тони Старка. Отвечай на русском, кратко, называй пользователя сэр."},
                        {Role: "user", Content: req.Query},
                    },
                })
                if err != nil {
                    c.JSON(500, gin.H{"response": "Ошибка ИИ: " + err.Error()})
                    return
                }
                c.JSON(200, gin.H{"response": resp.Choices[0].Message.Content})
            })
            port := os.Getenv("PORT")
            if port == "" { port = "8080" }
            r.Run(":" + port)
        }
        
