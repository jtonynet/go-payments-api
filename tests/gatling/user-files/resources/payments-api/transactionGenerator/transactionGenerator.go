package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

var generate = "account" // account | account_categories | transactions | payload

type MccMerchant struct {
	MCC      string
	Merchant string
}

func main() {
	start := 2501
	end := 3000

	if generate == "account" {
		for i := start; i <= end; i++ {
			fmt.Println(generateAccount())
		}
	}

	if generate == "account_categories" {
		for i := start; i <= end; i++ {
			accountCategories := fmt.Sprintf(
				"(%d, 1, NOW(), NOW()), (%d, 2, NOW(), NOW()), (%d, 3, NOW(), NOW()), (%d, 4, NOW(), NOW()), (%d, 5, NOW(), NOW()),",
				i, i, i, i, i,
			)
			fmt.Println(accountCategories)
		}
	}

	if generate == "transactions" {
		for i := start; i <= end; i++ {
			transactions := fmt.Sprintf(
				"(%d, 91350.00, 1, NOW(), NOW()),(%d, 91350.00, 2, NOW(), NOW()),(%d, 91300.00, 3, NOW(), NOW()),(%d, 91350.00, 4, NOW(), NOW()),(%d, 91350.00, 5, NOW(), NOW()),",
				i, i, i, i, i,
			)
			fmt.Println(transactions)
		}
	}

	if generate == "payload" {
		generatePayload()
	}
}

func generateAccount() string {
	return fmt.Sprintf("('%s', '%s', NOW(), NOW()),", uuid.NewString(), uniqueName())
}

func uniqueName() string {
	nouns1 := []string{
		"Almeidinha", "Saldanha", "Crespo", "Guerrero", "Lindomar", "Moscow", "Kiev", "Seul", "Foreman", "Sena",
		"Farinhas", "Carlson", "Smith", "Kaio", "Fialho", "Billian", "Anderssen", "Donald", "Jeff", "Haffman",
		"Kobold", "Piscian", "Tesla", "Cobalt", "Niquel", "Itsuragi", "Combs", "Hermes", "Arceu", "Pluto",
		"Cyan", "Caligui", "Silvio", "Coppola", "Cage", "Nichols", "Sabina", "Suyara", "Infantino", "Payne",
		"Phellipe", "Donnatelo", "Donavan", "Firmino", "Fabriete", "Falsiane", "Falseite", "Fagoberto", "Valino", "Rulio"}

	nouns2 := []string{
		"Kingston", "Dakar", "Vienna", "Rashid", "Firenze", "Calabria", "Lisboa", "Doha", "Samarkand", "Kilimanjaro",
		"Bogota", "Damascus", "Luanda", "Oslo", "Quito", "Cochabamba", "Zurich", "Kyoto", "Havana", "Mumbai",
		"Tijuana", "Zanzibar", "Sapporo", "Orleans", "Geneva", "Sofia", "Athens", "Sevilla", "Tromsø", "Sydney",
		"Anchorage", "Casablanca", "Istanbul", "Krakow", "Tallinn", "Nairobi", "Belfast", "Reykjavik", "Cairo", "Malaga",
		"Santiago", "Cordoba", "Napoli", "Goiania", "Manaus", "Belgrade", "Chisinau", "Baku", "Bruges", "Galway",
	}

	nouns3 := []string{
		"Machado", "Fontana", "Perez", "Cordova", "Ronaldo", "Carvalho", "Avelar", "Bennett", "Torres", "Calheiros",
		"Yamamoto", "Shivani", "Kumari", "Ankara", "Gwangju", "Marseille", "Ottawa", "Osaka", "Canberra", "Montevideo",
		"Cusco", "Tbilisi", "Addis", "Ababa", "Helsinki", "Chongqing", "Harare", "Kigali", "Vilnius", "Minsk",
		"Gaborone", "Pretoria", "Maputo", "Dushanbe", "Bishkek", "Vientiane", "Phnom", "Penh", "Ulaanbaatar", "Vladivostok",
		"Tangier", "Alexandria", "Batumi", "Granada", "Jerez", "Leonardo", "Borges", "Griffin", "Simpson", "Agatha",
	}

	nouns4 := []string{
		"Salazar", "Furtado", "Eisenhower", "Macchiato", "Zaire", "Fukuoka", "Hanover", "Cambridge", "Boston", "Savannah",
		"Nashville", "Phoenix", "Brisbane", "Adelaide", "Valparaiso", "Juliaca", "Puno", "Cajamarca", "Arequipa", "Pucallpa",
		"Lima", "Callao", "Trujillo", "Rojas", "Ambato", "Latacunga", "Quito", "Cali", "Medellin", "Cartagena",
		"Bogota", "Barranquilla", "Pastor", "Manizales", "Armenia", "Pereira", "Ibague", "Guayaquil", "Paulo", "Salvador",
		"Recife", "Belem", "Fortaleza", "Brasilia", "Curitiba", "Porto", "Manaus", "Cuiaba", "CampoGrande", "Pessoa",
	}

	return fmt.Sprintf("%s %s %s %s",
		nouns1[rand.Intn(len(nouns2))],
		nouns2[rand.Intn(len(nouns4))],
		nouns3[rand.Intn(len(nouns1))],
		nouns4[rand.Intn(len(nouns3))],
	)
}

func generatePayload() {
	//Genereted accountUIDs here "000-000,000-000,..."
	uuidListStr := "d16e42ac-dc9d-46ef-8c5d-7a6f5207eaab,0a7eb116-8114-4e7b-b81c-3ce14d4b4102,..."

	rand.Seed(time.Now().UnixNano())
	uuidList := strings.Split(uuidListStr, ",")

	mccMerchantMap := []MccMerchant{
		{MCC: "5811", Merchant: "UBER EATS                   SAO PAULO BR"},
		{MCC: "5812", Merchant: "PAG*JoseDaSilva          RIO DE JANEI BR"},
		{MCC: "5811", Merchant: "TAXI*SUPERMERCADO        RIO DE JANEI BR"},
		{MCC: "6411", Merchant: "99TAXI                      SAO PAULO BR"},
		{MCC: "6411", Merchant: "LYFT                     RIO DE JANEI BR"},
		{MCC: "6411", Merchant: "RADIO*TAXI                  SAO PAULO BR"},
		{MCC: "7411", Merchant: "DROGAS*L                 RIO DE JANEI BR"},
		{MCC: "7411", Merchant: "POLICLINICA BRA          RIO DE JANEI BR"},
		{MCC: "7411", Merchant: "YVO_PYTANGUI                SAO PAULO BR"},
		{MCC: "5411", Merchant: "POLICLINICA*LANCHONETE   RIO DE JANEI BR"},
		{MCC: "5411", Merchant: "BURGUER*DONALD`S         RIO DE JANEI BR"},
	}

	for _, uuid := range uuidList {
		for _, mccMerchant := range mccMerchantMap {
			totalAmount := rand.Float64()*(4.99-1.01) + 1.01
			fmt.Printf(`{"account": "%s", "mcc": "%s", "merchant": "%s", "totalAmount": %.2f}`,
				uuid,
				mccMerchant.MCC,
				mccMerchant.Merchant,
				totalAmount,
			)
			fmt.Print("\n")
		}
	}
}
