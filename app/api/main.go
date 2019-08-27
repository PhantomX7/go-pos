package main

import (
	"context"
	"fmt"
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/service/order_transaction"
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
	"github.com/PhantomX7/go-pos/service/return_transaction"
	"github.com/PhantomX7/go-pos/service/role"
	"github.com/PhantomX7/go-pos/service/stock_adjustment"
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

	orderInvoiceHTTP "github.com/PhantomX7/go-pos/service/order_invoice/delivery/http"
	orderInvoiceRepo "github.com/PhantomX7/go-pos/service/order_invoice/repository/mysql"
	orderInvoiceUsecase "github.com/PhantomX7/go-pos/service/order_invoice/usecase"

	returnTransactionHTTP "github.com/PhantomX7/go-pos/service/return_transaction/delivery/http"
	returnTransactionRepo "github.com/PhantomX7/go-pos/service/return_transaction/repository/mysql"
	returnTransactionUsecase "github.com/PhantomX7/go-pos/service/return_transaction/usecase"

	productHTTP "github.com/PhantomX7/go-pos/service/product/delivery/http"
	productRepo "github.com/PhantomX7/go-pos/service/product/repository/mysql"
	productUsecase "github.com/PhantomX7/go-pos/service/product/usecase"

	roleRepo "github.com/PhantomX7/go-pos/service/role/repository/mysql"

	stockAdjustmentRepo "github.com/PhantomX7/go-pos/service/stock_adjustment/repository/mysql"

	stockMutationRepo "github.com/PhantomX7/go-pos/service/stock_mutation/repository/mysql"

	orderTransactionHTTP "github.com/PhantomX7/go-pos/service/order_transaction/delivery/http"
	orderTransactionRepo "github.com/PhantomX7/go-pos/service/order_transaction/repository/mysql"
	orderTransactionUsecase "github.com/PhantomX7/go-pos/service/order_transaction/usecase"

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
	userRepository              user.UserRepository
	roleRepository              role.RoleRepository
	productRepository           product.ProductRepository
	customerRepository          customer.CustomerRepository
	invoiceRepository           invoice.InvoiceRepository
	stockAdjustmentRepository   stock_adjustment.StockAdjustmentRepository
	stockMutationRepository     stock_mutation.StockMutationRepository
	transactionRepository       transaction.TransactionRepository
	returnTransactionRepository return_transaction.ReturnTransactionRepository
	orderInvoiceRepository      order_invoice.OrderInvoiceRepository
	orderTransactionRepository  order_transaction.OrderTransactionRepository
}

type usecases struct {
	userUsecase              user.UserUsecase
	authUsecase              auth.AuthUsecase
	productUsecase           product.ProductUsecase
	customerUsecase          customer.CustomerUsecase
	invoiceUsecase           invoice.InvoiceUsecase
	transactionUsecase       transaction.TransactionUsecase
	returnTransactionUsecase return_transaction.ReturnTransactionUsecase
	orderInvoiceUsecase      order_invoice.OrderInvoiceUsecase
	orderTransactionUsecase  order_transaction.OrderTransactionUsecase
}

func main() {
	loadEnv()
	db := setupDatabase()

	// for transactions
	database.InitDB(db)

	// init custom validator
	validators.NewValidator(db)

	repositories := initRepositories(db)
	usecases := initUsecases(repositories)

	userHandler := resolveUserHandler(usecases)
	authHandler := resolveAuthHandler(usecases)
	productHandler := resolveProductHandler(usecases)
	customerHandler := resolveCustomerHandler(usecases)
	invoiceHandler := resolveInvoiceHandler(usecases)
	transactionHandler := resolveTransactionHandler(usecases)
	returnTransactionHandler := resolveReturnTransactionHandler(usecases)
	orderInvoiceHandler := resolveOrderInvoiceHandler(usecases)
	orderTransactionHandler := resolveOrderTransaction(usecases)

	startServer(
		userHandler,
		authHandler,
		productHandler,
		customerHandler,
		invoiceHandler,
		transactionHandler,
		returnTransactionHandler,
		orderInvoiceHandler,
		orderTransactionHandler,
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

func initRepositories(db *gorm.DB) repositories {
	return repositories{
		userRepository:              userRepo.New(db),
		roleRepository:              roleRepo.New(db),
		productRepository:           productRepo.New(db),
		customerRepository:          customerRepo.New(db),
		invoiceRepository:           invoiceRepo.New(db),
		stockAdjustmentRepository:   stockAdjustmentRepo.New(db),
		stockMutationRepository:     stockMutationRepo.New(db),
		transactionRepository:       transactionRepo.New(db),
		returnTransactionRepository: returnTransactionRepo.New(db),
		orderInvoiceRepository:      orderInvoiceRepo.New(db),
		orderTransactionRepository:  orderTransactionRepo.New(db),
	}
}

func initUsecases(r repositories) usecases {
	return usecases{
		userUsecase: userUsecase.New(r.userRepository),
		authUsecase: authUsecase.New(r.userRepository, r.roleRepository),
		invoiceUsecase: invoiceUsecase.New(
			r.invoiceRepository,
			r.customerRepository,
			r.transactionRepository,
			r.productRepository,
			r.returnTransactionRepository,
		),
		productUsecase:  productUsecase.New(r.productRepository, r.stockAdjustmentRepository),
		customerUsecase: customerUsecase.New(r.customerRepository),
		transactionUsecase: transactionUsecase.New(
			r.transactionRepository,
			r.stockMutationRepository,
			r.invoiceRepository,
			r.productRepository,
			r.returnTransactionRepository,
		),
		returnTransactionUsecase: returnTransactionUsecase.New(
			r.returnTransactionRepository,
			r.transactionRepository,
			r.stockMutationRepository,
			r.productRepository,
		),
		orderInvoiceUsecase: orderInvoiceUsecase.New(
			r.orderInvoiceRepository,
			r.orderTransactionRepository,
			r.productRepository,
		),
		orderTransactionUsecase: orderTransactionUsecase.New(
			r.orderTransactionRepository,
			r.stockMutationRepository,
			r.orderInvoiceRepository,
			r.productRepository,
		),
	}
}

func resolveUserHandler(usecases usecases) server.Handler {
	return userHTTP.New(usecases.userUsecase)
}

func resolveAuthHandler(usecases usecases) server.Handler {
	return authHTTP.New(usecases.authUsecase)
}

func resolveProductHandler(usecases usecases) server.Handler {
	return productHTTP.New(usecases.productUsecase)
}

func resolveCustomerHandler(usecases usecases) server.Handler {
	return customerHTTP.New(usecases.customerUsecase)
}

func resolveInvoiceHandler(usecases usecases) server.Handler {
	return invoiceHTTP.New(usecases.invoiceUsecase)
}

func resolveTransactionHandler(usecases usecases) server.Handler {
	return transactionHTTP.New(usecases.transactionUsecase, usecases.invoiceUsecase)
}

func resolveReturnTransactionHandler(usecases usecases) server.Handler {
	return returnTransactionHTTP.New(usecases.returnTransactionUsecase, usecases.invoiceUsecase)
}

func resolveOrderInvoiceHandler(usecases usecases) server.Handler {
	return orderInvoiceHTTP.New(usecases.orderInvoiceUsecase)
}

func resolveOrderTransaction(usecases usecases) server.Handler {
	return orderTransactionHTTP.New(usecases.orderTransactionUsecase, usecases.orderInvoiceUsecase)
}
