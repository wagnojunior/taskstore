package models

type Store struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// type StoreService struct {
// 	DB *sql.DB
// }

// func (ss *StoreService) Create(name string) (*Store, error) {
// 	store := Store{
// 		Name: name,
// 	}

// 	row := ss.DB.QueryRow(`
// 		INSERT INTO stores (name)
// 		VALUES ($1) RETURNING id`,
// 		name)

// 	err := row.Scan(&store.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("create store: %w", err)
// 	}

// 	log.Println("Store created!")
// 	return &store, nil
// }
