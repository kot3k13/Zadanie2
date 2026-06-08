# Budowanie środowiska
FROM golang:alpine AS builder

# Instalacja certyfikatów SSL oraz utworzenie bezpiecznego użytkownika non-root
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' -u 10001 appuser

WORKDIR /app

# Inicjalizacja modułu Go
RUN go mod init weatherapp

# Skopiowanie kodów źródłowych
COPY main.go index.html ./

# Kompilacja odchudzonej binarki
# CGO_ENABLED=0 pozwala na uruchomienie w scratch
# flaga -ldflags="-w -s" usuwa informacje debugowania
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o weatherapp .

FROM scratch

# Etykieta zgodna ze standardem OCI 
LABEL org.opencontainers.image.authors="Kacper Kot" \
      org.opencontainers.image.description="Aplikacja pogodowa w Go" \
      org.opencontainers.image.source="https://github.com/kot3k13/Zadanie2"

# Przeniesienie pliku konfiguracyjnego użytkownika
COPY --from=builder /etc/passwd /etc/passwd
# Skopiowanie certyfikatów SSL do autoryzacji HTTPS
COPY --from=builder --chown=10001:10001 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Skopiowanie gotowej aplikacji
COPY --from=builder --chown=10001:10001 /app/weatherapp /weatherapp
# Zmiana użytkownika na non-root przed uruchomieniem aplikacji
USER 10001

EXPOSE 8080

# Wymagany HEALTHCHECK z wykorzystaniem wbudowanej w program funkcji testującej
HEALTHCHECK --interval=30s --timeout=3s \
  CMD ["/weatherapp", "-health"]

# Uruchomienie aplikacji
CMD ["/weatherapp"]