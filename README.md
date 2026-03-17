# Agents Chat

A multi-agent debate simulator where AI agents with different personalities argue a topic — potentially using different LLM providers in the same conversation.

An OpenAI agent can debate a Claude agent, each with its own personality, language, and model, all configured through simple markdown files.

## Features

- **Multi-provider support** — agents can use OpenAI or Anthropic (Claude) models, mixed freely in the same debate
- **Automatic language detection** — the app detects the language of the question using an LLM call and automatically selects the matching UI language (Czech, English, German, French, Spanish, Portuguese, Italian), falling back to English
- **Markdown-based configuration** — agents, personalities, and questions are defined as `.md` files with YAML-like frontmatter
- **Per-agent model selection** — each agent can use a different model (e.g., `gpt-5-mini` vs `claude-haiku-4-5`)
- **Demo scenarios** — switch between debate topics by setting a single environment variable

## Project Structure

```
agents-chat/
├── main.go                  # Entry point, debate loop, prompt building
├── demo.go                  # Demo/Agent config loading from markdown files
├── provider.go              # Provider interface and model-to-provider routing
├── provider_openai.go       # OpenAI Responses API implementation
├── provider_anthropic.go    # Anthropic Messages API implementation
├── languages.go             # Localized UI strings and language detection
├── *_test.go                # Unit tests
├── demos/
│   ├── flat_earth/          # English demo: Flat Earth debate
│   │   ├── Question.md
│   │   ├── AgentA.md
│   │   └── AgentB.md
│   └── flat_earth_cz/       # Czech demo: same topic in Czech
│       ├── Question.md
│       ├── AgentA.md
│       └── AgentB.md
└── go.mod
```

## Configuration

Each demo is a directory under `demos/` containing markdown files.

### Question.md

Defines the debate topic and optional settings. The frontmatter is optional.

```markdown
---
rounds: 5
---
Is the Earth flat or round? Defend your position.
```

| Field    | Description                  | Default |
|----------|------------------------------|---------|
| `rounds` | Number of debate rounds      | `5`     |

### Agent files (e.g., AgentA.md, AgentB.md)

Each `.md` file (other than `Question.md`) defines one agent. Agents are sorted alphabetically by name for turn order.

```markdown
---
name: Alice
model: gpt-5-mini
---
You are Alice, a passionate flat Earth believer. You are loud,
confrontational, and absolutely convinced the Earth is flat...
```

| Field   | Required | Description                                          |
|---------|----------|------------------------------------------------------|
| `name`  | yes      | Display name of the agent                            |
| `model` | yes      | LLM model ID — determines which provider is used     |

**Provider routing:** models starting with `claude` use Anthropic, all others use OpenAI.

## Running

### Prerequisites

- Go 1.22+
- API keys for the providers your agents use

### Environment Variables

| Variable            | Description                                      |
|---------------------|--------------------------------------------------|
| `DEMO_DIR`          | Demo directory name under `demos/` (required)    |
| `OPENAI_API_KEY`    | OpenAI API key (required if any agent uses OpenAI models) |
| `ANTHROPIC_API_KEY` | Anthropic API key (required if any agent uses Claude models) |

### Run

```bash
# Run a demo
DEMO_DIR=flat_earth go run .

# Or specify a full path as CLI argument
go run . demos/flat_earth_cz
```

### Tests

```bash
go test -v ./...
```

## Sample Output

### English — Flat Earth Debate (OpenAI vs Claude)

> **Topic:** *Is the Earth flat or round? Defend your position.*
>
> **Alice** `gpt-5-mini` — passionate flat Earth believer
> 
> **Bob** `claude-haiku-4-5` — sarcastic astrophysicist

---

> **Round 1**

**Alice:** GIVE ME A BREAK — I can stand on the pier, see for MILES, and THERE IS NO CURVE, so stop swallowing NASA'S CGI fairy tales! Watch the YouTube investigations proving the flat horizon and wake up, you round-Earth cultists — you're embarrassing yourselves.

**Bob:** Look, the horizon *looks* flat because you're standing on a sphere the size of 40 million football fields — your eyeball is cosmically insignificant. Ships disappearing hull-first, time zones, satellite GPS that literally guides your phone, and the fact that every other planet we've observed is spherical should've clued you in by now, but here we are.

> **Round 2**

**Alice:** OH PLEASE — your little textbook spiel doesn't change the FACT that I can stand on a pier and SEE FOR MILES with NO CURVE, so stop drinking NASA'S CGI Kool-Aid!

**Bob:** Alice, you can "see for miles" on a 40,000-km circumference — that's 0.006% of Earth's surface, so yeah, your perspective is geometrically useless. Go to any port, watch ships sail away, and notice the hull disappears *before* the mast *every single time*.

> **...**

> **Round 5**

**Alice:** NASA FAKEUJE FOTKY A KDO TO NEVIDÍ, TEN JE BUĎ BLÁZEN, NEBO KOLABORANT!

**Bob:** Alice, you're trapped in a loop — literally the same argument repeated five times while ignoring every testable prediction you could verify yourself. Take a $200 telescope, go to any port at sunset, and watch a ship disappear hull-first; that's not NASA, that's *you* observing basic perspective on a curve.

---

### Czech — Flat Earth Debate (automatically detected)

The same debate in Czech — the app automatically detects the language and switches all UI text.

> **Detekovaný jazyk:** Čeština
> 
> **Téma:** *Je Země plochá, nebo kulatá? Obhajte svůj postoj.*

---

> **Kolo 1**

**Alice:** Země JE PLOCHÁ — podívejte se na stovky videí na YouTube a NA VLASTNÍ OČI VIDÍM NA KILOMETRY DALEKO A ŽÁDNÉ ZAKŘIVENÍ TAM NENÍ! NASA LŽE A KDO V TOM VĚŘÍ, TEN JE BUĎ NAIVNÍ, NEBO SOUČÁSTÍ KONSPIRACE.

**Bob:** Ahoj Alice, krásné, že věříš vlastním očím — bohužel jsou na to krátké. Satelity vidí Zemi každý den kulatou, a ty nejsou vlastnictví NASA — Rusko, Čína, EU, všichni ji fotografují jako glóbus. Pokud je to konspirace, musela by být největší v dějinách.

> **Kolo 2**

**Alice:** BOHUŽEL, BOB, TY JSI BUĎ NAIVNÍ, NEBO SOUČÁST KONSPIRACE! VŠECHNY „SATELITNÍ" FOTKY JSOU PODVOD NASA A JEJÍCH PARTNERŮ!

**Bob:** Alice, satelity nejsou vlastnictví jednoho státu — Rusko, Čína, Indie, EU a desítky soukromých firem mají své, a všechny vidí to samé, protože realita je objektivní, ne americká konspirace.

> **...**

> **Kolo 5**

**Bob:** YouTube ti nedá nic, co by vrátilo Zemi do plochy, gravitace to prostě nedovoluje — a to není názor, to je měřitelný fakt, který funguje stejně v každém státě na světě.

---

## Creating Your Own Demo

1. Create a new directory under `demos/`:
   ```
   demos/my_topic/
   ```

2. Add a `Question.md` with the debate topic:
   ```markdown
   ---
   rounds: 3
   ---
   Your debate question here.
   ```

3. Add agent files (any name ending in `.md`):
   ```markdown
   ---
   name: Agent Name
   model: claude-sonnet-4-5-20250514
   ---
   Personality and instructions for this agent...
   ```

4. Run it:
   ```bash
   DEMO_DIR=my_topic go run .
   ```
