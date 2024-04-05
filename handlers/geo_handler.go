package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"studentgit.kata.academy/SLK/go-swagger/models"
)

// HandlerGeo отправляет запрос на dadata.ru и возвращает полученый ответ клиенту.
// @Summary Получение данных адреса.
// @Description Поиск информации об адресе по его геолокационным данным.
// @Param input body models.GeocodeRequest  true "координаты".
// @Success      200  {object}  models.GeocodeResponse "информация об адресе"
// @Failure      400  {string}  string "400 :Неверный формат запроса"
// @Failure      500  {string}  string "500: Сервис dadata.ru не доступен"
// @Router       /api/address/geocode [post]
func HandlerGeo(w http.ResponseWriter, r *http.Request) {
	bodyJS, errIRA := io.ReadAll(r.Body)
	if errIRA != nil {
		fmt.Println(errIRA)
		return
	}

	client := &http.Client{}

	req, errNR := http.NewRequest("POST", "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", bytes.NewBuffer(bodyJS))
	if errNR != nil {
		fmt.Println(errNR)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token 479572091c2fed0fba3c3f6f2467ad005541911f")
	req.Header.Set("X-Secret", "c7e794ec88c793e0bf1f3322d43ee7a5bc1ea64e")

	resp, errCD := client.Do(req)
	if errCD != nil {
		fmt.Println(errCD)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Сервис dadata.ru не доступен"))
		return
	}
	defer resp.Body.Close()

	var geoCode models.GeoCode

	errJND := json.NewDecoder(resp.Body).Decode(&geoCode)
	if errJND != nil {
		fmt.Println(errJND)
		return
	}

	var res []models.Address
	for _, r := range geoCode.Suggestions {
		var address models.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		res = append(res, address)
	}

	resJSON, errJM := json.Marshal(res)
	if errJM != nil {
		fmt.Println(errJM)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(resJSON)
	if err != nil {
		fmt.Println(err)
		return
	}
}
