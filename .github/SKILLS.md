# Cadence as an Agent Skill

Cadence is compatible with the [open agent skills ecosystem](https://skills.sh/). This allows AI agents and tools to integrate Cadence's AI-detection capabilities.

## Skill Definition

Cadence provides two main skills:

### 1. `analyze-repository`
Analyzes an entire Git repository for suspicious AI-generated commits.

**Input:**
```json
{
  "repository": "https://github.com/user/repo",
  "branch": "main",
  "output_format": "json"
}
```

**Output:**
```json
{
  "suspicious_commits": [
    {
      "hash": "abc1234",
      "confidence": 0.85,
      "reasons": ["High velocity", "Large commit"],
      "ai_analysis": "likely AI-generated"
    }
  ],
  "statistics": {
    "total_commits": 150,
    "total_suspicious": 3
  }
}
```

### 2. `detect-suspicious-commit`
Analyzes a single commit for AI-generation patterns.

**Input:**
```json
{
  "commit_hash": "abc1234",
  "additions": "def hello_world():\n    print('Hello')\n    return True",
  "use_ai": true
}
```

**Output:**
```json
{
  "suspicious": true,
  "confidence": 0.75,
  "reasons": ["High velocity", "Unusual pattern"],
  "ai_analysis": "possibly AI-generated (confidence: 75%)"
}
```

## Installing as a Skill

```bash
npx skills add CodeMeAPixel/Cadence
```

Or specify a version:
```bash
npx skills add CodeMeAPixel/Cadence@0.1.0
```

## Integration Examples

### With Claude Code
```
Use Cadence to analyze if my repository has AI-generated commits.
Analyze https://github.com/user/repo
```

### With GitHub Copilot
Cadence can be used to validate code quality in pull request analysis.

### With Custom Agents
```go
// Agents can invoke Cadence skills programmatically
skill := cadence.GetSkill("detect-suspicious-commit")
result := skill.Call(context.Background(), map[string]interface{}{
    "commit_hash": "abc1234",
    "additions": codeAdditions,
    "use_ai": true,
})
```

## Configuration

Cadence skills respect the standard configuration:

```yaml
thresholds:
  suspicious_additions: 500
  suspicious_deletions: 1000
  max_additions_per_min: 100

ai:
  enabled: true
  provider: openai
  api_key: ${CADENCE_AI_KEY}
  model: gpt-4-mini
```

## Environment Variables

- `CADENCE_AI_KEY` - OpenAI API key for AI analysis
- `CADENCE_CONFIG` - Path to configuration file

## Skill Features

✅ **Token Efficient** - Only analyzes flagged commits
✅ **AI Powered** - Optional GPT-4 Mini analysis
✅ **Async Ready** - Works with agent job queues
✅ **Composable** - Chain with other skills
✅ **Well Documented** - Clear input/output schemas

## Contributing

To improve Cadence skills:
1. Fork the repository
2. Add new detection strategies
3. Improve prompting logic
4. Submit a pull request

See [skills.sh](https://skills.sh/) for more information about the ecosystem.
