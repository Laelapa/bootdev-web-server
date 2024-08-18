package main

import "strings"

func profanityFilter(strToFilter string, badword string) string {
	s := strings.Split(strToFilter, " ")
	for i, str := range s {
		if strings.ToLower(str) == badword {
			s[i] = strings.Replace(strings.ToLower(str), badword, "****", -1)
		}
	}

	return strings.Join(s, " ")

}
