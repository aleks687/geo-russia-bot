package main

import (
    "fmt"
    "log"
    "math/rand"
    "os"
    "time"
    "strings"

    "github.com/joho/godotenv"
    "gopkg.in/telebot.v3"
)

// Структуры данных
type Banknote struct {
    Nominal     string
    City        string
    Description string
    Facts       []string
}

type Question struct {
    Text    string
    Options []string
    Correct int
    Fact    string
}

// Хранилище состояний пользователей
var userStates = make(map[int64]*UserState)

type UserState struct {
    CurrentGame string
    Score       int
    Questions   []Question
    CurrentQ    int
    Banknote    *Banknote
    Cities      []string
    CorrectIdx  int
}

// База данных купюр (20 городов!)
var banknotes = []Banknote{
    // Оригинальные купюры
    {
        Nominal:     "5 рублей",
        City:        "Великий Новгород",
        Description: "Памятник 'Тысячелетие России' и Софийский собор",
        Facts: []string{
            "Памятник установлен в 1862 году",
            "Софийский собор - древнейший каменный храм России",
        },
    },
    {
        Nominal:     "10 рублей",
        City:        "Красноярск",
        Description: "Часовня Параскевы Пятницы и Красноярская ГЭС",
        Facts: []string{
            "Часовня изображена на купюре с 1997 года",
            "Красноярская ГЭС - одна из мощнейших в России",
        },
    },
    {
        Nominal:     "50 рублей",
        City:        "Санкт-Петербург",
        Description: "Ростральная колонна",
        Facts: []string{
            "Ростральные колонны использовались как маяки",
            "Высота колонн - 32 метра",
        },
    },
    {
        Nominal:     "100 рублей",
        City:        "Москва",
        Description: "Большой театр",
        Facts: []string{
            "Здание Большого театра открыто в 1825 году",
            "Квадрига Аполлона - символ русского искусства",
        },
    },
    {
        Nominal:     "200 рублей",
        City:        "Севастополь",
        Description: "Памятник затопленным кораблям",
        Facts: []string{
            "Памятник - символ Севастополя",
            "Установлен в 1905 году",
        },
    },
    {
        Nominal:     "500 рублей",
        City:        "Архангельск",
        Description: "Памятник Петру I",
        Facts: []string{
            "Памятник Петру I установлен в 1914 году",
            "Соловецкий монастырь основан в XV веке",
        },
    },
    {
        Nominal:     "1000 рублей",
        City:        "Ярославль",
        Description: "Церковь Иоанна Предтечи",
        Facts: []string{
            "Церковь внесена в список ЮНЕСКО",
            "Ярославль основан в 1010 году",
        },
    },
    {
        Nominal:     "2000 рублей",
        City:        "Владивосток",
        Description: "Мост на остров Русский и космодром Восточный",
        Facts: []string{
            "Введена в обращение в 2017 году",
            "Мост на Русский остров - один из самых длинных вантовых мостов в мире",
            "Космодром Восточный - первый российский гражданский космодром",
        },
    },
    {
        Nominal:     "5000 рублей",
        City:        "Хабаровск",
        Description: "Мост через Амур",
        Facts: []string{
            "Мост - самый длинный на Транссибе (3890 м)",
            "Открыт в 1916 году",
        },
    },
    // Добавляем больше городов для игры "Угадай город"
    {
        Nominal:     "10 рублей (старая)",
        City:        "Новосибирск",
        Description: "Оперный театр и метромост",
        Facts: []string{
            "Новосибирск - третий по численности город России",
            "Оперный театр - один из крупнейших в России",
        },
    },
    {
        Nominal:     "50 рублей (старая)",
        City:        "Екатеринбург",
        Description: "Памятник основателям города и бизнес-центр 'Высоцкий'",
        Facts: []string{
            "Екатеринбург - столица Урала",
            "Город основан в 1723 году",
        },
    },
    {
        Nominal:     "100 рублей (старая)",
        City:        "Ростов-на-Дону",
        Description: "Набережная и памятник Дмитрию Ростовскому",
        Facts: []string{
            "Ростов-на-Дону - порт пяти морей",
            "Основан в 1749 году",
        },
    },
    {
        Nominal:     "500 рублей (старая)",
        City:        "Нижний Новгород",
        Description: "Нижегородский кремль и ярмарка",
        Facts: []string{
            "Нижний Новгород основан в 1221 году",
            "Кремль - один из самых сохранившихся в России",
        },
    },
    {
        Nominal:     "1000 рублей (старая)",
        City:        "Казань",
        Description: "Казанский кремль и башня Сююмбике",
        Facts: []string{
            "Казань - столица Татарстана",
            "Кремль входит в список ЮНЕСКО",
        },
    },
    {
        Nominal:     "5000 рублей (старая)",
        City:        "Самара",
        Description: "Площадь Куйбышева и монумент Славы",
        Facts: []string{
            "Самара - крупный центр авиакосмической промышленности",
            "Площадь Куйбышева - одна из крупнейших в Европе",
        },
    },
    {
        Nominal:     "Памятная банкнота",
        City:        "Сочи",
        Description: "Олимпийский парк и горнолыжные курорты",
        Facts: []string{
            "Сочи принимал зимнюю Олимпиаду в 2014 году",
            "Город-курорт на Черном море",
        },
    },
    {
        Nominal:     "Памятная банкнота",
        City:        "Калининград",
        Description: "Кафедральный собор и Музей Мирового океана",
        Facts: []string{
            "Калининград - самый западный регион России",
            "Ранее назывался Кёнигсберг",
        },
    },
    {
        Nominal:     "Памятная банкнота",
        City:        "Иркутск",
        Description: "Набережная Ангары и деревянное зодчество",
        Facts: []string{
            "Иркутск - ворота Байкала",
            "Основан в 1661 году",
        },
    },
    {
        Nominal:     "Памятная банкнота",
        City:        "Петрозаводск",
        Description: "Набережная Онежского озера",
        Facts: []string{
            "Петрозаводск - столица Карелии",
            "Основан Петром I в 1703 году",
        },
    },
    {
        Nominal:     "Памятная банкнота",
        City:        "Мурманск",
        Description: "Памятник Защитникам Заполярья",
        Facts: []string{
            "Мурманск - крупнейший город за полярным кругом",
            "Алёша - памятник высотой 42 метра",
        },
    },
}

