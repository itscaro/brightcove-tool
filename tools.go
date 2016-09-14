package main

type Config struct {
	Token string
	Share []struct {
		ShareeAccountIds []int    `yaml:"sharee_account_ids"`
		Tags             []string `yaml:"tags"`
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
