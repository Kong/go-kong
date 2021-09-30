package kong

type auth struct {
	keyAuth *keyAuthOption
}

type keyAuthOption struct {
	key   string
	value string
}
