package intelligence

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubex-ecosystem/gemx/ghbex/internal/defs/gromptz"
	gl "github.com/kubex-ecosystem/gemx/ghbex/internal/module/logger"
)

// setCachedHealthStatus armazena status no cache de forma thread-safe
func setCachedHealthStatus(providerName string, isHealthy bool) {
	// Simplificado por enquanto
}

// performHealthCheck executa a verificação real baseada no tipo de provider
func performHealthCheck(ctx context.Context, provider gromptz.Provider) bool {
	// Health check baseado no tipo de provider
	switch provider.Name() {
	case "gemini":
		return checkGeminiHealth(ctx, provider)
	case "ollama":
		return checkOllamaHealth(ctx, provider)
	case "openai", "chatgpt":
		return checkOpenAIHealth(ctx, provider)
	case "claude":
		return checkClaudeHealth(ctx, provider)
	case "deepseek":
		return checkDeepSeekHealth(ctx, provider)
	default:
		// Para providers desconhecidos, assumir disponível se tem API key
		return provider.IsAvailable()
	}
}

// checkGeminiHealth verifica especificamente o Gemini 2.5 Flash
func checkGeminiHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "Gemini provider not available (no API key)")
		return false
	}

	// Teste mínimo e rápido para verificar conectividade
	testPrompt := "ping"

	// Use um canal para timeout rápido
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Gemini health check panic: %v", r))
				done <- false
			}
		}()

		// Teste simples - apenas verifica se consegue fazer uma chamada básica
		// Para Gemini 2.5 Flash, o response deve ser rápido
		response, err := provider.Execute(testPrompt)

		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Gemini health check failed: %v", err))
		}

		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", fmt.Sprintf("%s %s health check passed", strings.ToTitle(provider.Name()), provider.Version()))
		} else {
			gl.Log("warn", fmt.Sprintf("%s %s health check failed", strings.ToTitle(provider.Name()), provider.Version()))
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Gemini health check timeout (may be slow or unavailable)")
		return false
	}
}

// checkOllamaHealth verifica se o Ollama está rodando localmente
func checkOllamaHealth(ctx context.Context, provider gromptz.Provider) bool {
	// Ollama muitas vezes está configurado mas não rodando
	if !provider.IsAvailable() {
		gl.Log("debug", "Ollama provider not available (not configured)")
		return false
	}

	// Para Ollama, fazer um teste mais rigoroso já que é local
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Ollama health check panic: %v", r))
				done <- false
			}
		}()

		// Teste básico para ver se o Ollama responde
		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Ollama health check failed (may not be running): %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "Ollama health check passed (server running)")
		} else {
			gl.Log("warn", "Ollama health check failed (server may not be running)")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Ollama health check timeout (server may be slow)")
		return false
	}
}

// checkOpenAIHealth verifica OpenAI API
func checkOpenAIHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "OpenAI provider not available (no API key)")
		return false
	}

	// Similar ao Gemini, mas com especificidades do OpenAI
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("OpenAI health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("OpenAI health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "OpenAI health check passed")
		} else {
			gl.Log("warn", "OpenAI health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "OpenAI health check timeout")
		return false
	}
}

// checkClaudeHealth verifica Anthropic Claude
func checkClaudeHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "Claude provider not available (no API key)")
		return false
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Claude health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Claude health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "Claude health check passed")
		} else {
			gl.Log("warn", "Claude health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Claude health check timeout")
		return false
	}
}

// checkDeepSeekHealth verifica DeepSeek API
func checkDeepSeekHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "DeepSeek provider not available (no API key)")
		return false
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("DeepSeek health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("DeepSeek health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "DeepSeek health check passed")
		} else {
			gl.Log("warn", "DeepSeek health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "DeepSeek health check timeout")
		return false
	}
}
