package seeders

import (
	"log"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	var userCount int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := models.User{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Println("Admin user seeded: admin@example.com / password")
	}

	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		desc1 := "High performance cloud hosting with 99.99% uptime SLA and instant provisioning."
		desc2 := "Custom domain name registration including free WHOIS privacy protection."
		desc3 := "Industry-standard 256-bit SSL encryption certificate for web applications."
		desc4 := "Automated email newsletters, campaigns, and subscriber list management."
		desc5 := "Dedicated PostgreSQL database server with automated backups and read replicas."
		desc6 := "Ultra-low latency API gateway with rate limiting and analytics."

		products := []models.Product{
			{Name: "Cloud Hosting Plan", Slug: "cloud-hosting-plan", Description: &desc1, Price: 29.99, Status: true},
			{Name: "Domain Name Registration", Slug: "domain-name-registration", Description: &desc2, Price: 14.50, Status: true},
			{Name: "Enterprise SSL Certificate", Slug: "enterprise-ssl-certificate", Description: &desc3, Price: 79.00, Status: true},
			{Name: "Email Marketing Suite", Slug: "email-marketing-suite", Description: &desc4, Price: 45.00, Status: false},
			{Name: "Dedicated Database Server", Slug: "dedicated-database-server", Description: &desc5, Price: 120.00, Status: true},
			{Name: "API Gateway Pro", Slug: "api-gateway-pro", Description: &desc6, Price: 55.00, Status: true},
		}

		for _, p := range products {
			if err := db.Create(&p).Error; err != nil {
				return err
			}
		}
		log.Printf("Seeded %d products\n", len(products))
	}

	return nil
}
