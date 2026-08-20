package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type JarvisRequest struct {
	Query string `json:"query"`
}

type JarvisResponse struct {
	Response string `json:"response"`
}

func main() {
	// Устанавливаем режим релиза для скорости
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Получаем API ключ из настроек сервера (Environment Variables)
	apiKey := os.Getenv("OPENAI_API_KEY")

	r.POST("/query", func(c *gin.Context) {
		var req JarvisRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Если ключа нет, Джарвис вежливо об этом скажет
		if apiKey == "" {
			c.JSON(200, JarvisResponse{Response: "Сэр, мой модуль ИИ отключен. Пожалуйста, добавьте OPENAI_API_KEY в настройки Render."})
			return
		}

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "Ты — Джарвис, ИИ помощник Тони Старка. Отвечай кратко, дерзко, называй пользователя 'сэр'. Используй русский язык.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: req.Query,
					},
				},
			},
		)

		if err != nil {
			c.JSON(500, gin.H{"response": "Сэр, возникла ошибка в нейронных связях: " + err.Error()})
			return
		}

		c.JSON(200, JarvisResponse{Response: resp.Choices[0].Message.Content})
	})

	// Порт для Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
