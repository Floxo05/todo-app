# ToDo-Projekt Check24

## Komponenten

1. Go Api: Dies ist der Backend-Server, der in Go geschrieben ist. Er stellt die Haupt-API für die Interaktion mit der Datenbank bereit.
2. React Frontend: Dies ist die Benutzeroberfläche, die in React geschrieben ist. Es kommuniziert mit der Go-API, um Daten anzuzeigen und zu speichern.
3. MariaDB: Dies ist die Datenbank, in der alle Daten gespeichert werden.
4. Adminer: Dies ist ein Datenbankverwaltungstool, das verwendet wird, um auf die Datenbank zuzugreifen.

## Installation 

### Mit Docker

1. Stellen Sie sicher, dass Docker und Docker Compose auf Ihrem System installiert sind.
2. Klonen Sie dieses Repository.
3. Navigieren Sie in das Projektverzeichnis.
4. Führen Sie `docker-compose up` aus.

Nun sollten die Container über folgende Ports erreichbar sein:
- Go API: http://localhost:8080
- React Frontend: http://localhost:80
- Adminer: http://localhost:8081
- MariaDB: http://localhost:3306


### Ohne Docker

1. Stellen Sie sicher, dass Go und Node.js auf Ihrem System installiert sind.
2. Stellen Sie sicher, dass MariaDB auf Ihrem System installiert ist.
3. Klonen Sie dieses Repository.
4. Starten Sie den Go-Server
   1. Navigieren Sie in das `goApi`-Verzeichnis.
   2. Kopieren Sie die `.env.test`-Datei und benennen Sie sie in `.env` um. Passen Sie die Werte in der `.env`-Datei an.
   3. Installieren Sie die Abhängigkeiten, indem Sie `go mod download` ausführen.
   4. Führen Sie die Migrationen aus, indem Sie `go run cmd/migrate/up/up.go` ausführen.
   5. Starten Sie den Server, indem Sie `go run cmd/api/main.go` ausführen.
5. Starten Sie das React-Frontend
   1. Navigieren Sie in das `todo-react-frontend`-Verzeichnis.
   2. Installieren Sie die Abhängigkeiten, indem Sie `npm install` ausführen.
   3. Kopieren Sie die `.env.test`-Datei und benennen Sie sie in `.env` um. Passen Sie die Werte in der `.env`-Datei an.
   4. Starten Sie den Server, indem Sie `npm start` ausführen.

Nun sollten die Anwendungen über folgende Ports erreichbar sein:
- Go API: http://localhost:8080
- React Frontend: http://localhost:3000
