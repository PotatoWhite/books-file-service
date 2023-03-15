package util

import "strconv"

// AtoIOrNil converts a string to an int or returns nil if the string is empty or nil.
func AtoIOrNil(s *string) *int {
	if s == nil || *s == "" {
		return nil
	}
	i, err := strconv.Atoi(*s)
	if err != nil {
		return nil
	}
	return &i
}

func AtoUIOrNil(s *string) *uint {
	if s == nil || *s == "" {
		return nil
	}
	i, err := strconv.Atoi(*s)
	if err != nil {
		return nil
	}
	ui := uint(i)
	return &ui
}

func ItoAOrNil(id *int) *string {
	if id == nil {
		return nil
	}
	s := strconv.Itoa(*id)
	return &s
}

func UItoAOrNil(id *uint) *string {
	if id == nil {
		return nil
	}
	s := strconv.Itoa(int(*id))
	return &s
}
