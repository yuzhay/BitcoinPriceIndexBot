package main

var Config = struct {
	PriceIndexURI    string `required:"true"`
	TelegramBotToken string `required:"true"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `default:""`
		Port     uint   `default:"3306"`
	}
}{}
