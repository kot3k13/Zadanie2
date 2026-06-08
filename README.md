# Technologie Chmurowe - Zadanie 2
**Autor:** Kacper Kot

## 1. Opis Architektury CI/CD
W ramach zadania opracowano w pełni zautomatyzowany łańcuch CI/CD w usłudze GitHub Actions. Łańcuch ten buduje, testuje pod kątem bezpieczeństwa i publikuje zoptymalizowany obraz kontenera z aplikacją pogodową napisaną w języku Go.

Zgodnie z wymaganiami proces spełnia następujące warunki:
* **Wieloplatformowość:** Obraz budowany jest ze wsparciem dla dwóch architektur: `linux/amd64` oraz `linux/arm64` poprzez użycie QEMU oraz silnika Docker Buildx.
* **Zarządzanie Pamięcią Podręczną (Cache):** Dane eksportowane są do publicznego repozytorium na DockerHub (`weather-cache`). Wykorzystano backend cache w trybie `mode=max`, co zapewnia zapisywanie wszystkich warstw budowania, maksymalizując optymalizację przy kolejnych kompilacjach.
* **Testy Bezpieczeństwa (CVE):** Do skanowania użyto skanera **Trivy**. Ustawiono twardy blok (`exit-code: 1`), który przerywa łańcuch po wykryciu luk klasy `CRITICAL` lub `HIGH`. Ponadto, zgodnie z wytycznymi, skanowany jest **docelowy obraz** kontenera (budowany lokalnie), a dopiero po uzyskaniu wyniku pozytywnego proces przechodzi do budowy multi-architekturowej i przesyła obraz do GHCR.

## 2. Uzasadnienie Strategii Tagowania
Do profesjonalnego zarządzania wersjami użyto zautomatyzowanej akcji `docker/metadata-action`. Przyjęto następujący hybrydowy system tagowania:
1. **Tagowanie bazujące na SHA Git (`type=sha`):** Przy każdym wdrożeniu do gałęzi głównej (`main`) generowany jest tag oparty na kalce kryptograficznej commita (np. `sha-a1b2c3d`). Gwarantuje to pełną identyfikowalność i pozwala precyzyjnie powiązać kontener na produkcji z konkretnym kodem źródłowym.
2. **Tagowanie SemVer (`type=semver`):** Gdy w repozytorium utworzony zostanie tag wersji (np. `v1.0.0`), obraz otrzyma ładny, oficjalny tag wydania zgodny z regułami Semantic Versioning. Rozdziela to techniczne wersje robocze (SHA) od stabilnych, docelowych wersji produkcyjnych dla użytkowników.

## 3. Optymalizacje Bezpieczeństwa (Dockerfile)
Aby środowisko DevOps było w 100% bezpieczne, plik `Dockerfile` wykorzystuje:
- Model **Multi-stage build** z wykorzystaniem obrazu bazowego `scratch` (brak shella i narzutu systemu operacyjnego zmniejsza wektor ataków niemal do zera).
- Zasadę najmniejszych uprawnień (**Least Privilege**). Użytkownik root służy tylko do zbudowania binarki, aplikacja uruchamiana jest z poziomu nieuprzywilejowanego konta (ID `10001`).