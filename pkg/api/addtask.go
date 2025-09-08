package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DtIeNaAcRhA/go_final_project/pkg/db"
)

// addTaskResponse - структура для формирования JSON, котороый содержит id, успешно созданной задачи.
type addTaskResponse struct {
	ID string `json:"id"`
}

// checkDate - функция для проверки даты задачи.
// Проверяет формат даты, актуальность даты и переопределяет дату в случае её неактуальности.
// Принимает структуру db.Task, возвращает ошибку.
func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DataFormat)
		return nil
	}
	t, err := time.Parse(DataFormat, task.Date)
	if err != nil {
		return fmt.Errorf("Дата представлена в формате, отличном от 20060102")
	}
	// если сегодня (now) больше task.Date (t)
	if now.After(t) && now.Format(DataFormat) != task.Date {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(DataFormat)
			return nil
		}
		// в противном случае, берём следующую дату
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Правило повторения указано в неправильном формате")
		}
		task.Date = next
	}
	return nil
}

// addTaskHandler - хендлер для добавления задачи.
// В теле запроса ожидается JSON с параметрами задачи,
// В качестве ответа возвращает JSON с параметром id, успешно добавленной задачи или JSON c ошибкой - парметр error.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var addResponse addTaskResponse
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, errorResp{Error: "Ошибка десериализации JSON"}, false)
		return
	}

	if task.Title == "" {
		writeJson(w, errorResp{Error: "Не указан заголовок задачи"}, false)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addResponse.ID = strconv.Itoa(int(id))
	writeJson(w, addResponse, true)
}
