package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PhantomX7/go-pos/service/auth"
	"github.com/PhantomX7/go-pos/service/customer"
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/role"
	"github.com/PhantomX7/go-pos/service/stock_mutation"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/service/user"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/utils/database"
	"github.com/PhantomX7/go-pos/utils/validators"

	authHTTP "github.com/PhantomX7/go-pos/service/auth/delivery/http"
	authUsecase "github.com/PhantomX7/go-pos/service/auth/usecase"

	customerHTTP "github.com/PhantomX7/go-pos/service/customer/delivery/http"
	customerRepo "github.com/PhantomX7/go-pos/service/customer/repository/mysql"
	customerUsecase "github.com/PhantomX7/go-pos/service/customer/usecase"

	invoiceHTTP "github.com/PhantomX7/go-pos/service/invoice/delivery/http"
	invoiceRepo "github.com/PhantomX7/go-pos/service/invoice/repository/mysql"
	invoiceUsecase "github.com/PhantomX7/go-pos/service/invoice/usecase"

	productHTTP "github.com/PhantomX7/go-pos/service/product/delivery/http"
	productRepo "github.com/PhantomX7/go-pos/service/product/repository/mysql"
	productUsecase "github.com/PhantomX7/go-pos/service/product/usecase"

	roleRepo "github.com/PhantomX7/go-pos/service/role/repository/mysql"

	stockMutationRepo "github.com/PhantomX7/go-pos/service/stock_mutation/repository/mysql"

	transactionHTTP "github.com/PhantomX7/go-pos/service/transaction/delivery/http"
	transactionRepo "github.com/PhantomX7/go-pos/service/transaction/repository/mysql"
	transactionUsecase "github.com/PhantomX7/go-pos/service/transaction/usecase"

	userHTTP "github.com/PhantomX7/go-pos/service/user/delivery/http"
	userRepo "github.com/PhantomX7/go-pos/service/user/repository/mysql"
	userUsecase "github.com/PhantomX7/go-pos/service/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
)

type repositories struct {
	userRepository          user.UserRepository
	roleRepository          role.RoleRepository
	productRepository       product.ProductRepository
	customerRepository      customer.CustomerRepository
	invoiceRepository       invoice.InvoiceRepository
	stockMutationRepository stock_mutation.StockMutationRepository
	transactionRepository   transaction.TransactionRepository
}

type usecases struct {
	userUsecase        user.UserUsecase
	authUsecase        auth.AuthUsecase
	productUsecase     product.ProductUsecase
	customerUsecase    customer.CustomerUsecase
	invoiceUsecase     invoice.InvoiceUsecase
	transactionUsecase transaction.TransactionUsecase
}

func main() {
	loadEnv()
	db := setupDatabase()

	// for transactions
	database.InitDB(db)

	// init custom validator
	validators.NewValidator(db)

	repositories := initRepository(db)
	usecases := initUsecase(repositories)

	userHandler := resolveUserHandler(usecases)
	authHandler := resolveAuthHandler(usecases)
	productHandler := resolveProductHandler(usecases)
	customerHandler := resolveCustomerHandler(usecases)
	invoiceHandler := resolveInvoiceHandler(usecases)
	transactionHandler := resolveTransactionHandler(usecases)

	startServer(
		userHandler,
		authHandler,
		productHandler,
		customerHandler,
		invoiceHandler,
		transactionHandler,
	)
}

func startServer(handlers ...server.Handler) {

	mConfig := middleware.Config{
	}
	m := middleware.New(mConfig)

	h := server.BuildHandler(m, handlers...)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      h,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		log.Printf("api is available at %s\n", s.Addr)
		if serr := s.ListenAndServe(); serr != http.ErrServerClosed {
			log.Fatal(serr)
		}
	}(s)

	<-sigChan

	log.Println("\nSignal received. Waiting for readiness check...")

	//  wait 15s (kube readiness period second) before shutting down http
	sleep()

	log.Println("\nShutting down the api...")

	err := s.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Something wrong when stopping server : ", err)
		return
	}

	log.Println("api gracefully stopped")
}

func loadEnv() {
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func setupDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func sleep() {
	waitingTime := time.Duration(15)
	if os.Getenv("APP_ENV") != "production" {
		waitingTime = 1
	}
	time.Sleep(waitingTime * time.Second)
}

func resolveUserHandler(usecases usecases) server.Handler {
	return userHTTP.NewUserHandler(usecases.userUsecase)
}

func resolveAuthHandler(usecases usecases) server.Handler {
	return authHTTP.NewAuthHandler(usecases.authUsecase)
}

func resolveProductHandler(usecases usecases) server.Handler {
	return productHTTP.NewProductHandler(usecases.productUsecase)
}

func resolveCustomerHandler(usecases usecases) server.Handler {
	return customerHTTP.NewCustomerHandler(usecases.customerUsecase)
}

func resolveInvoiceHandler(usecases usecases) server.Handler {
	return invoiceHTTP.NewInvoiceHandler(usecases.invoiceUsecase)
}

func resolveTransactionHandler(usecases usecases) server.Handler {
	return transactionHTTP.NewTransactionHandler(usecases.transactionUsecase, usecases.invoiceUsecase)
}

func initRepository(db *gorm.DB) repositories {
	return repositories{
		userRepository:          userRepo.NewUserRepository(db),
		roleRepository:          roleRepo.NewRoleRepository(db),
		productRepository:       productRepo.NewProductRepository(db),
		customerRepository:      customerRepo.NewCustomerRepository(db),
		invoiceRepository:       invoiceRepo.NewInvoiceRepository(db),
		stockMutationRepository: stockMutationRepo.NewStockMutationRepository(db),
		transactionRepository:   transactionRepo.NewTransactionRepository(db),
	}
}

func initUsecase(repositories repositories) usecases {
	return usecases{
		userUsecase: userUsecase.NewUserUsecase(repositories.userRepository),
		authUsecase: authUsecase.NewAuthUsecase(repositories.userRepository, repositories.roleRepository),
		invoiceUsecase: invoiceUsecase.NewInvoiceUsecase(
			repositories.invoiceRepository,
			repositories.customerRepository,
			repositories.transactionRepository,
			repositories.productRepository,
		),
		productUsecase:  productUsecase.NewProductUsecase(repositories.productRepository),
		customerUsecase: customerUsecase.NewCustomerUsecase(repositories.customerRepository),
		transactionUsecase: transactionUsecase.NewTransactionUsecase(
			repositories.transactionRepository,
			repositories.stockMutationRepository,
			repositories.invoiceRepository,
			repositories.productRepository,
		),
	}
}
