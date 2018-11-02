package structures

type Publisher struct {
	Name string `json:"name"valid:"required,length(1|150)"`
}

type NotificationsCount struct {
	Total int
	PublisherId int
}