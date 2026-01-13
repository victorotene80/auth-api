package policy

import "time"

type PasswordPolicy struct {
	MinLength          int
	RequireUppercase   bool
	RequireLowercase   bool
	RequireNumbers     bool
	RequireSpecialChar bool
	MaxAge             time.Duration 
	PreventReuse       int          
}

func DefaultPasswordPolicy() PasswordPolicy {
	return PasswordPolicy{
		MinLength:          8,
		RequireUppercase:   true,
		RequireLowercase:   true,
		RequireNumbers:     true,
		RequireSpecialChar: true,
		MaxAge:             90 * 24 * time.Hour,
		PreventReuse:       5,
	}
}
