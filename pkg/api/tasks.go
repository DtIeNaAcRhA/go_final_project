package api

import (
	"net/http"
	"time"

	"github.com/DtIeNaAcRhA/go_final_project/pkg/db"
)

// limit устанавливает лимит списка задач
var limit = 50

// tasksResp структура для формирования JSON с параметром tasks - спиок задач.
type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler - хендлер для получения списка задач с параметрами поиска, если в параметрах URL запроса установлно занчение search, и без них.
// Реализует метод Get.
// В теле ответа возвращает JSON с списком задач - tasks или ошибку - error.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, errorResp{Error: "Method not allowed"}, false)
		return
	}

	search := r.URL.Query().Get("search")
	tasks := make([]*db.Task, 0, limit)

	var err error
	if search != "" {

		parsedDate, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = db.TasksWithSearch(limit, search, false)

		} else {
			tasks, err = db.TasksWithSearch(limit, parsedDate.Format(DataFormat), true)

		}
	} else {
		tasks, err = db.TasksWithOutSearch(limit)

	}
	if err != nil {
		writeJson(w, errorResp{
			Error: err.Error(),
		}, false)
		return
	}
	writeJson(w, tasksResp{
		Tasks: tasks,
	}, true)
}
