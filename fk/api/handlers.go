package api

import (
	"fmt"
	"net/http"
)

// Index write index info
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		------------------------------------------------------
		F**K TOOLS REST API:
		------------------------------------------------------
		Metrics: 
			GET: /api/v1/metrics
		-------------------------------
		Rules:
			GET: /api/v1/rules	
			POST: /api/v1/rules
			PUT: /api/v1/rules/{rule id}/
			DELTE: /api/v1/rules/{rule id}
		-------------------------------
		Issues:
			GET: /api/v1/issues
			GET: /api/v1/issues?q={query}
		------------------------------------------------------
	`)
}
