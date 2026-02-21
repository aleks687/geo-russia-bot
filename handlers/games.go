package handlers

import (
    "fmt"
    "geo-russia-bot/models"
    "gopkg.in/telebot.v3"
    "math/rand"
    "strconv"
    "strings"
)

// Хранилище для игры "Найди пару"
var matchGames = make(map[int64]*MatchGameState)

type MatchGameState struct {
    SelectedCity string
    Score        int
}

// GuessCityGame - игра "Угадай город"
func GuessCityGame(c telebot.Context) error {
    userID := c.Sender().ID
    banknote := models.GetRandomBanknote()
    
    models.GameStorage[userID] = &models.UserState{
        CurrentGame: "guess_city",
        Score:       0,
    }
    
    // Варианты ответов (правильный город + 3 случайных)
    cities := []string{banknote.City}
    for len(cities) < 4 {
        randomBanknote := models.GetRandomBanknote()
        if !contains(cities, randomBanknote.City) && randomBanknote.City != banknote.City {
            cities = append(cities, randomBanknote.City)
        }
    }
    
    // Перемешиваем варианты
    rand.Shuffle(len(cities), func(i, j int) {
        cities[i], cities[j] = cities[j], cities[i]
    })
    
    text := fmt.Sprintf("🎯 *Угадай город*\n\n") +
        fmt.Sprintf("На купюре *%s* изображён:\n", banknote.Nominal) +
        fmt.Sprintf("_%s_\n\n", banknote.Description) +
        "Какой это город?"
    
    inline := &telebot.ReplyMarkup{}
    var rows []telebot.Row
    
    // Находим индекс правильного ответа после перемешивания
    correctIdx := 0
    for i, city := range cities {
        if city == banknote.City {
            correctIdx = i
        }
        // Создаем callback data с правильным индексом
        data := fmt.Sprintf("guess_%d_%d", correctIdx, i)
        btn := inline.Data(city, data)
        rows = append(rows, inline.Row(btn))
    }
    
    exitBtn := inline.Data("🚪 Выйти", "game_exit")
    rows = append(rows, inline.Row(exitBtn))
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}

// HandleGuessAnswer - обработка ответа в игре
func HandleGuessAnswer(c telebot.Context, data string) error {
    // Парсим data вида "guess_2_1"
    parts := strings.Split(data, "_")
    if len(parts) != 3 {
        return nil
    }
    
    correctIdx, _ := strconv.Atoi(parts[1])
    selectedIdx, _ := strconv.Atoi(parts[2])
    
    userID := c.Sender().ID
    state, exists := models.GameStorage[userID]
    
    if !exists || state.CurrentGame != "guess_city" {
        return c.Send("Игра не найдена. Начните новую через меню 'Игры'")
    }
    
    var responseText string
    if correctIdx == selectedIdx {
        state.Score++
        responseText = "✅ *Верно!* Молодец!\n\n"
    } else {
        responseText = "❌ *Не угадал...*\n\n"
    }
    
    // Очищаем состояние игры
    delete(models.GameStorage, userID)
    
    responseText += fmt.Sprintf("Твой счёт: *%d*\n\nХочешь сыграть ещё?", state.Score)
    
    inline := &telebot.ReplyMarkup{}
    btnAgain := inline.Data("🔄 Ещё раз", "guess_again")
    btnMenu := inline.Data("🏠 Меню", "main_menu")
    inline.Inline(inline.Row(btnAgain), inline.Row(btnMenu))
    
    return c.Send(responseText, inline, telebot.ModeMarkdown)
}

// MatchGame - игра "Найди пару" (начало)
func MatchGame(c telebot.Context) error {
    text := "🧩 *Найди пару*\n\n" +
        "Сопоставь город и купюру, на которой он изображён.\n\n" +
        "Нажми на кнопку с городом:"
    
    // Создаем inline клавиатуру с городами
    inline := &telebot.ReplyMarkup{}
    
    rows := []telebot.Row{
        inline.Row(
            inline.Data("🏙️ Москва", "match_city_moscow"),
            inline.Data("🏙️ Санкт-Петербург", "match_city_spb"),
        ),
        inline.Row(
            inline.Data("🏙️ Красноярск", "match_city_krasnoyarsk"),
            inline.Data("🏙️ Великий Новгород", "match_city_novgorod"),
        ),
        inline.Row(
            inline.Data("🏙️ Хабаровск", "match_city_khabarovsk"),
            inline.Data("🏙️ Ярославль", "match_city_yaroslavl"),
        ),
        inline.Row(
            inline.Data("🏙️ Севастополь", "match_city_sevastopol"),
            inline.Data("🏙️ Архангельск", "match_city_arkhangelsk"),
        ),
        inline.Row(
            inline.Data("🏠 В меню", "main_menu"),
        ),
    }
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}

