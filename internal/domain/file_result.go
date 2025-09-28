package domain

type FileResult struct {
	URL   string `json:"url"`             //ссылка на скачиваемый файл
	Path  string `json:"path,omitempty"`  //путь по которому хранится скаченный файл
	Error error  `json:"error,omitempty"` //ошибка обработки ссылки
	Ok    bool   `json:"ok"`              //флаг обработки ссылки
}
