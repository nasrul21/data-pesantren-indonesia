package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type Provinsi struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type Kabupaten struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type Pesantren struct {
	ID       string    `json:"id"`
	Nama     string    `json:"nama"`
	NSPP     string    `json:"nspp"`
	Alamat   string    `json:"alamat"`
	Kyai     string    `json:"kyai"`
	KabKota  Kabupaten `json:"kab_kota"`
	Provinsi Provinsi  `json:"provinsi"`
}

var client *resty.Client = resty.New()
var listProvinsi []Provinsi

func main() {
	res, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36").SetDoNotParseResponse(true).Get("https://ditpdpontren.kemenag.go.id/pdpp/search")
	if err != nil {
		log.Fatal(err)
	}
	// defer res.RawBody().Close()
	if res.StatusCode() != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode(), res.Status())
	}

	defer res.RawBody().Close()

	fmt.Println(string(res.Body()))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.RawBody())
	if err != nil {
		log.Fatal("aduh : ", err)
	}

	doc.Find(".form-search .provinsi option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		name := s.Text()

		if value != "" {
			provinsi := Provinsi{
				ID:   value,
				Nama: name,
			}
			listProvinsi = append(listProvinsi, provinsi)
			kabupatenKota := getKabupatenKota(provinsi)
			saveToJsonFile(kabupatenKota, fmt.Sprintf("data/kabupaten/%s.json", value))
			fmt.Printf("%v\n", kabupatenKota)
			for _, kabKota := range kabupatenKota {
				fmt.Printf("%s\n", kabKota)
				getPesantren(provinsi, kabKota)
				time.Sleep(1 * time.Second)
			}
		}
	})

	saveToJsonFile(listProvinsi, "data/provinsi.json")
}

func getKabupatenKota(provinsi Provinsi) []Kabupaten {
	var kabKota map[string]string
	res, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36").SetResult(&kabKota).Get("https://ditpdpontren.kemenag.go.id/pdpp/getkabupaten/" + provinsi.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer res.RawBody().Close()
	if res.StatusCode() != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode(), res.Status())
	}

	var listKabupaten []Kabupaten

	for key, value := range kabKota {
		if key != "0" {
			listKabupaten = append(listKabupaten, Kabupaten{
				ID:   key,
				Nama: value,
			})
		}
	}

	return listKabupaten
}

func getPesantren(provinsi Provinsi, kabKota Kabupaten) {
	url := fmt.Sprintf("https://ditpdpontren.kemenag.go.id/pdpp/loadpp?provinsi_id_provinsi=%s&kabupaten_id_kabupaten=%s", provinsi.ID, kabKota.ID)
	fmt.Println(url)
	res, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36").SetDoNotParseResponse(true).Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.RawBody().Close()
	if res.StatusCode() != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode(), res.Status())
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.RawBody())
	if err != nil {
		log.Fatal(err)
	}

	var listPesantren []Pesantren

	doc.Find(".search-result").Each(func(i int, s *goquery.Selection) {
		id := strings.Split(s.Find(".nama-pondok-search a").AttrOr("href", ""), "/")[1]
		nama := s.Find(".nama-pondok-search a").Text()
		desc := strings.TrimSpace(s.Find(".des-search").Text())
		nspp := regexp.MustCompile("NSPP ([0-9]*?) berdiri").FindStringSubmatch(desc)[1]
		alamat := strings.TrimSpace(regexp.MustCompile(`beralamat di ([\s\S]*)`).FindStringSubmatch(desc)[1])
		kyai := strings.TrimSpace(s.Find(".footer-result div").First().Text())

		fmt.Printf("%s\n%s\n%s\n%s\n%s\n%v\n%v\n", id, nama, nspp, alamat, kyai, provinsi, kabKota)
		fmt.Println("==============================================")

		listPesantren = append(listPesantren, Pesantren{
			ID:       id,
			Nama:     nama,
			NSPP:     nspp,
			Alamat:   alamat,
			Kyai:     kyai,
			KabKota:  kabKota,
			Provinsi: provinsi,
		})

	})

	paginationElement := doc.Find(".pagination li")
	lastIndex := paginationElement.Length() - 2
	lastPage, _ := strconv.Atoi(paginationElement.Eq(lastIndex).Text())
	fmt.Printf("LAST PAGE: %d\n", lastPage)
	if lastPage != 0 {
		time.Sleep(1 * time.Second)
		for i := 2; i <= lastPage; i++ {
			url := fmt.Sprintf("https://ditpdpontren.kemenag.go.id/pdpp/loadpp?provinsi_id_provinsi=%s&kabupaten_id_kabupaten=%s&page=%d", provinsi.ID, kabKota.ID, i)
			fmt.Println(url)
			res, err := client.R().SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36").SetDoNotParseResponse(true).Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.RawBody().Close()
			if res.StatusCode() != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode(), res.Status())
			}

			// Load the HTML document
			docPagination, err := goquery.NewDocumentFromReader(res.RawBody())
			if err != nil {
				log.Fatal(err)
			}

			docPagination.Find(".search-result").Each(func(i int, s *goquery.Selection) {
				id := strings.Split(s.Find(".nama-pondok-search a").AttrOr("href", ""), "/")[1]
				nama := s.Find(".nama-pondok-search a").Text()
				desc := strings.TrimSpace(s.Find(".des-search").Text())
				nspp := regexp.MustCompile("NSPP ([0-9]*?) berdiri").FindStringSubmatch(desc)[1]
				alamat := strings.TrimSpace(regexp.MustCompile(`beralamat di ([\s\S]*)`).FindStringSubmatch(desc)[1])
				kyai := strings.TrimSpace(s.Find(".footer-result div").First().Text())

				fmt.Printf("%s\n%s\n%s\n%s\n%s\n%v\n%v\n", id, nama, nspp, alamat, kyai, provinsi, kabKota)
				fmt.Println("==============================================")

				listPesantren = append(listPesantren, Pesantren{
					ID:       id,
					Nama:     nama,
					NSPP:     nspp,
					Alamat:   alamat,
					Kyai:     kyai,
					KabKota:  kabKota,
					Provinsi: provinsi,
				})
			})

			time.Sleep(1 * time.Second)
		}
	}

	saveToJsonFile(listPesantren, fmt.Sprintf("data/pesantren/%s.json", kabKota.ID))
}

func saveToJsonFile(data interface{}, filename string) {
	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile(filename, file, 0644)
}