// HandleMatchCity - обработка выбора города
func HandleMatchCity(c telebot.Context, data string) error {
    userID := c.Sender().ID
    city := strings.TrimPrefix(data, "match_city_")
    
    // Сохраняем выбранный город
    matchGames[userID] = &MatchGameState{
        SelectedCity: city,
    }
    
    // Показываем кнопки с купюрами
    text := "Теперь выбери купюру, на которой изображён этот город:"
    
    inline := &telebot.ReplyMarkup{}
    rows := []telebot.Row{
        inline.Row(
            inline.Data("5 ₽", "match_banknote_5"),
            inline.Data("10 ₽", "match_banknote_10"),
            inline.Data("50 ₽", "match_banknote_50"),
        ),
        inline.Row(
            inline.Data("100 ₽", "match_banknote_100"),
            inline.Data("200 ₽", "match_banknote_200"),
            inline.Data("500 ₽", "match_banknote_500"),
        ),
        inline.Row(
            inline.Data("1000 ₽", "match_banknote_1000"),
            inline.Data("2000 ₽", "match_banknote_2000"),
            inline.Data("5000 ₽", "match_banknote_5000"),
        ),
        inline.Row(
            inline.Data("🔙 Назад к городам", "match_back"),
        ),
    }
    
    inline.Inline(rows...)
    
    return c.Edit(text, inline, telebot.ModeMarkdown)
}

// HandleMatchBanknote - обработка выбора купюры
func HandleMatchBanknote(c telebot.Context, data string) error {
    userID := c.Sender().ID
    game, exists := matchGames[userID]
    
    if !exists {
        return c.Send("Сначала выбери город")
    }
    
    banknote := strings.TrimPrefix(data, "match_banknote_")
    
    // Правильные соответствия город -> номинал
    correctPairs := map[string]string{
        "moscow":      "100",
        "spb":         "50",
        "krasnoyarsk": "10",
        "novgorod":    "5",
        "khabarovsk":  "5000",
        "yaroslavl":   "1000",
        "sevastopol":  "200",
        "arkhangelsk": "500",
    }
    
    // Названия городов для отображения
    cityNames := map[string]string{
        "moscow":      "Москва",
        "spb":         "Санкт-Петербург",
        "krasnoyarsk": "Красноярск",
        "novgorod":    "Великий Новгород",
        "khabarovsk":  "Хабаровск",
        "yaroslavl":   "Ярославль",
        "sevastopol":  "Севастополь",
        "arkhangelsk": "Архангельск",
    }
    
    var responseText string
    if correctPairs[game.SelectedCity] == banknote {
        game.Score++
        responseText = fmt.Sprintf("✅ *Правильно!* Город %s изображён на купюре %s ₽\n\n", 
            cityNames[game.SelectedCity], banknote)
        responseText += fmt.Sprintf("Твой счёт: *%d*", game.Score)
    } else {
        responseText = fmt.Sprintf("❌ *Неверно!*\nГород %s изображён на купюре *%s ₽*, а не %s ₽\n\n", 
            cityNames[game.SelectedCity], correctPairs[game.SelectedCity], banknote)
        responseText += fmt.Sprintf("Текущий счёт: *%d*", game.Score)
    }
    
    // Очищаем состояние игры
    delete(matchGames, userID)
    
    inline := &telebot.ReplyMarkup{}
    btnAgain := inline.Data("🔄 Играть ещё", "guess_again")
    btnMenu := inline.Data("🏠 Меню", "main_menu")
    inline.Inline(inline.Row(btnAgain), inline.Row(btnMenu))
    
    return c.Edit(responseText, inline, telebot.ModeMarkdown)
}

// HandleMatchBack - возврат к выбору города
func HandleMatchBack(c telebot.Context) error {
    userID := c.Sender().ID
    delete(matchGames, userID)
    
    return MatchGame(c)
}

// contains - вспомогательная функция для проверки наличия элемента в срезе
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}