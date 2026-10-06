package session

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/config"
)

func init() {
	gob.Register(gonertia.Flash{})
	gob.Register(gonertia.ValidationErrors{})
	gob.Register(map[string]any{})
	gob.Register(map[string]string{})
	gob.Register(uint(0))
}

type Manager struct {
	*scs.SessionManager
}

func NewManager(cfg *config.Config) *Manager {
	sm := scs.New()
	sm.Lifetime = 24 * time.Hour
	sm.Cookie.Name = "starter_kit_session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Path = "/"
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = cfg.IsProduction()

	return &Manager{SessionManager: sm}
}

type FlashProvider struct {
	sm *scs.SessionManager
}

func NewFlashProvider(sm *scs.SessionManager) *FlashProvider {
	return &FlashProvider{sm: sm}
}

const (
	flashErrorsKey       = "_gonertia_errors"
	flashMessagesKey     = "_gonertia_flash"
	flashClearHistoryKey = "_gonertia_clear_history"
)

func (fp *FlashProvider) FlashErrors(ctx context.Context, errors gonertia.ValidationErrors) error {
	fp.sm.Put(ctx, flashErrorsKey, errors)
	return nil
}

func (fp *FlashProvider) GetErrors(ctx context.Context) (gonertia.ValidationErrors, error) {
	val := fp.sm.Pop(ctx, flashErrorsKey)
	if val == nil {
		return nil, nil
	}
	if errs, ok := val.(gonertia.ValidationErrors); ok {
		return errs, nil
	}
	if m, ok := val.(map[string]any); ok {
		return gonertia.ValidationErrors(m), nil
	}
	return nil, nil
}

func (fp *FlashProvider) Flash(ctx context.Context, flash gonertia.Flash) error {
	fp.sm.Put(ctx, flashMessagesKey, flash)
	return nil
}

func (fp *FlashProvider) GetFlash(ctx context.Context) (gonertia.Flash, error) {
	val := fp.sm.Pop(ctx, flashMessagesKey)
	if val == nil {
		return nil, nil
	}
	if f, ok := val.(gonertia.Flash); ok {
		return f, nil
	}
	if m, ok := val.(map[string]any); ok {
		return gonertia.Flash(m), nil
	}
	return nil, nil
}

func (fp *FlashProvider) FlashClearHistory(ctx context.Context) error {
	fp.sm.Put(ctx, flashClearHistoryKey, true)
	return nil
}

func (fp *FlashProvider) ShouldClearHistory(ctx context.Context) (bool, error) {
	return fp.sm.PopBool(ctx, flashClearHistoryKey), nil
}
