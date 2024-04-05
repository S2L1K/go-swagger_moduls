package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"studentgit.kata.academy/SLK/go-swagger/models"
)

// HandlerSearch отправляет запрос на dadata.ru и возвращает получнный ответ клиенту
// @Summary Получение данных адреса
// @Description Поиск информации об адресе по его фактическому адресу
// @Param input body models.SearchRequest  true "фактический адрес"
// @Success      200  {object}  models.SearchResponse "информация об адресе"
// @Failure      400  {string}  string "400 :Неверный формат запроса"
// @Failure      500  {string}  string "500: Сервис dadata.ru не доступен"
// @Router       /api/address/search [post]
func HandlerSearch(w http.ResponseWriter, r *http.Request) {
	bodyJS, errIRA := io.ReadAll(r.Body)
	if errIRA != nil {
		fmt.Println(errIRA, 6)
		return
	}

	var searchRequest models.SearchRequest
	errJUM := json.Unmarshal(bodyJS, &searchRequest)
	if errJUM != nil {
		fmt.Println(errJUM, 5)
		return
	}

	servReq := []string{searchRequest.Query}
	servReqJS, errJM := json.Marshal(servReq)
	if errJM != nil {
		fmt.Println(errJM, 4)
		return
	}

	client := &http.Client{}

	req, errNR := http.NewRequest("POST", "https://cleaner.dadata.ru/api/v1/clean/address", bytes.NewBuffer(servReqJS))
	if errNR != nil {
		fmt.Println(errNR, 123123)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token 479572091c2fed0fba3c3f6f2467ad005541911f")
	req.Header.Set("X-Secret", "c7e794ec88c793e0bf1f3322d43ee7a5bc1ea64e")

	res, errCD := client.Do(req)

	if errCD != nil {
		fmt.Println(errCD)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Сервис dadata.ru не доступен"))
		return
	}
	defer res.Body.Close()

	body, errIRA := io.ReadAll(res.Body)
	if errIRA != nil {
		fmt.Println(errIRA, 3)
	}

	var serachResp models.SearchResponse

	errJUM = json.Unmarshal(body, &serachResp.Addresses)
	if errJUM != nil {
		fmt.Println(errJUM, 2)
		return
	}

	jsonout, errJM := json.Marshal(serachResp.Addresses)
	if errJM != nil {
		fmt.Println(errJM, 1)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(jsonout)
	if err != nil {
		fmt.Println(err, 124151234)
		return
	}
}
