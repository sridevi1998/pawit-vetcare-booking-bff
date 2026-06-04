package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("PawIt Booking BFF started", "port", port)
	if err := http.ListenAndServe(":"+port, newServer()); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func newServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", bookingPage)
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/v1/public/slots", slots)
	return mux
}

func bookingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Book a PawIt VetCare Visit</title>
  <style>
    body { margin: 0; font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; color: #0f172a; background: #f8fafc; }
    main { max-width: 920px; margin: 0 auto; padding: 48px 20px; }
    h1 { margin: 0; font-size: clamp(2rem, 6vw, 4rem); line-height: 1; }
    p { color: #475569; font-size: 1.05rem; line-height: 1.7; }
    .panel { margin-top: 32px; border: 1px solid #e2e8f0; border-radius: 8px; background: white; overflow: hidden; }
    .slot { display: grid; gap: 8px; padding: 20px; border-top: 1px solid #e2e8f0; }
    .slot:first-child { border-top: 0; }
    .meta { color: #2563eb; font-size: .85rem; font-weight: 700; text-transform: uppercase; }
    .name { font-size: 1.1rem; font-weight: 800; }
    button { width: fit-content; border: 0; border-radius: 8px; background: #0f172a; color: white; padding: 10px 14px; font-weight: 700; }
  </style>
</head>
<body>
  <main>
    <h1>Book a PawIt VetCare visit</h1>
    <p>Choose an available appointment window. Online booking confirmation is being connected to clinic scheduling; for now, these are the public slots exposed by the booking boundary.</p>
    <section class="panel" id="slots" aria-live="polite"></section>
  </main>
  <script>
    const target = document.getElementById("slots");
    fetch("/api/v1/public/slots")
      .then((response) => response.json())
      .then((data) => {
        target.innerHTML = data.items.map((slot) => ` + "`" + `
          <article class="slot">
            <div class="meta">${slot.type} &middot; ${slot.time}</div>
            <div class="name">${slot.clinic}</div>
            <p>${slot.vet}</p>
            <button type="button">Request slot</button>
          </article>
        ` + "`" + `).join("");
      })
      .catch(() => {
        target.innerHTML = '<article class="slot"><div class="name">Slots are unavailable</div><p>Please try again shortly.</p></article>';
      });
  </script>
</body>
</html>`))
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "pawit-booking-bff"})
}

func slots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": []map[string]string{
		{"id": "slot_001", "clinic": "PawIt VetCare", "vet": "Dr. Asha Rao", "time": "09:30", "type": "Check-up"},
	}})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
