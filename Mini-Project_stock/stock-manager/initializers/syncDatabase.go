package initializers

import "example.com/stock-manager/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Token{})

	// Create foreign keys
	DB.Debug().Exec("ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories(id);")
	DB.Debug().Exec("ALTER TABLE tokens ADD FOREIGN KEY (user_id) REFERENCES users(id);")
}
