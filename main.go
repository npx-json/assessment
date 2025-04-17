package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ipcheck/api"
	"ipcheck/internal/geoip"
	"ipcheck/internal/grpc"

	"database/sql"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	grpcLib "google.golang.org/grpc"
)

func main() {
	var sqlDb *sql.DB
	dbPath := os.Getenv("GEOLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "./GeoLite2-Country.mmdb"
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	server, err := geoip.NewServer(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	// go server.StartAutoUpdate()

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/check", api.MakeCheckHandler(server))

	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: httpMux,
	}

	httpServer.SetKeepAlivesEnabled(false)

	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpcLib.NewServer()
	grpc.Register(grpcServer, server)

	go func() {
		log.Println("HTTP server running on port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("gRPC server running on port", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Connect to MySQL database
	sqlDb, err = connectToMySQL()
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer sqlDb.Close()

	api.SetDB(sqlDb)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("HTTP server shut down")
	}

	log.Println("Shutdown complete")
}

func connectToMySQL() (*sql.DB, error) {

	// Replace with your MySQL connection details
	dsn := os.Getenv("MYSQL_DSN") // Example: "user:password@tcp(localhost:3306)/dbname"
	if dsn == "" {
		dsn = "avoxi:root@tcp(mysql:3306)/avoxi"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to MySQL database")
	return db, nil
}
