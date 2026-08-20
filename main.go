package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// Структуры для общения с Android-приложением
type JarvisRequest struct {
	Query string `json:"query"`
}

type JarvisResponse struct {
	Response string `json:"response"`
}

const (
	OpenAIKey = "YOUR_OPENAI_API_KEY" // Вставьте ваш ключ сюда
)

func main() {
	r := gin.Default()

	r.POST("/query", func(c *gin.Context) {
		var req JarvisRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Сэр, получен запрос: %s\n", req.Query)

		// 1. Логика "поиска" в интернете (упрощенная)
		// Здесь можно подключить API типа Google Search или Serper
		searchContext := ""
		if strings.Contains(strings.ToLower(req.Query), "найди") || strings.Contains(strings.ToLower(req.Query), "новости") {
			searchContext = " (Контекст: сегодня 20 августа 2026 года, в мире технологий бум квантовых вычислений)"
		}

		// 2. Обращение к ИИ-мозгу (OpenAI)
		client := openai.NewClient(OpenAIKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT4,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "Ты — Джарвис, ИИ Тони Старка. Будь дерзким, называй пользователя 'сэр'. Ты самый умный.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: req.Query + searchContext,
					},
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Сэр, произошел сбой в нейронных цепях"})
			return
		}

		aiResponse := resp.Choices[0].Message.Content
		c.JSON(http.StatusOK, JarvisResponse{Response: aiResponse})
	})

	// Запуск на всех интерфейсах для доступа с телефона
	fmt.Println("Джарвис онлайн на порту 8000...")
	r.Run("0.0.0.0:8000")
}
