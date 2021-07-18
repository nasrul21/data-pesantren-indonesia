# API Data Pesantren Indonesia
API data pesantren from Kemenag Indonesia

## API Endpoint
```
Base URL: https://api-pesantren-indonesia.vercel.app
```
|Endpoint|Method|URL Example|
|---|---|---|
|[/provinsi.json](#get-all-provinsi)|GET|https://api-pesantren-indonesia.vercel.app/provinsi.json|
|[/kabupaten/{id_provinsi}.json](#get-kabupaten-kota-by-id-provinsi)|GET|https://api-pesantren-indonesia.vercel.app/kabupaten/32.json|
|[/pesantren/{id_kab_kota}.json](#get-pesantren-by-id-kabupaten-kota)|GET|https://api-pesantren-indonesia.vercel.app/pesantren/3206.json|
### Get All Provinsi
```
GET /provinsi.json
```
```js
// example : /provinsi.json
[
  {
    "id": "11",
    "nama": "Aceh"
  },
  {
    "id": "12",
    "nama": "Sumatera Utara"
  },
  {
    "id": "13",
    "nama": "Sumatera Barat"
  },
  {
    "id": "14",
    "nama": "Riau"
  },
  ... // and more
]
```

### Get Kabupaten Kota by Id Provinsi
```
GET /kabupaten/{id_provinsi}.json
```
```js
// example: `/kabupaten/32.json` (Jawa Barat)
[
  {
    "id": "3210",
    "nama": "Majalengka"
  },
  {
    "id": "3273",
    "nama": "Kota Bandung"
  },
  {
    "id": "3278",
    "nama": "Kota Tasikmalaya"
  },
  {
    "id": "3203",
    "nama": "Cianjur"
  },
  {
    "id": "3204",
    "nama": "Bandung"
  },
  ... // and more
]
```

### Get Pesantren By Id Kabupaten Kota
```
GET /pesantren/{id_kab_kota}.json
```
```js
// example: `/pesantren/3206.json` (Kab. Tasikmalaya)
[
  {
    "id": "157",
    "nama": "Al-Amanah",
    "nspp": "510032061290",
    "alamat": "Pajaten Ds. Pasirhuni.",
    "kyai": "Aj. Sahid Idrus",
    "kab_kota": {
      "id": "3206",
      "nama": "Tasikmalaya"
    },
    "provinsi": {
      "id": "32",
      "nama": "Jawa Barat"
    }
  },
  {
    "id": "158",
    "nama": "Nurul Huda",
    "nspp": "510032060957",
    "alamat": "Ciaren Ds. Karangjaya.",
    "kyai": "Aj Komarudin",
    "kab_kota": {
      "id": "3206",
      "nama": "Tasikmalaya"
    },
    "provinsi": {
      "id": "32",
      "nama": "Jawa Barat"
    }
  },
  ... // and more
]
```

## Scraper
this scraper will get pesantren data from `ditpdpontren.kemenag.go.id`
```
go run scraper/scraper.go
```
