package main

// UTIL
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
