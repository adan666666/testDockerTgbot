package rebot

type Conf struct {
	System struct {
		Name string `yaml:"name"`
	} `yaml:"system"`

	TgBot struct {
		Token  string `yaml:"Token"`
		Hour   int    `yaml:"hour"`
		Min    int    `yaml:"min"`
		Sec    int    `yaml:"sec"`
		ChatID int64  `yaml:"chatID"`
	}
}
