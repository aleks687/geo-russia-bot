package handlers

import (
    "fmt"
    "geo-russia-bot/models"
    "gopkg.in/telebot.v3"
    "strconv"
    "strings"
)

// StartQuiz - начало викторины
func StartQuiz(c telebot.Context) error {
    userID := c.Sender().ID
    
    questions := models.GetRandomQuestions(5)
    models.GameStorage[userID] = &models.UserState{
        CurrentGame:   "quiz",
        Score:         0,
        QuestionIndex: 0,
    }
    
    return sendQuizQuestion(c, userID, questions[0], 0, 5)
}

// sendQuizQuestion - отправить вопрос
func sendQuizQuestion(c telebot.Context, userID int64, question models.Question, index, total int) error {
    text := fmt.Sprintf("❓ *Вопрос %d/%d*\n\n%s", index+1, total, question.Text)
    
    inline := &telebot.ReplyMarkup{}
    var rows []telebot.Row
    
    // Создаем кнопки для каждого варианта ответа
    for i, option := range question.Options {
        data := fmt.Sprintf("quiz_%d_%d", index, i)
        btn := inline.Data(option, data)
        rows = append(rows, inline.Row(btn))
    }
    
    // Кнопка выхода
    exitBtn := inline.Data("🚪 Выйти", "quiz_exit")
    rows = append(rows, inline.Row(exitBtn))
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}

// HandleQuizAnswer - обработка ответа
func HandleQuizAnswer(c telebot.Context, data string) error {
    // Парсим data вида "quiz_0_2"
    parts := strings.Split(data, "_")
    if len(parts) != 3 {
        return nil
    }
    
    questionIdx, _ := strconv.Atoi(parts[1])
    answerIdx, _ := strconv.Atoi(parts[2])
    
    userID := c.Sender().ID
    state, exists := models.GameStorage[userID]
    
    if !exists || state.CurrentGame != "quiz" {
        return c.Send("Викторина не найдена. Начните новую с /quiz")
    }
    
    // Получаем вопросы (в реальном приложении нужно сохранять вопросы в состоянии)
    questions := models.GetRandomQuestions(5)
    if questionIdx >= len(questions) {
        return nil
    }
    
    question := questions[questionIdx]
    
    // Формируем ответ
    var responseText string
    if answerIdx == question.Correct {
        state.Score++
        responseText = "✅ *Правильно!*\n\n"
    } else {
        responseText = fmt.Sprintf("❌ *Неверно!*\nПравильный ответ: *%s*\n\n", 
            question.Options[question.Correct])
    }
    
    responseText += "📚 *Факт:* " + question.Fact
    
    // Отправляем результат (новым сообщением)
    c.Send(responseText, telebot.ModeMarkdown)
    
    // Проверяем, есть ли еще вопросы
    if questionIdx+1 >= len(questions) {
        return finishQuiz(c, userID, state.Score, len(questions))
    }
    
    // Отправляем следующий вопрос
    return sendQuizQuestion(c, userID, questions[questionIdx+1], questionIdx+1, len(questions))
}

// finishQuiz - завершение викторины
func finishQuiz(c telebot.Context, userID int64, score, total int) error {
    delete(models.GameStorage, userID)
    
    percentage := float64(score) / float64(total) * 100
    
    var emoji, comment string
    switch {
    case percentage == 100:
        emoji, comment = "🏆", "Потрясающе! Ты знаток!"
    case percentage >= 80:
        emoji, comment = "🌟", "Отлично!"
    case percentage >= 60:
        emoji, comment = "👍", "Хорошо!"
    default:
        emoji, comment = "📚", "Попробуй ещё раз!"
    }
    
    text := fmt.Sprintf("%s *Викторина завершена!*\n\n", emoji) +
        fmt.Sprintf("Результат: *%d/%d* (%.0f%%)\n\n", score, total, percentage) +
        fmt.Sprintf("%s\n\nХочешь ещё?", comment)
    
    inline := &telebot.ReplyMarkup{}
    btnAgain := inline.Data("🔄 Ещё раз", "quiz_again")
    btnMenu := inline.Data("🏠 Меню", "main_menu")
    inline.Inline(inline.Row(btnAgain), inline.Row(btnMenu))
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}