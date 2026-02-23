// Package service - LLM integration for optional AI-powered explanations.
// App works fully without LLM. If LLM_API_KEY exists, enhanced explanations are generated.
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/careermanifest/backend/internal/config"
	"github.com/careermanifest/backend/internal/dto"
)

// LLMService handles optional LLM integration (Groq Llama3 or Claude).
type LLMService struct {
	cfg    *config.Config
	client *http.Client
}

// NewLLMService creates a new LLMService.
func NewLLMService(cfg *config.Config) *LLMService {
	return &LLMService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// IsEnabled checks if LLM integration is configured.
func (s *LLMService) IsEnabled() bool {
	return s.cfg.IsLLMEnabled()
}

// GenerateExplanation produces an AI-powered personalized career explanation.
func (s *LLMService) GenerateExplanation(result *dto.AssessmentResult) (string, error) {
	if !s.IsEnabled() {
		return s.generateTemplateExplanation(result), nil
	}

	prompt := buildPrompt(result)

	switch strings.ToLower(s.cfg.LLMProvider) {
	case "groq":
		return s.callGroq(prompt)
	case "claude":
		return s.callClaude(prompt)
	default:
		return s.generateTemplateExplanation(result), nil
	}
}

// buildPrompt creates the prompt for the LLM.
func buildPrompt(result *dto.AssessmentResult) string {
	return fmt.Sprintf(`You are a career counselor specializing in Indian education and career paths.

Based on the following career assessment result, provide:
1. A personalized explanation (2-3 paragraphs) of why this career path suits the student
2. A detailed 1-year preparation plan with monthly milestones
3. An overview of relevant exam syllabi they should prepare for

Assessment Result:
- Best Career Path: %s
- Risk Level: %s (Score: %.1f)
- Top 3 Career Scores: %s (%.0f%%), %s (%.0f%%), %s (%.0f%%)

Keep the tone encouraging but realistic. Focus on actionable Indian-specific advice.
Include specific Indian exam names, colleges, and salary expectations in INR.
Format with clear headings and bullet points.`,
		result.BestCareerPath,
		result.Risk.Level, result.Risk.Score,
		result.Scores[0].Category, result.Scores[0].Percentage,
		result.Scores[1].Category, result.Scores[1].Percentage,
		result.Scores[2].Category, result.Scores[2].Percentage,
	)
}

// callGroq calls the Groq API (Llama3 compatible with OpenAI format).
func (s *LLMService) callGroq(prompt string) (string, error) {
	body := map[string]interface{}{
		"model": s.cfg.LLMModel,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a career counselor for Indian students."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+s.cfg.LLMApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("groq API call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse groq response: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from groq")
}

// callClaude calls the Anthropic Claude API.
func (s *LLMService) callClaude(prompt string) (string, error) {
	body := map[string]interface{}{
		"model":      s.cfg.LLMModel,
		"max_tokens": 2000,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("x-api-key", s.cfg.LLMApiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("claude API call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse claude response: %w", err)
	}

	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", fmt.Errorf("no response from claude")
}

// generateTemplateExplanation produces a structured explanation without LLM.
func (s *LLMService) generateTemplateExplanation(result *dto.AssessmentResult) string {
	return fmt.Sprintf(`## Your CareerManifest Analysis

### Recommended Path: %s

Based on your academic background, financial situation, personality traits, and career interests, **%s** emerges as your strongest career match with a compatibility score of %.0f%%.

### Risk Assessment: %s (Score: %.1f/10)
Your risk profile indicates a %s risk level. This means %s

### Why This Path?
Your responses indicate strong alignment with the skills, temperament, and goals required for success in %s. The scoring engine evaluated your answers across 6 major career categories, and this path scored highest based on weighted analysis of 30 factors.

### Next Steps
Follow the preparation roadmap provided below. Focus on building the required skills and preparing for the suggested exams. Remember, career decisions are personal — use this analysis as a guide, not a verdict.

*This analysis was generated by CareerManifest's rule-based scoring engine. For a more personalized AI-powered explanation, contact your institution about enabling the AI module.*`,
		result.BestCareerPath,
		result.BestCareerPath,
		result.Scores[0].Percentage,
		result.Risk.Level,
		result.Risk.Score,
		strings.ToLower(result.Risk.Level),
		getRiskExplanation(result.Risk.Level),
		result.BestCareerPath,
	)
}

// Chat handles a free-form chat question in the context of a career assessment result.
func (s *LLMService) Chat(message string, result *dto.AssessmentResult) (string, error) {
	if !s.IsEnabled() {
		return s.generateTemplateChatResponse(message, result), nil
	}

	prompt := fmt.Sprintf(`You are a friendly, knowledgeable career counselor chatbot for Indian students. The student has completed a career assessment on CareerManifest.

Their assessment results:
- Best Career Path: %s
- Risk Level: %s (Score: %.1f/10)
- Top Scores: %s (%.0f%%), %s (%.0f%%)

The student is asking: "%s"

Provide a helpful, concise, and encouraging response. Keep it under 250 words.
Focus on actionable Indian-specific advice (exams, colleges, salary in INR, timeline).
Use bullet points for clarity. Be warm and supportive.
If the question is unrelated to career/education, gently redirect them.`,
		result.BestCareerPath,
		result.Risk.Level, result.Risk.Score,
		result.Scores[0].Category, result.Scores[0].Percentage,
		result.Scores[1].Category, result.Scores[1].Percentage,
		message,
	)

	var reply string
	var err error
	switch strings.ToLower(s.cfg.LLMProvider) {
	case "groq":
		reply, err = s.callGroq(prompt)
	case "claude":
		reply, err = s.callClaude(prompt)
	default:
		return s.generateTemplateChatResponse(message, result), nil
	}

	// Fallback to template if AI fails
	if err != nil {
		return s.generateTemplateChatResponse(message, result), nil
	}
	return reply, nil
}

// generateTemplateChatResponse produces a helpful response without LLM.
func (s *LLMService) generateTemplateChatResponse(message string, result *dto.AssessmentResult) string {
	lowerMsg := strings.ToLower(message)
	career := result.BestCareerPath

	switch {
	case strings.Contains(lowerMsg, "exam") || strings.Contains(lowerMsg, "prepare"):
		return fmt.Sprintf("For %s, here are the key exams you should focus on:\n\nYour top career match is %s with a %.0f%% compatibility score. Check the 'Exams to Prepare For' section in your results above for specific exam recommendations. The preparation roadmap also has a step-by-step timeline.\n\nWould you like to know more about a specific exam?", career, career, result.Scores[0].Percentage)
	case strings.Contains(lowerMsg, "salary") || strings.Contains(lowerMsg, "earning") || strings.Contains(lowerMsg, "pay") || strings.Contains(lowerMsg, "money"):
		return fmt.Sprintf("Great question about earnings! Check the '5-Year Salary Projection' section in your results above for detailed year-by-year salary expectations for %s.\n\nRemember, actual salaries depend on your skills, location, and the specific company/role. The figures shown are average ranges for Indian professionals.", career)
	case strings.Contains(lowerMsg, "college") || strings.Contains(lowerMsg, "university") || strings.Contains(lowerMsg, "institute"):
		return fmt.Sprintf("For %s, see the 'Suggested Institutions' section in your results above for top recommended colleges.\n\nI'd recommend researching each institution's placement records, faculty, and alumni network to find the best fit for you.", career)
	case strings.Contains(lowerMsg, "risk") || strings.Contains(lowerMsg, "safe"):
		return fmt.Sprintf("Your risk level is %s (Score: %.1f/10).\n\n%s\n\nCheck the Risk Assessment section above for a detailed breakdown of the factors.", result.Risk.Level, result.Risk.Score, getRiskExplanation(result.Risk.Level))
	case strings.Contains(lowerMsg, "skill") || strings.Contains(lowerMsg, "learn"):
		return fmt.Sprintf("For %s, check the 'Skills You Need' section in your results above.\n\nI recommend starting with the most fundamental skills first, then building up to advanced ones. The preparation roadmap gives you a timeline for when to learn each skill.", career)
	case strings.Contains(lowerMsg, "start") || strings.Contains(lowerMsg, "begin") || strings.Contains(lowerMsg, "next step") || strings.Contains(lowerMsg, "today"):
		return fmt.Sprintf("Here's how to start your journey towards %s:\n\n1. Check the Preparation Roadmap in your results — it has a step-by-step plan\n2. Start with Step 1 today\n3. Build the skills listed in 'Skills You Need'\n4. Register for the exams listed in your results\n\nConsistency is key. Even 1-2 hours daily can make a huge difference over 6-12 months!", career)
	default:
		return fmt.Sprintf("Thanks for your question! Based on your assessment, %s is your best career match with %.0f%% compatibility.\n\nYour results page above has detailed information about:\n• 📊 Score breakdown across all careers\n• 🗺️ Step-by-step preparation roadmap\n• 🛠️ Skills you need to build\n• 📝 Exams to prepare for\n• 🏫 Suggested institutions\n• 💰 Salary projections\n\nFeel free to ask about any specific topic!", career, result.Scores[0].Percentage)
	}
}

func getRiskExplanation(level string) string {
	switch level {
	case "Low":
		return "you have a stable foundation to pursue this career path with moderate pace. You can afford to take calculated risks in your career planning."
	case "Medium":
		return "you should balance ambition with pragmatism. Consider having a backup plan while pursuing your primary career goal."
	case "High":
		return "financial or family pressures require careful planning. Prioritize paths that offer quicker returns while keeping long-term goals in sight."
	default:
		return "consider consulting a career counselor for personalized guidance."
	}
}
