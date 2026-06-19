package components

import "os"

func activeClass(current, link string) string {
	if current == link {
		return "active"
	}
	return ""
}

func FormspreeID() string {
	return os.Getenv("FORMSPREE_FORM_ID")
}
