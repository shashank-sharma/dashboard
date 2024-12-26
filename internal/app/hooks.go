package app

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func (app *Application) registerHooks() {
	app.Pb.OnRecordCreate("tokens").BindFunc(func(e *core.RecordEvent) error {
		if accessToken := e.Record.GetString("access_token"); accessToken != "" {
			encrypted, err := security.Encrypt([]byte(accessToken), app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				return fmt.Errorf("failed to encrypt access_token: %w", err)
			}
			e.Record.Set("access_token", encrypted)
		}

		// Encrypt refresh_token if present
		if refreshToken := e.Record.GetString("refresh_token"); refreshToken != "" {
			encrypted, err := security.Encrypt([]byte(refreshToken), app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				return fmt.Errorf("failed to encrypt refresh_token: %w", err)
			}
			e.Record.Set("refresh_token", encrypted)
		}

		return e.Next()
	})

	app.Pb.OnRecordViewRequest("tokens").BindFunc(func(e *core.RecordRequestEvent) error {
		if encryptedToken := e.Record.GetString("access_token"); encryptedToken != "" {
			decrypted, err := security.Decrypt(encryptedToken, app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				log.Printf("Failed to decrypt access_token: %v\n", err)
				e.Record.Set("access_token", "[decryption failed]")
			} else {
				e.Record.Set("access_token", decrypted)
			}
		}

		// Decrypt refresh_token if present
		if encryptedToken := e.Record.GetString("refresh_token"); encryptedToken != "" {
			decrypted, err := security.Decrypt(encryptedToken, app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				log.Printf("Failed to decrypt refresh_token: %v\n", err)
				e.Record.Set("refresh_token", "[decryption failed]")
			} else {
				e.Record.Set("refresh_token", decrypted)
			}
		}

		return e.Next()
	})

	// Handle token updates
	app.Pb.OnRecordUpdate("tokens").BindFunc(func(e *core.RecordEvent) error {
		// Encrypt access_token if it's being updated
		if accessToken := e.Record.GetString("access_token"); accessToken != "" {
			encrypted, err := security.Encrypt([]byte(accessToken), app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				return fmt.Errorf("failed to encrypt access_token: %w", err)
			}
			e.Record.Set("access_token", encrypted)
		}

		// Encrypt refresh_token if it's being updated
		if refreshToken := e.Record.GetString("refresh_token"); refreshToken != "" {
			encrypted, err := security.Encrypt([]byte(refreshToken), app.Pb.Store().Get("ENCRYPTION_KEY").(string))
			if err != nil {
				return fmt.Errorf("failed to encrypt refresh_token: %w", err)
			}
			e.Record.Set("refresh_token", encrypted)
		}

		return e.Next()
	})
}
