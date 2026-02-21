package keyboards

import (
    "fmt"
    "gopkg.in/telebot.v3"
)

// MainMenu - главное меню
func MainMenu() *telebot.ReplyMarkup {
    menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
    
    btn1 := menu.Text("🏦 Все купюры")
    btn2 := menu.Text("❓ Викторина")
    btn3 := menu.Text("🎮 Игры")
    btn4 := menu.Text("ℹ️ О проекте")
    
    menu.Reply(
        menu.Row(btn1, btn2),
        menu.Row(btn3, btn4),
    )
    
    return menu
}

// GamesMenu - меню игр
func GamesMenu() *telebot.ReplyMarkup {
    menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
    
    btn1 := menu.Text("🎯 Угадай город")
    btn2 := menu.Text("🧩 Найди пару")
    btn3 := menu.Text("🔙 Назад")
    
    menu.Reply(
        menu.Row(btn1, btn2),
        menu.Row(btn3),
    )
    
    return menu
}

// InlineQuizOptions - inline клавиатура для вариантов ответов в викторине
func InlineQuizOptions(questionIndex int, options []string) *telebot.ReplyMarkup {
    inline := &telebot.ReplyMarkup{}
    
    var rows []telebot.Row
    
    for i, option := range options {
        // Создаем уникальный callback data для каждого варианта
        data := fmt.Sprintf("quiz_%d_%d", questionIndex, i)
        btn := inline.Data(option, data)
        rows = append(rows, inline.Row(btn))
    }
    
    inline.Inline(rows...)
    return inline
}

// InlineGameOptions - inline клавиатура для игр
func InlineGameOptions(gameType string) *telebot.ReplyMarkup {
    inline := &telebot.ReplyMarkup{}
    
    switch gameType {
    case "match_cities":
        // Для игры "Найди пару"
        rows := []telebot.Row{
            inline.Row(
                inline.Data("Москва", "match_100"),
                inline.Data("СПб", "match_50"),
            ),
            inline.Row(
                inline.Data("Красноярск", "match_10"),
                inline.Data("Новгород", "match_5"),
            ),
            inline.Row(
                inline.Data("Хабаровск", "match_5000"),
                inline.Data("Ярославль", "match_1000"),
            ),
            inline.Row(
                inline.Data("Севастополь", "match_200"),
                inline.Data("Архангельск", "match_500"),
            ),
        }
        inline.Inline(rows...)
    }
    
    return inline
}

// InlineNavigation - навигация (дальше, назад, меню)
func InlineNavigation() *telebot.ReplyMarkup {
    inline := &telebot.ReplyMarkup{}
    
    btn1 := inline.Data("⬅️ Назад", "nav_back")
    btn2 := inline.Data("🏠 Меню", "nav_menu")
    btn3 := inline.Data("➡️ Далее", "nav_next")
    
    inline.Inline(
        inline.Row(btn1, btn2, btn3),
    )
    
    return inline
}