package main

import "snibbage.borkware.com/internal/models"

type templateData struct {
	Snippet models.Snippet 
	Snippets []models.Snippet
}
