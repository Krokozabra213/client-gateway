package constants

type TelegramMarkup string

const (
	StartupMarkup          TelegramMarkup = "startup"
	RegisterMarkup         TelegramMarkup = "register"
	GenderMarkup           TelegramMarkup = "gender_markup"
	BirthdayMarkup         TelegramMarkup = "birthday_markup"
	RegisterSuccessMarkup  TelegramMarkup = "register_success"
	MainMenuMarkup         TelegramMarkup = "main_menu"
	MainMenuWithBotsMarkup TelegramMarkup = "main_menu_with_bots"

	ConnectBotConfirmMarkup TelegramMarkup = "connect_bot_confirm"
	BalanceMarkup           TelegramMarkup = "balance"

	DepositMarkup TelegramMarkup = "deposit"

	ShopMarkup TelegramMarkup = "shop"
)
