package url

import (
	"encoding/json"
	"fmt"
	"os"
)

type urlRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// save сохраняет все сокращённые URL в JSON-файл.
func (r *Repository) save() error {
	records := make([]urlRecord, 0, len(r.urls))

	for id, originalURL := range r.urls {
		records = append(records, urlRecord{
			UUID:        id,
			ShortURL:    id,
			OriginalURL: originalURL,
		})
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := r.filePath + ".tmp"

	if err = os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("запись во временный файл: %w", err)
	}

	if err = os.Rename(tmpPath, r.filePath); err != nil {
		return fmt.Errorf("работа с основным файлом: %w", err)
	}

	return nil
}

// load загружает сохранённые URL из JSON-файла при запуске сервера.
func (r *Repository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	var records []urlRecord

	if err = json.Unmarshal(data, &records); err != nil {
		return err
	}

	for _, record := range records {
		r.urls[record.ShortURL] = record.OriginalURL
	}

	return nil
}
