package handlers

import (
	"fmt"
	"github.com/IoanaAdrian/HubA/src/api/models"
	"github.com/IoanaAdrian/HubA/src/api/services"
	"github.com/IoanaAdrian/HubA/src/api/services/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HomePost(w http.ResponseWriter, r *http.Request) {
	amountString := r.URL.Query().Get("amount")
	creationDateString := r.URL.Query().Get("creation_date")
	isIncome := r.URL.Query().Get("is_income")
	description := r.URL.Query().Get("description")

	amountValue, err := strconv.Atoi(amountString)
	if services.LogError(services.LogErrorArgs{Err: err, W: w}) {
		return
	}

	var creationDate time.Time

	if len(creationDateString) > 0 {
		creationDate, err = time.Parse("02-01-2006", creationDateString)
		if services.LogError(services.LogErrorArgs{Err: err, W: w}) {
			return
		}
	} else {
		creationDate = time.Now()
	}

	transaction := models.Transaction{Amount: uint(amountValue), CreationDate: creationDate, IsIncome: isIncome == "true", Description: description}
	result := database.DB.Create(&transaction)

	if services.LogError(services.LogErrorArgs{Err: result.Error, W: w, Fatal: true}) {
		return
	}

	log.Print("Successfully inserted data!")

}
func HomeGet(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction
	result := database.DB.Order("creation_date desc").Find(&transactions)

	if services.LogError(services.LogErrorArgs{Err: result.Error, W: w, Fatal: true}) {
		return
	}

	// TODO: Replace with table that stores amount of income and spendings

	var income, spending uint

	err := database.DB.Table("transactions").Select("sum(CASE WHEN is_income = 1 THEN amount ELSE 0 END)").Row().Scan(&income)
	if services.LogError(services.LogErrorArgs{Err: err, W: w, Fatal: true}) {
		return
	}

	err = database.DB.Table("transactions").Select("sum(CASE WHEN is_income = 0 THEN amount ELSE 0 END)").Row().Scan(&spending)
	if services.LogError(services.LogErrorArgs{Err: err, W: w, Fatal: true}) {
		return
	}


	fmt.Fprintf(w, fmt.Sprintf("INCOME %d", income))
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, fmt.Sprintf("SPENDING %d", spending))
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "---------------------------")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "INCOME TRANSACTIONS")
	fmt.Fprintf(w, "\n")

	for _, transaction := range transactions {
		if transaction.IsIncome {
			fmt.Fprintf(w, fmt.Sprintf("%s   %s  %s", transaction.CreationDate, strconv.Itoa(int(transaction.Amount)), transaction.Description))
			fmt.Fprintf(w, "\n")
		}
	}

	fmt.Fprintf(w, "---------------------------")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "SPENDING TRANSACTIONS")
	fmt.Fprintf(w, "\n")

	for _, transaction := range transactions {
		if !transaction.IsIncome {
			fmt.Fprintf(w, fmt.Sprintf("%s   %s  %s", transaction.CreationDate, strconv.Itoa(int(transaction.Amount)), transaction.Description))
			fmt.Fprintf(w, "\n")
		}
	}
}
