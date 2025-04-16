package geoip

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
)

type Server struct {
	db      *geoip2.Reader
	dbMutex sync.RWMutex
}

func NewServer(dbPath string) (*Server, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Server{db: db}, nil
}

func (s *Server) CheckIP(ipStr string, allowed []string) (bool, string, error) {
	// ip := net.ParseIP("194.33.45.162")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, "", fmt.Errorf("invalid IP")
	}
	s.dbMutex.RLock()
	record, err := s.db.Country(ip)
	s.dbMutex.RUnlock()
	if err != nil {
		return false, "", err
	}
	code := record.Country.IsoCode
	for _, a := range allowed {
		if strings.EqualFold(a, code) {
			return true, code, nil
		}
	}
	return false, code, nil
}

func (s *Server) StartAutoUpdate() {
	licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
	if licenseKey == "" {
		log.Println("No license key for GeoIP DB update")
		return
	}
	url := fmt.Sprintf("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=%s&suffix=tar.gz", licenseKey)
	go func() {
		for {
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != 200 {
				log.Printf("GeoIP download failed: %v", err)
				time.Sleep(24 * time.Hour)
				continue
			}
			gz, _ := gzip.NewReader(resp.Body)
			tmpFile := "GeoLite2-Country.mmdb"
			f, _ := os.Create(tmpFile)
			io.Copy(f, gz)
			f.Close()
			s.dbMutex.Lock()
			s.db.Close()
			s.db, _ = geoip2.Open(tmpFile)
			s.dbMutex.Unlock()
			resp.Body.Close()
			log.Println("GeoIP DB updated")
			time.Sleep(7 * 24 * time.Hour)
		}
	}()
}
