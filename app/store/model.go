package store

type Player struct {
	PlayerId       string  `json:"playerId" gorm:"primaryKey;unique;<-:create;column:player_id;size:50"`
	BankName       string  `json:"bankName" gorm:"size:50;not null;default:'';column:bank_name"`
	AccountNumber  string  `json:"accountNumber" gorm:"size:50;not null;default:'';column:account_number"`
	AccountName    string  `json:"accountName" gorm:"size:50;not null;default:'';column:account_name"`
	Wallet         float64 `json:"wallet" gorm:"not null;default:0;column:wallet"`
	HashedPassword string  `json:"hashedPassword" gorm:"not null;default:'';column:hashed_password"`
}

type GetAllPlayersRequest struct {
	PlayerId         string
	AccountName      string
	AccountNumber    string
	BankName         string
	RemainingBalance float64
}

type AddBankAccountRequest struct {
	PlayerId      string `json:"playerId"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

type LoginRequest struct {
	PlayerId string
	JwtToken string
}
