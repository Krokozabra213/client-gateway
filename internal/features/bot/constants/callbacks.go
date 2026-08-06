package constants

type TelegramCallback string

const (
	RegisterCallback          TelegramCallback = "register"
	ApproveRegisterCallback   TelegramCallback = "approve_register"
	CancelRegisterCallback    TelegramCallback = "cancel_register"
	MaleCallback              TelegramCallback = "male"
	FemaleCallback            TelegramCallback = "female"
	ShowMainMenuCallback      TelegramCallback = "show_main_menu"
	ShowHelpCallback          TelegramCallback = "show_help"
	ConnectBotCallback        TelegramCallback = "connect_bot"
	ManagementCallback        TelegramCallback = "management"
	StatisticsCallback        TelegramCallback = "statistics"
	BalanceCallback           TelegramCallback = "balance"
	BuySubscriptionCallback   TelegramCallback = "buy_subscription"
	HelpCallback              TelegramCallback = "help"
	ConnectBotConfirmCallback TelegramCallback = "connect_bot_confirm"

	MyBotsCallback  TelegramCallback = "my_bots"
	BotInfoCallback TelegramCallback = "bot_info"

	DepositCallback            TelegramCallback = "deposit"
	TGDepositCallback          TelegramCallback = "tg_deposit"
	TbankDepositCallback       TelegramCallback = "tbank_deposit"
	TGDepositAmountCallback    TelegramCallback = "tg_deposit_amount"
	TbankDepositAmountCallback TelegramCallback = "tbank_deposit_amount"

	MyOrdersCallback        TelegramCallback = "my_orders"
	MySubscriptionsCallback TelegramCallback = "my_subscriptions"
	ShopCallback            TelegramCallback = "shop"
	ShopItemCallback        TelegramCallback = "shop_item"
	BuyProductCallback      TelegramCallback = "buy_product"
)
