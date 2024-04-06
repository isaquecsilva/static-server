package utils

type ValidationActions map[string]func(arg ...any) (error, int)