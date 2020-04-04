package edbo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/yevchuk-kostiantyn/stats-export/model"
)

const (
	host   = "vstup2018.edbo.gov.ua"
	origin = "https://vstup2018.edbo.gov.ua"
)

func (c *Client) GetUniversities(edBase, spec, region string) ([]model.University, error) {
	urlReq := fmt.Sprint(c.baseURL)
	form := url.Values{}
	form.Add("action", "universities")
	form.Add("qualification", "1")
	form.Add("education-base", edBase)
	form.Add("speciality", spec)
	form.Add("region", region)
	form.Add("education-form", "1")
	form.Add("course", "1")
	req, err := http.NewRequest(http.MethodPost, urlReq, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", host)
	req.Header.Add("Origin", origin)
	referer := fmt.Sprintf("https://vstup2018.edbo.gov.ua/offers/?qualification=1&education-base=%s&speciality=%s&education-form=1&course=1&region=%s",
		edBase, spec, region)
	req.Header.Add("Referer", referer)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("request err:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		log.Println("decode err:", err)
		return nil, err
	}
	data, ok := result["universities"]
	if !ok {
		log.Println("cannot interpret resp")
		// TODO fix error msg
		return nil, fmt.Errorf("")
	}

	universities := arrToStruct(data.([]interface{}))
	return universities, nil
}

func (c *Client) GetVolumesPerOffer(edBase, spec, region string, offer []string) (map[float64][]float64, error) {
	urlReq := fmt.Sprint(c.baseURL)
	form := url.Values{}
	form.Add("action", "offers")
	form.Add("usids", strings.Join(offer, ","))
	req, err := http.NewRequest(http.MethodPost, urlReq, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", host)
	req.Header.Add("Origin", origin)
	referer := fmt.Sprintf("https://vstup2018.edbo.gov.ua/offers/?qualification=1&education-base=%s&speciality=%s&education-form=1&course=1&region=%s",
		edBase, spec, region)
	req.Header.Add("Referer", referer)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("request err:", err)
		return nil, err
	}

	var response model.EDBOOfferResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	result := make(map[float64][]float64)
	for _, offer := range response.Offers {
		budgetVolume := offer.([]interface{})[14].(float64)
		licenseVolume := offer.([]interface{})[16].(float64)
		contractVolume := offer.([]interface{})[17].(float64)
		result[offer.([]interface{})[0].(float64)] = []float64{budgetVolume, licenseVolume, contractVolume}
	}
	return result, nil
}

func arrToStruct(in []interface{}) []model.University {
	var out []model.University
	for _, u := range in {
		out = append(out, model.University{
			Code:   u.([]interface{})[0].(float64),
			Name:   u.([]interface{})[1].(string),
			Number: u.([]interface{})[2].(float64),
			USIDs:  strings.Split(u.([]interface{})[3].(string), ","),
		})
	}
	return out
}
