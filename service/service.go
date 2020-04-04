package service

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/yevchuk-kostiantyn/stats-export/model"

	"github.com/gorilla/mux"
	"github.com/yevchuk-kostiantyn/stats-export/rest/edbo"
)

var regions = map[string]string{
	"Вінницька":         "05",
	"КИЇВ":              "80",
	"Волинська":         "07",
	"Дніпропетровська":  "12",
	"Донецька":          "14",
	"Житомирська":       "18",
	"Закарпатська":      "21",
	"Запорізька":        "23",
	"Івано-Франківська": "26",
	"Київська":          "32",
	"Кіровоградська":    "35",
	"Луганська":         "44",
	"Львівська":         "46",
	"Миколаївська":      "48",
	"Одеська":           "51",
	"Полтавська":        "53",
	"Рівненська":        "56",
	"Сумська":           "59",
	"Тернопільска":      "61",
	"Харківська":        "63",
	"Херсонська":        "65",
	"Хмельницька":       "68",
	"Черкаська":         "71",
	"Чернівецька":       "73",
	"Чернігівська":      "74",
}

type Service struct {
	EDBO *edbo.Client
}

func NewService(edbo *edbo.Client) *Service {
	return &Service{EDBO: edbo}
}

func (s *Service) Run(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/export", s.ExportStats).Methods(http.MethodGet)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Printf("[error] %s, exiting now ...", err.Error())
		return
	}
}

func (s *Service) ExportStats(w http.ResponseWriter, r *http.Request) {
	// 40 and 122
	//eduBase := r.URL.Query().Get("eduBase")
	eduBase := "40"
	//if eduBase == "" {
	//	http.Error(w, fmt.Sprint("education-base is not set"), http.StatusBadRequest)
	//	return
	//}

	//speciality := r.URL.Query().Get("speciality")
	speciality := "122"
	//if speciality == "" {
	//	http.Error(w, fmt.Sprint("speciality is not set"), http.StatusBadRequest)
	//	return
	//}
	log.Printf("export stats - edu-base:%s, speciality:%s", eduBase, speciality)
	result, err := s.fetchStats(eduBase, speciality)
	if err != nil {
		// TODO refactor status code definition
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for a, b := range result {
		log.Println("a:", a)
		log.Println("b:", b)
	}
	encodeResult(w, result)
}

func (s *Service) fetchStats(eduBase, speciality string) ([]model.ExportResponse, error) {
	var result []model.ExportResponse
	for name, code := range regions {
		log.Printf("fetching data for %s region", name)
		unis, err := s.EDBO.GetUniversities(eduBase, speciality, code)
		if err != nil {
			return nil, err
		}
		log.Printf("universities list fetched: %d", len(unis))
		var totalBudgetVolume float64
		var totalLicenseVolume float64
		var totalContractVolume float64
		for _, uni := range unis {
			res, err := s.EDBO.GetVolumesPerOffer(eduBase, speciality, code, uni.USIDs)
			if err != nil {
				return nil, err
			}
			for _, b := range res {
				totalBudgetVolume += b[0]
				totalLicenseVolume += b[1]
				totalContractVolume += b[2]
			}
			// increase limit for 80 and 63 regions
			if code == "80" || code == "63" {
				time.Sleep(20 * time.Second)
				continue
			}
			time.Sleep(10 * time.Second)
		}
		log.Printf("name:%s budget: %f license: %f contract: %f", name, totalBudgetVolume, totalLicenseVolume, totalContractVolume)
		res := model.ExportResponse{
			Region:         name,
			Speciality:     speciality,
			LicenseVolume:  totalLicenseVolume,
			BudgetVolume:   totalBudgetVolume,
			ContractVolume: totalContractVolume,
		}
		result = append(result, res)
		// increase limit for 80 and 63 regions
		if code == "80" || code == "63" {
			time.Sleep(40 * time.Second)
			continue
		}
		time.Sleep(25 * time.Second)
	}
	return result, nil
}

func encodeResult(w http.ResponseWriter, res interface{}) {
	if res == nil || reflect.ValueOf(res).IsNil() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println("error occurred while encoding object")
	}
}
