package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"hustlex/internal/config"
	"hustlex/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database holds database connections
type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Connect to PostgreSQL
	db, err := connectPostgres(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Connect to Redis
	rdb, err := connectRedis(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Database{
		DB:    db,
		Redis: rdb,
	}, nil
}

// connectPostgres establishes a PostgreSQL connection
func connectPostgres(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Configure GORM logger
	logLevel := logger.Silent
	if cfg.SSLMode == "disable" { // likely development
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}

// connectRedis establishes a Redis connection
func connectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Println("✅ Connected to Redis")
	return rdb, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	log.Println("🔄 Running database migrations...")

	// Enable UUID extension
	d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")

	// Auto migrate models
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Skill{},
		&models.UserSkill{},
		&models.Gig{},
		&models.GigProposal{},
		&models.GigContract{},
		&models.GigReview{},
		&models.SavingsCircle{},
		&models.CircleMember{},
		&models.Contribution{},
		&models.Wallet{},
		&models.Transaction{},
		&models.CreditScore{},
		&models.Loan{},
		&models.LoanRepayment{},
		&models.Course{},
		&models.Enrollment{},
		&models.OTPCode{},
		&models.Notification{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Create indexes
	d.createIndexes()

	// Seed initial data
	d.seedData()

	log.Println("✅ Database migrations completed")
	return nil
}

// createIndexes creates additional indexes for performance
func (d *Database) createIndexes() {
	// Composite indexes for common queries
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_gigs_status_created 
		ON gigs(status, created_at DESC);
	`)
	
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_gigs_category_status 
		ON gigs(category, status);
	`)
	
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_wallet_created 
		ON transactions(wallet_id, created_at DESC);
	`)
	
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_contributions_circle_round 
		ON contributions(circle_id, round);
	`)
	
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_notifications_user_unread 
		ON notifications(user_id, is_read) WHERE is_read = false;
	`)

	// Full-text search index for gigs
	d.DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_gigs_search 
		ON gigs USING gin(to_tsvector('english', title || ' ' || description));
	`)
}

// seedData seeds initial data
func (d *Database) seedData() {
	// Check if skills already exist
	var count int64
	d.DB.Model(&models.Skill{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("🌱 Seeding initial data...")

	// Seed skill categories
	skills := []models.Skill{
		// Digital Services
		{Name: "Graphic Design", Category: "digital", Description: "Logo design, branding, social media graphics", Icon: "palette"},
		{Name: "Content Writing", Category: "digital", Description: "Blog posts, copywriting, academic writing", Icon: "edit"},
		{Name: "Digital Marketing", Category: "digital", Description: "Social media management, ads, SEO", Icon: "trending_up"},
		{Name: "Video Editing", Category: "digital", Description: "YouTube, TikTok, promotional videos", Icon: "movie"},
		{Name: "Web Development", Category: "digital", Description: "Website design and development", Icon: "code"},
		{Name: "Mobile Development", Category: "digital", Description: "iOS and Android app development", Icon: "smartphone"},
		{Name: "UI/UX Design", Category: "digital", Description: "User interface and experience design", Icon: "design_services"},
		{Name: "Virtual Assistance", Category: "digital", Description: "Admin support, scheduling, data entry", Icon: "support_agent"},
		{Name: "Data Entry", Category: "digital", Description: "Data input and management", Icon: "keyboard"},
		{Name: "Transcription", Category: "digital", Description: "Audio to text transcription", Icon: "hearing"},
		
		// Creative Services
		{Name: "Photography", Category: "creative", Description: "Events, portraits, product photography", Icon: "camera_alt"},
		{Name: "Videography", Category: "creative", Description: "Event coverage, promotional videos", Icon: "videocam"},
		{Name: "Music Production", Category: "creative", Description: "Beat making, mixing, mastering", Icon: "music_note"},
		{Name: "Voice Over", Category: "creative", Description: "Narration, ads, podcasts", Icon: "mic"},
		{Name: "Animation", Category: "creative", Description: "2D/3D animation, motion graphics", Icon: "animation"},
		
		// Physical Services
		{Name: "Event Planning", Category: "physical", Description: "Parties, weddings, corporate events", Icon: "event"},
		{Name: "Tutoring", Category: "physical", Description: "Academic and skill tutoring", Icon: "school"},
		{Name: "Beauty Services", Category: "physical", Description: "Makeup, hair styling", Icon: "face"},
		{Name: "Fashion/Tailoring", Category: "physical", Description: "Custom clothing, alterations", Icon: "checkroom"},
		{Name: "Delivery/Errands", Category: "physical", Description: "Package delivery, personal errands", Icon: "local_shipping"},
		{Name: "Home Cleaning", Category: "physical", Description: "Residential cleaning services", Icon: "cleaning_services"},
		{Name: "Handyman", Category: "physical", Description: "Home repairs and maintenance", Icon: "build"},
		
		// Professional Services
		{Name: "Accounting", Category: "professional", Description: "Bookkeeping, tax preparation", Icon: "account_balance"},
		{Name: "Legal Services", Category: "professional", Description: "Document preparation, consultation", Icon: "gavel"},
		{Name: "Business Consulting", Category: "professional", Description: "Strategy, planning, advisory", Icon: "business_center"},
		{Name: "Translation", Category: "professional", Description: "Document and live translation", Icon: "translate"},
		{Name: "HR Services", Category: "professional", Description: "Recruitment, training", Icon: "people"},
	}

	for _, skill := range skills {
		skill.IsActive = true
		d.DB.Create(&skill)
	}

	log.Println("✅ Initial data seeded")
}

// Close closes all database connections
func (d *Database) Close() error {
	// Close Redis
	if err := d.Redis.Close(); err != nil {
		return err
	}

	// Close PostgreSQL
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks database connections
func (d *Database) HealthCheck(ctx context.Context) error {
	// Check PostgreSQL
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("postgres error: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres ping failed: %w", err)
	}

	// Check Redis
	if err := d.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}
