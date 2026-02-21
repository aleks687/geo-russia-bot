package models

import "math/rand"

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

type UserState struct {
    CurrentGame   string
    Score         int
    QuestionIndex int
}

var GameStorage = make(map[int64]*UserState)

var Banknotes = []Banknote{
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
        Nominal:     "5000 рублей",
        City:        "Хабаровск",
        Description: "Мост через Амур",
        Facts: []string{
            "Мост - самый длинный на Транссибе (3890 м)",
            "Открыт в 1916 году",
        },
    },
}

var QuizQuestions = []Question{
    {
        Text:    "Какой город на 100-рублёвой купюре?",
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
}

func GetRandomBanknote() Banknote {
    return Banknotes[rand.Intn(len(Banknotes))]
}

func GetRandomQuestions(count int) []Question {
    if count > len(QuizQuestions) {
        count = len(QuizQuestions)
    }
    
    shuffled := make([]Question, len(QuizQuestions))
    copy(shuffled, QuizQuestions)
    rand.Shuffle(len(shuffled), func(i, j int) {
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    })
    
    return shuffled[:count]
}