// Расширенная база вопросов для викторины
var quizQuestions = []Question{
    {
        Text:    "Какой город изображён на 100-рублёвой купюре?",
        Options: []string{"Москва", "Санкт-Петербург", "Казань", "Новгород"},
        Correct: 0,
        Fact:    "На 100 рублях изображён Большой театр в Москве",
    },
    {
        Text:    "Что изображено на 10-рублёвой купюре?",
        Options: []string{"Мост", "Часовня", "Собор", "Кремль"},
        Correct: 1,
        Fact:    "На 10 рублях - часовня Параскевы Пятницы в Красноярске",
    },
    {
        Text:    "Где находится Ростральная колонна?",
        Options: []string{"Москва", "Санкт-Петербург", "Кронштадт", "Выборг"},
        Correct: 1,
        Fact:    "Ростральная колонна в Санкт-Петербурге на 50 рублях",
    },
    {
        Text:    "Какой город на 5000-рублёвой купюре?",
        Options: []string{"Москва", "Хабаровск", "Владивосток", "Новосибирск"},
        Correct: 1,
        Fact:    "На 5000 рублях изображён мост через Амур в Хабаровске",
    },
    {
        Text:    "Что изображено на 200-рублёвой купюре?",
        Options: []string{"Космодром", "Мост", "Памятник кораблям", "Вокзал"},
        Correct: 2,
        Fact:    "На 200 рублях - памятник затопленным кораблям в Севастополе",
    },
    {
        Text:    "Какой памятник на 5-рублёвой купюре?",
        Options: []string{"Минину и Пожарскому", "Тысячелетие России", "Медный всадник", "Ленину"},
        Correct: 1,
        Fact:    "На 5 рублях - памятник 'Тысячелетие России' в Великом Новгороде",
    },
    {
        Text:    "Что изображено на 2000-рублёвой купюре?",
        Options: []string{"Мост в Крым", "Мост на Русский остров", "Космодром Восточный", "Оба варианта"},
        Correct: 3,
        Fact:    "На 2000 рублях изображены мост на Русский остров и космодром Восточный во Владивостоке",
    },
    {
        Text:    "На какой купюре изображён Красноярск?",
        Options: []string{"5 рублей", "10 рублей", "50 рублей", "100 рублей"},
        Correct: 1,
        Fact:    "Красноярск с часовней Параскевы Пятницы на 10 рублях",
    },
    {
        Text:    "Что изображено на 1000-рублёвой купюре?",
        Options: []string{"Церковь Иоанна Предтечи", "Успенский собор", "Храм Василия Блаженного", "Спас на Крови"},
        Correct: 0,
        Fact:    "На 1000 рублях - церковь Иоанна Предтечи в Ярославле",
    },
    {
        Text:    "Какой город на 500-рублёвой купюре?",
        Options: []string{"Мурманск", "Архангельск", "Петрозаводск", "Вологда"},
        Correct: 1,
        Fact:    "На 500 рублях - Архангельск с памятником Петру I",
    },
    {
        Text:    "Какой город является столицей Урала?",
        Options: []string{"Челябинск", "Екатеринбург", "Пермь", "Тюмень"},
        Correct: 1,
        Fact:    "Екатеринбург - неофициальная столица Урала",
    },
    {
        Text:    "Какой город стоит на реке Волга?",
        Options: []string{"Нижний Новгород", "Новосибирск", "Екатеринбург", "Красноярск"},
        Correct: 0,
        Fact:    "Нижний Новгород стоит на месте слияния Волги и Оки",
    },
}

