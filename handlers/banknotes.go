package handlers

import (
    "fmt"
    "geo-russia-bot/models"
    "gopkg.in/telebot.v3"
    "strings"
)

// ShowAllBanknotes - показать все купюры
func ShowAllBanknotes(c telebot.Context) error {
    text := "🏦 *Купюры России*\n\nВыбери номинал:"
    
    inline := &telebot.ReplyMarkup{}
    var rows []telebot.Row
    
    nominals := []string{"5", "10", "50", "100", "200", "500", "1000", "2000", "5000"}
    var row telebot.Row
    
    for i, nominal := range nominals {
        // Важно: создаем кнопку с уникальным callback data
        btn := inline.Data(nominal+" ₽", "banknote_"+nominal)
        row = append(row, btn)
        
        // По 3 кнопки в ряд
        if (i+1)%3 == 0 || i == len(nominals)-1 {
            rows = append(rows, inline.Row(row...))
            row = telebot.Row{}
        }
    }
    
    // Кнопка возврата в меню
    menuBtn := inline.Data("🏠 Главное меню", "main_menu")
    rows = append(rows, inline.Row(menuBtn))
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}

// ShowBanknoteDetails - показать детали купюры
func ShowBanknoteDetails(c telebot.Context, data string) error {
    nominal := strings.TrimPrefix(data, "banknote_")
    
    // Ищем купюру по номиналу
    var banknote *models.Banknote
    for i, bn := range models.Banknotes {
        if strings.TrimSuffix(bn.Nominal, " рублей") == nominal {
            banknote = &models.Banknotes[i]
            break
        }
    }
    
    if banknote == nil {
        return c.Send("Купюра не найдена")
    }
    
    text := fmt.Sprintf("💵 *%s*\n\n", banknote.Nominal) +
        fmt.Sprintf("🏙️ *Город:* %s\n", banknote.City) +
        fmt.Sprintf("🏛️ *Что изображено:* %s\n\n", banknote.Description) +
        "📚 *Интересные факты:*\n"
    
    for i, fact := range banknote.Facts {
        text += fmt.Sprintf("%d. %s\n", i+1, fact)
    }
    
    inline := &telebot.ReplyMarkup{}
    btnMenu := inline.Data("🏠 Меню", "main_menu")
    btnRandom := inline.Data("🎲 Случайная", "random_banknote")
    inline.Inline(inline.Row(btnMenu, btnRandom))
    
    // Используем Edit, так как это ответ на callback
    return c.Edit(text, inline, telebot.ModeMarkdown)
}

// ShowRandomBanknote - показать случайную купюру
func ShowRandomBanknote(c telebot.Context) error {
    banknote := models.GetRandomBanknote()
    
    text := fmt.Sprintf("🎲 *Случайная купюра*\n\n") +
        fmt.Sprintf("💵 *%s*\n", banknote.Nominal) +
        fmt.Sprintf("🏙️ *Город:* %s\n", banknote.City) +
        fmt.Sprintf("🏛️ *Что изображено:* %s\n\n", banknote.Description) +
        "📚 *Факт:* " + banknote.Facts[0]
    
    inline := &telebot.ReplyMarkup{}
    btnMore := inline.Data("🎲 Ещё раз", "random_banknote")
    btnMenu := inline.Data("🏠 Меню", "main_menu")
    inline.Inline(inline.Row(btnMore), inline.Row(btnMenu))
    
    return c.Edit(text, inline, telebot.ModeMarkdown)
}