package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate - возвращает следующую дату задачи при указаном парвиле повторения и ошибку,
// принимает текущую дату, строку - дата задачи в установленном формате, правило повторения.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	//проверка repeat на пустоту значения
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is not specified")
	}

	//парсинг времени задачи и проверка его на корректность формата
	startData, err := time.Parse(DataFormat, dstart)
	if err != nil {
		return "", err
	}

	//парсинг правила
	repeatArray := strings.Split(repeat, " ")

	switch repeatArray[0] {
	case "y":
		return caseY(now, startData), nil
	case "d":
		return caseD(now, startData, repeatArray)
	case "w":
		return caseW(now, startData, repeatArray)
	case "m":
		return caseM(now, startData, repeatArray)
	default:
		return "", fmt.Errorf("incorrect format")
	}

}

// caseY - реализует правило повторения задачи - ежигодно.
// Принимает текущую дату и дату задачи (time.Time), возвращает строку - следующую дату.
// Вызывается в функции NextDate.
func caseY(now, startData time.Time) string {
	for {
		startData = startData.AddDate(1, 0, 0)
		if startData.After(now) {
			break
		}
	}
	return startData.Format(DataFormat)
}

// caseD - реализует правило повторения задачи - с интервалом в днях.
// Принимает текущую дату, дату задачи (time.Time), слайс строк - правло повтарения, возвращает строку - следующую дату и ошибку.
// Вызывается в функции NextDate.
func caseD(now, startData time.Time, repeatArray []string) (string, error) {

	if len(repeatArray) < 2 {
		return "", fmt.Errorf("incorrect format")
	}
	days, err := strconv.Atoi(repeatArray[1])
	if err != nil {
		return "", err
	}
	if days > 400 {
		return "", fmt.Errorf("the maximum allowed interval has been exceeded")
	}
	for {
		startData = startData.AddDate(0, 0, days)
		if startData.After(now) {
			break
		}
	}
	return startData.Format(DataFormat), nil
}

// caseW - реализует правило повторения задачи - еженеделно.
// Принимает текущую дату, дату задачи (time.Time), слайс строк - правло повтарения, возвращает строку - следующую дату и ошибку.
// Вызывается в функции NextDate.
func caseW(now, startData time.Time, repeatArray []string) (string, error) {
	//проверка корректонсти формата
	if len(repeatArray) < 2 {
		return "", fmt.Errorf("incorrect format")
	}

	//парсинг дней недели
	weekdaysString := strings.Split(repeatArray[1], ",")
	weekdays := make([]int, 0, len(weekdaysString))
	for _, v := range weekdaysString {
		weekday, err := strconv.Atoi(strings.TrimSpace(v))
		//проверка допустимых значений
		if err != nil || weekday < 1 || weekday > 7 {
			return "", fmt.Errorf("incorrect format")
		}

		weekday = weekday % 7 //воскресенье ==0
		weekdays = append(weekdays, weekday)
	}
	sort.Ints(weekdays)

	//релизация логики переноса
	if startData.Before(now) {
		startData = now
	}

	check := false
	for _, weekday := range weekdays {
		//если в списке есть будущие дни недели на текущей неделе
		if weekday-int(startData.Weekday()) > 0 {
			startData = startData.AddDate(0, 0, weekday-int(startData.Weekday()))
			check = true
			break
		}
	}
	if !check {
		difference := 7 - int(startData.Weekday()) //сколько дней до понедельника след. недели
		startData = startData.AddDate(0, 0, difference+weekdays[0])
	}

	return startData.Format(DataFormat), nil
}

// caseM - реализует правило повторения задачи - еженемесячно.
// Принимает текущую дату, дату задачи (time.Time), слайс строк - правло повтарения, возвращает строку - следующую дату и ошибку.
// Вызывается в функции NextDate.
func caseM(now, startData time.Time, repeatArray []string) (string, error) {
	if len(repeatArray) < 2 {
		return "", fmt.Errorf("incorrect format")
	}

	var days [32]bool
	var months [13]bool
	LastDays := make([]int, 0, 2)

	//парсинг чисел
	daysString := strings.Split(repeatArray[1], ",")
	for _, v := range daysString {
		day, err := strconv.Atoi(strings.TrimSpace(v))
		//проверка допустимых значений
		if err != nil || day < -2 || day > 31 || day == 0 {
			return "", fmt.Errorf("incorrect format")
		}
		if day == -1 {
			LastDays = append(LastDays, -1)
		} else if day == -2 {
			LastDays = append(LastDays, -2)
		} else {
			days[day] = true
		}
	}
	sort.Ints(LastDays) //сначала предпоследний день, потом последний

	switch len(repeatArray) {
	case 2:
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
	case 3:
		//парсинг месяцев
		monthsString := strings.Split(repeatArray[2], ",")
		for _, v := range monthsString {
			month, err := strconv.Atoi(strings.TrimSpace(v))
			//проверка допустимых значений
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("incorrect format")
			}
			months[month] = true
		}
	default:
		return "", fmt.Errorf("incorrect format")
	}

	//реализация логики переноса
	if startData.Before(now) {
		startData = now
	}
	finishYear := startData.Year() + 1
	for startData.Year() <= finishYear {
		startData = startData.AddDate(0, 0, 1)
		m := int(startData.Month())
		if !months[m] {
			firstOfNextMonth := time.Date(startData.Year(), startData.Month(), 1, 0, 0, 0, 0, startData.Location()).AddDate(0, 1, 0)
			newData := firstOfNextMonth.AddDate(0, 0, -1)
			startData = newData
			continue
		}
		if days[startData.Day()] {
			return startData.Format(DataFormat), nil
		}
		if len(LastDays) > 0 {
			firstOfNextMonth := time.Date(startData.Year(), startData.Month(), 1, 0, 0, 0, 0, startData.Location()).AddDate(0, 1, 0)
			for _, i := range LastDays { //-2,-1 (отсортировано от меньшего к большиму)
				wantedDay := firstOfNextMonth.AddDate(0, 0, i)
				if startData.Day() == wantedDay.Day() {
					return startData.Format(DataFormat), nil
				}
			}
		}

	}
	if days[29] && months[2] {
		startData = time.Date(startData.Year(), time.Month(2), 29, 0, 0, 0, 0, startData.Location())
		for startData.Day() != 29 && int(startData.Month()) != 2 {
			startData = time.Date(startData.Year(), time.Month(2), 29, 0, 0, 0, 0, startData.Location()).AddDate(1, 0, 0)
		}
		return startData.Format(DataFormat), nil
	}
	return "", fmt.Errorf("incorrect format")
}

// nextDayHandler - хендлер, который изменят дату задачи, если установлено правило повторения.
// Реализует метод Get.
// Ожидает в запросе форму с параметрами для определения следующей даты задачи.
// В теле ответа возвращает следубщую дату или пустую строку.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из запроса с помощью FormValue
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "Missing required parameter: date", http.StatusBadRequest)
		return
	}

	if repeatParam == "" {
		http.Error(w, "Missing required parameter: repeat", http.StatusBadRequest)
		return
	}

	// Определяем текущее время (now)
	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DataFormat, nowParam)
		if err != nil {
			http.Error(w, "Invalid now parameter format", http.StatusBadRequest)
			return
		}
	}

	// Вызываем функцию NextDate
	result, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