func main() {
    rand.Seed(time.Now().UnixNano())

    godotenv.Load()
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN не установлен")
    }

    pref := telebot.Settings{
        Token:  token,
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    }

    bot, err := telebot.NewBot(pref)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Бот запущен...")

    // === СОЗДАНИЕ МЕНЮ ===
    
    // Главное меню
    mainMenu := &telebot.ReplyMarkup{ResizeKeyboard: true}
    btnBanknotes := mainMenu.Text("🏦 Все купюры")
    btnQuiz := mainMenu.Text("❓ Викторина")
    btnGames := mainMenu.Text("🎮 Игры")
    btnAbout := mainMenu.Text("ℹ️ О проекте")
    mainMenu.Reply(
        mainMenu.Row(btnBanknotes, btnQuiz),
        mainMenu.Row(btnGames, btnAbout),
    )

    // Меню игр
    gamesMenu := &telebot.ReplyMarkup{ResizeKeyboard: true}
    btnGuess := gamesMenu.Text("🎯 Угадай город")
    btnBackToMain := gamesMenu.Text("🏠 Главное меню")
    gamesMenu.Reply(
        gamesMenu.Row(btnGuess),
        gamesMenu.Row(btnBackToMain),
    )

    // === ОБРАБОТЧИКИ КОМАНД ===
    
    // /start
    bot.Handle("/start", func(c telebot.Context) error {
        return c.Send("🇷🇺 *Добро пожаловать!*\n\nЯ бот о географии России на денежных знаках. Выбери раздел:", mainMenu, telebot.ModeMarkdown)
    })

    // === ОБРАБОТЧИКИ ГЛАВНОГО МЕНЮ ===
    
    // Все купюры
    bot.Handle(&btnBanknotes, func(c telebot.Context) error {
        text := "🏦 *Купюры России*\n\nВыбери номинал:"
        
        inline := &telebot.ReplyMarkup{}
        btn5 := inline.Data("5 ₽", "banknote_5")
        btn10 := inline.Data("10 ₽", "banknote_10")
        btn50 := inline.Data("50 ₽", "banknote_50")
        btn100 := inline.Data("100 ₽", "banknote_100")
        btn200 := inline.Data("200 ₽", "banknote_200")
        btn500 := inline.Data("500 ₽", "banknote_500")
        btn1000 := inline.Data("1000 ₽", "banknote_1000")
        btn2000 := inline.Data("2000 ₽", "banknote_2000")
        btn5000 := inline.Data("5000 ₽", "banknote_5000")
        btnMenu := inline.Data("🏠 Главное меню", "main_menu")
        
        inline.Inline(
            inline.Row(btn5, btn10, btn50),
            inline.Row(btn100, btn200, btn500),
            inline.Row(btn1000, btn2000, btn5000),
            inline.Row(btnMenu),
        )
        
        return c.Send(text, inline, telebot.ModeMarkdown)
    })

    // О проекте
    bot.Handle(&btnAbout, func(c telebot.Context) error {
        text := "ℹ️ *О проекте*\n\n" +
            "Этот бот создан для изучения географии России через денежные знаки.\n\n" +
            "*Команды:*\n" +
            "/start - Главное меню\n" +
            "🏦 Все купюры - список купюр\n" +
            "❓ Викторина - проверь знания\n" +
            "🎮 Игры - увлекательные игры\n\n" +
            "В игре 'Угадай город' теперь 20 городов!\n\n" +
            "Версия: 2.0"
        return c.Send(text, mainMenu, telebot.ModeMarkdown)
    })

    // Игры
    bot.Handle(&btnGames, func(c telebot.Context) error {
        return c.Send("🎮 *Выбери игру:*", gamesMenu, telebot.ModeMarkdown)
    })

    // Назад в главное меню (из меню игр)
    bot.Handle(&btnBackToMain, func(c telebot.Context) error {
        return c.Send("Главное меню:", mainMenu, telebot.ModeMarkdown)
    })

    // Викторина
    bot.Handle(&btnQuiz, func(c telebot.Context) error {
        userID := c.Sender().ID
        
        // Берем 5 случайных вопросов из 12
        questions := make([]Question, 5)
        indices := rand.Perm(len(quizQuestions))[:5]
        for i, idx := range indices {
            questions[i] = quizQuestions[idx]
        }
        
        userStates[userID] = &UserState{
            CurrentGame: "quiz",
            Questions:   questions,
            CurrentQ:    0,
            Score:       0,
        }
        
        return sendQuizQuestion(c, userID, questions[0], 0, 5)
    })

    // Угадай город (с 20 городами)
    bot.Handle(&btnGuess, func(c telebot.Context) error {
        userID := c.Sender().ID
        
        // Выбираем случайную купюру из 20
        banknote := banknotes[rand.Intn(len(banknotes))]
        
        // Создаем варианты ответов (4 города)
        cities := []string{banknote.City}
        
        // Список всех городов для вариантов
        allCities := []string{
            "Москва", "Санкт-Петербург", "Казань", "Нижний Новгород",
            "Новосибирск", "Екатеринбург", "Самара", "Омск",
            "Челябинск", "Ростов-на-Дону", "Уфа", "Красноярск",
            "Пермь", "Воронеж", "Волгоград", "Краснодар",
            "Саратов", "Тюмень", "Тольятти", "Ижевск",
            "Барнаул", "Ульяновск", "Иркутск", "Хабаровск",
            "Ярославль", "Владивосток", "Томск", "Оренбург",
            "Кемерово", "Новокузнецк", "Рязань", "Астрахань",
            "Набережные Челны", "Пенза", "Липецк", "Киров",
            "Чебоксары", "Калининград", "Брянск", "Курск",
            "Иваново", "Магнитогорск", "Тверь", "Ставрополь",
            "Симферополь", "Севастополь", "Архангельск", "Владимир",
            "Смоленск", "Мурманск", "Петрозаводск", "Великий Новгород",
        }
        
        for len(cities) < 4 {
            randomCity := allCities[rand.Intn(len(allCities))]
            if !contains(cities, randomCity) && randomCity != banknote.City {
                cities = append(cities, randomCity)
            }
        }
        
        // Перемешиваем варианты
        rand.Shuffle(len(cities), func(i, j int) {
            cities[i], cities[j] = cities[j], cities[i]
        })
        
        // Находим индекс правильного ответа
        correctIdx := 0
        for i, city := range cities {
            if city == banknote.City {
                correctIdx = i
                break
            }
        }
        
        // Сохраняем состояние
        userStates[userID] = &UserState{
            CurrentGame: "guess_city",
            Banknote:    &banknote,
            Cities:      cities,
            CorrectIdx:  correctIdx,
            Score:       0,
        }
        
        text := fmt.Sprintf("🎯 *Угадай город* (из 20 городов)\n\n") +
            fmt.Sprintf("На купюре *%s* изображён:\n", banknote.Nominal) +
            fmt.Sprintf("_%s_\n\n", banknote.Description) +
            "Какой это город?"
        
        inline := &telebot.ReplyMarkup{}
        var rows []telebot.Row
        
        for i, city := range cities {
            data := fmt.Sprintf("guess_%d", i)
            btn := inline.Data(city, data)
            rows = append(rows, inline.Row(btn))
        }
        
        exitBtn := inline.Data("🚪 Выйти", "game_exit")
        rows = append(rows, inline.Row(exitBtn))
        
        inline.Inline(rows...)
        
        return c.Send(text, inline, telebot.ModeMarkdown)
    })

    // === УНИВЕРСАЛЬНЫЙ ОБРАБОТЧИК CALLBACK ===
    bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
        rawData := c.Data()
        data := strings.TrimSpace(rawData)
        
        log.Printf("Callback: %s", data)
        
        c.Respond()
        
        userID := c.Sender().ID
        state := userStates[userID]
        
        switch {
        case data == "main_menu":
            return c.Edit("Главное меню:", mainMenu, telebot.ModeMarkdown)
            
        case data == "game_exit":
            delete(userStates, userID)
            c.Delete()
            return c.Send("Игра завершена.", gamesMenu, telebot.ModeMarkdown)
            
        case data == "quiz_exit":
            delete(userStates, userID)
            c.Delete()
            return c.Send("Викторина завершена.", mainMenu, telebot.ModeMarkdown)
            
        case data == "quiz_again":
            delete(userStates, userID)
            return startQuiz(c)
            
        case strings.HasPrefix(data, "guess_"):
            if state == nil || state.CurrentGame != "guess_city" {
                return c.Send("Игра не найдена. Начните новую.")
            }
            
            var answerIdx int
            fmt.Sscanf(data, "guess_%d", &answerIdx)
            
            var responseText string
            if answerIdx == state.CorrectIdx {
                state.Score++
                responseText = "✅ *Правильно!* Молодец!\n\n"
            } else {
                responseText = fmt.Sprintf("❌ *Неверно!*\nПравильный ответ: *%s*\n\n", 
                    state.Banknote.City)
            }
            
            responseText += fmt.Sprintf("Твой счёт: *%d*\n\nХочешь сыграть ещё?", state.Score)
            
            delete(userStates, userID)
            
            inline := &telebot.ReplyMarkup{}
            btnAgain := inline.Data("🔄 Ещё раз", "guess_again")
            btnMenu := inline.Data("🏠 Меню", "main_menu")
            inline.Inline(inline.Row(btnAgain), inline.Row(btnMenu))
            
            return c.Edit(responseText, inline, telebot.ModeMarkdown)
            
        case data == "guess_again":
            delete(userStates, userID)
            return startGuessGame(c)
            
        case strings.HasPrefix(data, "quiz_"):
            if state == nil || state.CurrentGame != "quiz" {
                return c.Send("Викторина не найдена. Начните новую.")
            }
            
            parts := strings.Split(data, "_")
            if len(parts) != 2 {
                return nil
            }
            
            var answerIdx int
            fmt.Sscanf(parts[1], "%d", &answerIdx)
            
            question := state.Questions[state.CurrentQ]
            
            var responseText string
            if answerIdx == question.Correct {
                state.Score++
                responseText = "✅ *Правильно!*\n\n"
            } else {
                responseText = fmt.Sprintf("❌ *Неверно!*\nПравильный ответ: *%s*\n\n", 
                    question.Options[question.Correct])
            }
            
            responseText += "📚 *Факт:* " + question.Fact
            
            c.Send(responseText, telebot.ModeMarkdown)
            
            state.CurrentQ++
            
            if state.CurrentQ >= len(state.Questions) {
                return finishQuiz(c, userID, state.Score, len(state.Questions))
            }
            
            return sendQuizQuestion(c, userID, 
                state.Questions[state.CurrentQ], 
                state.CurrentQ, 
                len(state.Questions))
            
        case strings.HasPrefix(data, "banknote_"):
            nominal := strings.TrimPrefix(data, "banknote_")
            
            var banknote *Banknote
            for i, bn := range banknotes {
                if strings.TrimSuffix(bn.Nominal, " рублей") == nominal {
                    banknote = &banknotes[i]
                    break
                }
            }
            
            if banknote == nil {
                return nil
            }
            
            text := fmt.Sprintf("💵 *%s*\n\n", banknote.Nominal) +
                fmt.Sprintf("🏙️ *Город:* %s\n", banknote.City) +
                fmt.Sprintf("🏛️ *Что изображено:* %s\n\n", banknote.Description) +
                "📚 *Интересные факты:*\n"
            
            for i, fact := range banknote.Facts {
                text += fmt.Sprintf("%d. %s\n", i+1, fact)
            }
            
            inline := &telebot.ReplyMarkup{}
            btnBack := inline.Data("⬅️ Назад к списку", "back_to_banknotes")
            btnMenu := inline.Data("🏠 Меню", "main_menu")
            inline.Inline(inline.Row(btnBack), inline.Row(btnMenu))
            
            return c.Edit(text, inline, telebot.ModeMarkdown)
            
        case data == "back_to_banknotes":
            text := "🏦 *Купюры России*\n\nВыбери номинал:"
            
            inline := &telebot.ReplyMarkup{}
            btn5 := inline.Data("5 ₽", "banknote_5")
            btn10 := inline.Data("10 ₽", "banknote_10")
            btn50 := inline.Data("50 ₽", "banknote_50")
            btn100 := inline.Data("100 ₽", "banknote_100")
            btn200 := inline.Data("200 ₽", "banknote_200")
            btn500 := inline.Data("500 ₽", "banknote_500")
            btn1000 := inline.Data("1000 ₽", "banknote_1000")
            btn2000 := inline.Data("2000 ₽", "banknote_2000")
            btn5000 := inline.Data("5000 ₽", "banknote_5000")
            btnMenu := inline.Data("🏠 Главное меню", "main_menu")
            
            inline.Inline(
                inline.Row(btn5, btn10, btn50),
                inline.Row(btn100, btn200, btn500),
                inline.Row(btn1000, btn2000, btn5000),
                inline.Row(btnMenu),
            )
            
            return c.Edit(text, inline, telebot.ModeMarkdown)
        }
        
        return nil
    })

    log.Println("✅ Бот готов к работе!")
    bot.Start()
}

