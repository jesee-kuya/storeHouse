package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

// Router sets up the HTTP router with all handlers
func Router(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Initialize handlers
	accountHandler := NewAccountHandler(db)
	transactionHandler := NewTransactionHandler(db)
	memberHandler := NewMemberHandler(db)
	userHandler := NewUserHandler(db)
	expenditureHandler := NewExpenditureHandler(db)
	transferHandler := NewTransferHandler(db)
	receiptHandler := NewReceiptHandler(db)
	membersGroupHandler := NewMembersGroupHandler(db)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Accounts
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", accountHandler.GetAllAccounts)
			r.Post("/", accountHandler.CreateAccount)
			r.Get("/{id}", accountHandler.GetAccount)
			r.Put("/{id}", accountHandler.UpdateAccount)
			r.Delete("/{id}", accountHandler.DeactivateAccount)
		})

		// Transactions
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", transactionHandler.GetAllTransactions)
			r.Post("/", transactionHandler.CreateTransaction)
			r.Get("/{id}", transactionHandler.GetTransaction)
			r.Get("/ref/{ref}", transactionHandler.GetTransactionByRef)
			r.Get("/account/{accountID}", transactionHandler.GetTransactionsByAccount)
			r.Get("/member/{memberID}", transactionHandler.GetTransactionsByMember)
			r.Get("/type/{type}", transactionHandler.GetTransactionsByType)
			r.Get("/date-range", transactionHandler.GetTransactionsByDateRange)
			r.Put("/{id}", transactionHandler.UpdateTransaction)
			r.Delete("/{id}", transactionHandler.DeleteTransaction)
		})

		// Members
		r.Route("/members", func(r chi.Router) {
			r.Get("/", memberHandler.GetAllMembers)
			r.Get("/search", memberHandler.SearchMembers)
			r.Post("/", memberHandler.CreateMember)
			r.Get("/{id}", memberHandler.GetMember)
			r.Get("/phone/{phone}", memberHandler.GetMemberByPhone)
			r.Get("/email/{email}", memberHandler.GetMemberByEmail)
			r.Get("/group/{groupID}", memberHandler.GetMembersByGroup)
			r.Put("/{id}", memberHandler.UpdateMember)
			r.Delete("/{id}", memberHandler.DeleteMember)
		})

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAllUsers)
			r.Get("/active", userHandler.GetActiveUsers)
			r.Get("/role/{role}", userHandler.GetUsersByRole)
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Get("/username/{username}", userHandler.GetUserByUsername)
			r.Get("/email/{email}", userHandler.GetUserByEmail)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
			r.Post("/{id}/deactivate", userHandler.DeactivateUser)
			r.Post("/authenticate", userHandler.AuthenticateUser)
			r.Post("/{id}/change-password", userHandler.ChangePassword)
		})

		// Expenditures
		r.Route("/expenditures", func(r chi.Router) {
			r.Get("/", expenditureHandler.GetAllExpenditures)
			r.Post("/", expenditureHandler.CreateExpenditure)
			r.Get("/{id}", expenditureHandler.GetExpenditure)
			r.Get("/transaction/{transactionID}", expenditureHandler.GetExpendituresByTransaction)
			r.Put("/{id}", expenditureHandler.UpdateExpenditure)
			r.Delete("/{id}", expenditureHandler.DeleteExpenditure)
		})

		// Transfers
		r.Route("/transfers", func(r chi.Router) {
			r.Get("/", transferHandler.GetAllTransfers)
			r.Post("/", transferHandler.CreateTransfer)
			r.Get("/{id}", transferHandler.GetTransfer)
			r.Get("/transaction/{transactionID}", transferHandler.GetTransfersByTransaction)
			r.Get("/credit-account/{accountID}", transferHandler.GetTransfersByCreditAccount)
			r.Get("/credit-account/{accountID}/total", transferHandler.GetTotalTransfersByCreditAccount)
			r.Get("/date-range", transferHandler.GetTransfersByDateRange)
			r.Get("/date-range/{accountID}", transferHandler.GetTotalTransfersByDateRange)
			r.Put("/{id}", transferHandler.UpdateTransfer)
			r.Delete("/{id}", transferHandler.DeleteTransfer)
		})

		// Receipts
		r.Route("/receipts", func(r chi.Router) {
			r.Get("/", receiptHandler.GetAllReceipts)
			r.Post("/", receiptHandler.CreateReceipt)
			r.Get("/{id}", receiptHandler.GetReceipt)
			r.Get("/transaction/{transactionID}", receiptHandler.GetReceiptsByTransaction)
			r.Get("/account/{accountID}", receiptHandler.GetReceiptsByAccount)
			r.Get("/account/{accountID}/total", receiptHandler.GetTotalReceiptsByAccount)
			r.Get("/date-range", receiptHandler.GetReceiptsByDateRange)
			r.Get("/date-range/{accountID}", receiptHandler.GetTotalReceiptsByDateRange)
			r.Put("/{id}", receiptHandler.UpdateReceipt)
			r.Delete("/{id}", receiptHandler.DeleteReceipt)
		})

		// Members Groups
		r.Route("/groups", func(r chi.Router) {
			r.Get("/", membersGroupHandler.GetAllGroups)
			r.Get("/with-count", membersGroupHandler.GetGroupsWithMemberCount)
			r.Post("/", membersGroupHandler.CreateGroup)
			r.Get("/{id}", membersGroupHandler.GetGroup)
			r.Get("/{id}/count", membersGroupHandler.GetGroupMemberCount)
			r.Get("/name/{name}", membersGroupHandler.GetGroupByName)
			r.Put("/{id}", membersGroupHandler.UpdateGroup)
			r.Delete("/{id}", membersGroupHandler.DeleteGroup)
		})
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return router
}

// StartServer starts the HTTP server
func StartServer(db *sqlx.DB) {
	router := Router(db)

	port := ":8080"
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📋 API Documentation available at http://localhost%s/api/v1", port)
	log.Printf("💚 Health check available at http://localhost%s/health", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