// Вспомогательные функции
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func sendQuizQuestion(c telebot.Context, userID int64, question Question, index, total int) error {
    text := fmt.Sprintf("❓ *Вопрос %d/%d*\n\n%s", index+1, total, question.Text)
    
    inline := &telebot.ReplyMarkup{}
    var rows []telebot.Row
    
    for i, option := range question.Options {
        data := fmt.Sprintf("quiz_%d", i)
        btn := inline.Data(option, data)
        rows = append(rows, inline.Row(btn))
    }
    
    exitBtn := inline.Data("🚪 Выйти", "quiz_exit")
    rows = append(rows, inline.Row(exitBtn))
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}

func finishQuiz(c telebot.Context, userID int64, score, total int) error {
    delete(userStates, userID)
    
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

func startQuiz(c telebot.Context) error {
    userID := c.Sender().ID
    
    questions := make([]Question, 5)
    indices := rand.Perm(len(quizQuestions))[:5]
    for i, idx := range indices {
        questions[i] = quizQuestions[idx]
    }
    
    userStates[userID] = &UserState{
        CurrentGame: "quiz",
        Questions:   questions,
        CurrentQ:    0,
        Score:       0,
    }
    
    return sendQuizQuestion(c, userID, questions[0], 0, 5)
}

func startGuessGame(c telebot.Context) error {
    userID := c.Sender().ID
    
    banknote := banknotes[rand.Intn(len(banknotes))]
    
    cities := []string{banknote.City}
    
    allCities := []string{
        "Москва", "Санкт-Петербург", "Казань", "Нижний Новгород",
        "Новосибирск", "Екатеринбург", "Самара", "Омск",
        "Челябинск", "Ростов-на-Дону", "Уфа", "Красноярск",
        "Пермь", "Воронеж", "Волгоград", "Краснодар",
        "Саратов", "Тюмень", "Тольятти", "Ижевск",
        "Барнаул", "Ульяновск", "Иркутск", "Хабаровск",
        "Ярославль", "Владивосток", "Томск", "Оренбург",
        "Кемерово", "Новокузнецк", "Рязань", "Астрахань",
        "Набережные Челны", "Пенза", "Липецк", "Киров",
        "Чебоксары", "Калининград", "Брянск", "Курск",
        "Иваново", "Магнитогорск", "Тверь", "Ставрополь",
        "Симферополь", "Севастополь", "Архангельск", "Владимир",
        "Смоленск", "Мурманск", "Петрозаводск", "Великий Новгород",
    }
    
    for len(cities) < 4 {
        randomCity := allCities[rand.Intn(len(allCities))]
        if !contains(cities, randomCity) && randomCity != banknote.City {
            cities = append(cities, randomCity)
        }
    }
    
    rand.Shuffle(len(cities), func(i, j int) {
        cities[i], cities[j] = cities[j], cities[i]
    })
    
    correctIdx := 0
    for i, city := range cities {
        if city == banknote.City {
            correctIdx = i
            break
        }
    }
    
    userStates[userID] = &UserState{
        CurrentGame: "guess_city",
        Banknote:    &banknote,
        Cities:      cities,
        CorrectIdx:  correctIdx,
        Score:       0,
    }
    
    text := fmt.Sprintf("🎯 *Угадай город* (из 20 городов)\n\n") +
        fmt.Sprintf("На купюре *%s* изображён:\n", banknote.Nominal) +
        fmt.Sprintf("_%s_\n\n", banknote.Description) +
        "Какой это город?"
    
    inline := &telebot.ReplyMarkup{}
    var rows []telebot.Row
    
    for i, city := range cities {
        data := fmt.Sprintf("guess_%d", i)
        btn := inline.Data(city, data)
        rows = append(rows, inline.Row(btn))
    }
    
    exitBtn := inline.Data("🚪 Выйти", "game_exit")
    rows = append(rows, inline.Row(exitBtn))
    
    inline.Inline(rows...)
    
    return c.Send(text, inline, telebot.ModeMarkdown)
}