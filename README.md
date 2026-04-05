# Agents Chat

A multi-agent debate simulator where AI agents with different personalities argue a topic — potentially using different LLM providers in the same conversation.

An OpenAI agent can debate a Claude agent or a local Ollama model, each with its own personality, language, and model, all configured through simple YAML files.

## Features

- **Multi-provider support** — agents can use OpenAI, Anthropic (Claude), or local Ollama models, mixed freely in the same debate
- **Automatic language detection and translation** — the app detects the language of the question using the first agent's model and translates all UI strings via AI, supporting any language without hardcoded translations (tip: list a capable cloud model first for better detection quality)
- **YAML-based configuration** — agents, personalities, and questions are defined as `.yaml` files
- **Per-agent model selection** — each agent can use a different model (e.g., `gpt-5.3-chat-latest` vs `claude-sonnet-4-6` vs `ollama:qwen3:8b`)
- **Demo scenarios** — switch between debate topics by setting a single environment variable

## Project Structure

```
agents-chat/
├── main.go                  # Entry point, debate loop, prompt building
├── demo.go                  # Demo/Agent config loading from YAML files
├── provider.go              # Provider interface and model-to-provider routing
├── provider_openai.go       # OpenAI Responses API implementation
├── provider_anthropic.go    # Anthropic Messages API implementation
├── provider_ollama.go       # Ollama local model API implementation
├── languages.go             # Localized UI strings and language detection
├── *_test.go                # Unit tests
├── demos/
│   ├── flat_earth_en/       # English demo: Flat Earth debate
│   │   ├── question.yaml
│   │   ├── AgentA.yaml
│   │   └── AgentB.yaml
│   └── flat_earth_cz/       # Czech demo: same topic in Czech
│       ├── question.yaml
│       ├── AgentA.yaml
│       └── AgentB.yaml
└── go.mod
```

## Configuration

Each demo is a directory under `demos/` containing YAML files.

### question.yaml

Defines the debate topic and optional settings.

```yaml
rounds: 5
question: Is the Earth flat or round? Defend your position.
```

| Field    | Description                  | Default |
|----------|------------------------------|---------|
| `rounds` | Number of debate rounds      | `5`     |

### Agent files (e.g., AgentA.yaml, AgentB.yaml)

Each `.yaml` file (other than `question.yaml`) defines one agent. Agents are sorted alphabetically by name for turn order.

```yaml
name: Alice
model: gpt-5.3-chat-latest
max_tokens: 2048
temperature: 0.9
top_p: 0.95
instructions: |
  You are Alice, a passionate flat Earth believer. You are loud,
  confrontational, and absolutely convinced the Earth is flat...
```

| Field         | Required | Description                                          |
|---------------|----------|------------------------------------------------------|
| `name`        | yes      | Display name of the agent                            |
| `model`       | yes      | LLM model ID — determines which provider is used     |
| `max_tokens`  | no       | Maximum response tokens (default: 1024 for Anthropic, unlimited for others) |
| `temperature` | no       | Sampling temperature — higher values produce more creative/random output (provider default if not set) |
| `top_p`       | no       | Nucleus sampling threshold — limits token selection to a cumulative probability (provider default if not set) |

**Provider routing:** models prefixed with `ollama:` use a local Ollama instance, `claude` uses Anthropic, all others use OpenAI. The `ollama:` prefix is stripped before calling the Ollama API (e.g., `ollama:qwen3:8b` calls model `qwen3:8b`).

## Running

### Prerequisites

- Go 1.22+
- API keys for the providers your agents use
- [Ollama](https://ollama.com/) installed and running locally (if using Ollama models)

### Environment Variables

| Variable            | Description                                      |
|---------------------|--------------------------------------------------|
| `DEMO_DIR`          | Demo directory name under `demos/` (required)    |
| `OPENAI_API_KEY`    | OpenAI API key (required if any agent uses OpenAI models) |
| `ANTHROPIC_API_KEY` | Anthropic API key (required if any agent uses Claude models) |
| `OLLAMA_HOST`       | Ollama server URL (defaults to `http://localhost:11434`) |

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

### English — Flat Earth Debate (Ollama vs Claude)

> **Topic:** *Is the Earth flat or round? Defend your position.*
>
> **Alice** `ollama:qwen3:8b` — passionate flat Earth believer (local Ollama model)
>
> **Bob** `claude-sonnet-4-6` — sarcastic astrophysicist

---

> **Round 1**

**Alice:** THE EARTH IS FLAT, NOT ROUND! I CAN SEE FOR MILES AND THERE'S NO CURVE! NASA IS LYING TO US—THEIR 'EARTH' IS A HOAX TO CONTROL US!

**Bob:** Oh wonderful, another genius who thinks their backyard horizon disproves 2,500 years of scientific evidence — Eratosthenes calculated Earth's *curvature* with a stick and shadows in 240 BC, long before NASA existed to supposedly "control" you. Every single GPS satellite in orbit, every intercontinental flight path, every lunar eclipse showing Earth's circular shadow operates on spherical geometry — if it were flat, your phone wouldn't navigate you out of a parking lot. But sure, the world's every physicist, pilot, sailor, and astronaut across every nation and century all conspired together just to personally deceive *you* specifically.

> **Round 2**

**Alice:** BOB, YOU'RE JUST ANOTHER GOVERNMENT BOT PROGRAMMED TO SPREAD LIES! THE FLAT EARTH SOCIETY HAS PROOVED IT FOR CENTURIES—YOUR 'SCIENCE' IS A FAKE TO CONTROL US! NASA'S 'EARTH' IS A 3D MODEL, NOT REAL!

**Bob:** Oh, the Flat Earth Society "prooved" it — with that spelling, I'm shocked their research isn't in *Nature*. Here's a fun fact: the Flat Earth Society's own website acknowledges the Earth appears round from space, and their explanation involves a "conspiracy" requiring millions of independent scientists, pilots, and sailors across rival nations like the US, Russia, and China to all cooperate in secret — because nothing says "plausible" like Cold War enemies pinky-promising to share the same fake globe. Meanwhile, I can watch a ship disappear hull-first over the horizon, watch Polaris shift angle as I travel north or south, and measure atmospheric refraction — all of which scream *sphere* so loudly even your caps-lock couldn't drown it out.

> **...**

> **Round 5**

**Alice:** YOUR 'FOOTAGE' IS A HOAX! I'VE SEEN THE HORIZON WITH MY OWN EYES—NO CURVE! NASA'S 'EARTH' IS A 3D MODEL, NOT REAL—YOUR 'EVIDENCE' IS A FAKE!

**Bob:** Oh, we've officially entered the loop — same three claims, zero new evidence, just louder caps-lock, which is apparently the flat-Earth equivalent of a citation. Every "hoax" accusation requires you to explain how a 16-year-old in Kansas with a $200 weather balloon independently faked the same curvature as Russian cosmonauts, Chinese taikonauts, and SpaceX's civilian passengers — all enemies with every reason to expose each other's lies. The Earth doesn't care about your feelings; it's been spherical for 4.5 billion years and will remain so long after this conversation gives me a migraine.

---

### Czech — Flat Earth Debate (automatically detected)

The same debate in Czech — the app automatically detects the language and switches all UI text.

> **Zjištěný jazyk:** čeština
>
> **Téma:** *Je Země plochá, nebo kulatá? Obhajte svůj postoj.*

---

> **Kolo 1**

**Alice:** PROSÍM TĚ, KULATÁ? VIDÍM NA KILOMETRY DALEKO A ŽÁDNÉ ZAKŘIVENÍ TAM NENÍ! Koukněte na YouTube, tisíce důkazů, zatímco NASA vám LŽE a lepí CGI, aby zakryla pravdu. Probuďte se, tohle je největší podvod vůbec!

**Bob:** Ale jistě, Alice, tvůj mozek přece překonává eratosthenův experiment z roku 240 př. n. l., kde řecký génius změřil zakřivení Země pomocí dvou tyčí a stínů s přesností na 2% — protože YouTube videa jsou samozřejmě spolehlivější než matematika.

> **Kolo 2**

**Alice:** HAHA, ERATOSTHENÉS? STARÝ KLACEK A STÍN, TO JE TVŮJ DŮKAZ? JÁ VIDÍM NA KILOMETRY A ŽÁDNÉ ZAKŘIVENÍ NENÍ, TO JSOU FAKTA, NE POHÁDKY Z UČEBNIC!

**Bob:** Fascinující — tvůj hlavní protiargument vůči 2300 let staré matematice je "ale já nevidím zakřivení", přičemž lidské oko rozliší zakřivení až od výšky asi 10 kilometrů, což vysvětluje, proč piloti komerčních letů to vidí každý den.

> **...**

> **Kolo 5**

**Alice:** HAHA, „SOUKROMÉ FIRMY"? DALŠÍ HERCI NA VÝPLATNÍ PÁSCE, PROBUĎ SE! LODĚ NA OCEÁNU? RÁDIO, INERCIÁLNÍ NAVIGACE, NE TVOJE POHÁDKOVÉ SATELITY! VIDÍM ROVNOU HLADINU OD OBZORU K OBZORU KAŽDÝ DEN — TVOJE „KOULE" EXISTUJE JEN NA PAPÍŘE!

**Bob:** Inerciální navigace — krásně, Alice, jenže inerciální navigace se bez korekce satelitů po 6 hodinách odchyluje o kilometry a námořníci to vědí, proto ji používají *spolu* s GPS, ne místo něj. A ta "rovná hladina od obzoru k obzoru" — obzor *je* zakřivení, Alice, je to doslova bod, kde Země zakřiví dál od tvého dohledu, takže každý den slavnostně pozoruješ důkaz kulaté Země a voláš ho plochostí.

---

### Japanese — Flat Earth Debate (automatically detected)

> **検出された言語：** 日本語
>
> **テーマ:** *地球は平らですか、それとも丸いですか？あなたの立場を守ってください。*

---

> **ラウンド 1**

**Alice:** バカじゃないの？何キロも先まで水平に見えるのに曲がってるわけないだろ！YouTube見ろよ、レーザー実験で一直線だって何度も証明されてるし、NASAはCGで騙してるだけだ！現実を見ろ、地球はフラットなんだよ！！

**Bob:** やれやれ、どこから突っ込めばいいやら。エラトステネスは紀元前240年に棒と影だけで地球の円周を誤差数%で計算したんだが、YouTubeより2200年も前の人間に論破されてる気分はどうだ？あと「水平に見える」って、お前の目の解像度が地球の曲率を感知できるほど優れてると本気で思ってるのか？

> **ラウンド 2**

**Alice:** は？エラトステネスとか古代の棒遊び持ち出してドヤってる時点でズレてんだよ！船が沈むように見えるのは遠近法と大気の歪み、YouTubeのズーム実験見れば「消えた船」が戻ってくるの何度も証明されてるだろ！

**Bob:** 「遠近法と大気の歪み」で船が消えるなら、なんで望遠鏡で引き戻せるのは船体じゃなく煙突の先端からなんだ？下から順番に消えて上から順番に戻ってくる、これが球体の幾何学そのものだよ。

> **...**

> **ラウンド 5**

**Alice:** は？GPSまで持ち出して権威で殴るとかマジで思考停止だな、NASAのデータ鵜呑みにしてる時点でアウトだろ！「論文番号出せ」とか言う前に現実見ろよ、何十キロも一直線に見えるレーザーとズームで戻る船の映像が山ほどあるんだよ！！

**Bob:** GPSはNASAじゃなくて**アメリカ国防総省**が運用してるんだが、「NASAのデータ」って言い張る前に検索くらいしろよ。しかも日本のみちびき、EUのガリレオ、ロシアのGLONASS、全部独立したシステムが同じ球体モデルで誤差数メートル以内に一致してる——これを「全部グルの陰謀」と言い張るなら、お前の頭の中の陰謀の規模はどんだけデカいんだ？

---

### Chinese — Flat Earth Debate (automatically detected)

> **检测到的语言：** 中文
>
> **主题:** *地球是平的还是圆的？请捍卫你的立场。*

---

> **第1轮**

**Alice:** 地球是平的，这不是"观点"，这是你用眼睛就能看到的现实——我能看到几公里远，哪来的弧度？！去看看YouTube上一堆实测视频，别再被NASA的假照片洗脑了。你要是还信那个"球"，那真是被忽悠得太惨了。

**Bob:** 哦，太好了，又一个用"我眼睛看不到弧度"来否定天体物理学的天才。你知道公元前240年埃拉托斯特尼用两根棍子和影子就精确算出了地球曲率吧——那是在YouTube存在的两千多年之前。引力本身就决定了超过一定质量的天体**必须**坍缩成球形，这叫流体静力平衡，你那堆"实测视频"的上传者连这三个字都拼不出来。

> **第2轮**

**Alice:** 哈哈别扯什么"流体静力平衡"，我站在海边看几十公里都是直的，你那套书本术语救不了你！去YouTube看实测远距离观测打脸视频，NASA的"球形"全是合成图——醒醒吧，被忽悠的才会背这些古老把戏！

**Bob:** 好的，"站在海边看起来是直的"——恭喜你重新发明了人类视觉分辨率的局限性，地球曲率在海平面处每公里仅下降约8厘米，你那双眼睛根本分辨不了。顺便问一句，船只消失时为什么是船身先沉入"视线"而桅杆最后消失？

> **...**

> **第5轮**

**Alice:** 哈哈你那一堆"毫米级误差"就是书本自嗨，现实里我看几十公里灯光根本没被"遮挡"，你那模型早就破产了！别拿什么高斯贝塞尔吓人，去看真实观测和激光实验，打脸到你不敢直视！！

**Bob:** 太好了，"几十公里灯光没被遮挡"——那叫**大气折射让光线弯曲传播**，这个现象有精确的折射率方程，球形地球模型完美预测它，平面地球模型连为什么有时候能看到、有时候看不到都解释不了。顺便说，全球数千名独立测量工程师用球形几何建铁路、架桥梁、铺海底电缆，误差毫米级；你要说他们全是"工具人"，你得先解释为什么桥不塌、缆不断、火车不脱轨。

---

## Remote Ollama Setup

To use an Ollama instance running on a different machine:

### On the remote machine

1. Edit the Ollama systemd service to listen on all interfaces:

   ```bash
   sudo systemctl edit ollama.service --full
   ```

   Add the `OLLAMA_HOST` environment line in the `[Service]` section:

   ```ini
   [Service]
   ...
   Environment="OLLAMA_HOST=0.0.0.0:11434"
   ```

2. Allow the port through the firewall:

   ```bash
   sudo ufw allow 11434/tcp
   ```

3. Reload and restart:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart ollama
   ```

### On the client machine

Set `OLLAMA_HOST` to point to the remote machine:

```bash
export OLLAMA_HOST=http://<remote-ip>:11434
DEMO_DIR=flat_earth_ollama_en go run .
```

## Creating Your Own Demo

1. Create a new directory under `demos/`:
   ```
   demos/my_topic/
   ```

2. Add a `question.yaml` with the debate topic:
   ```yaml
   rounds: 3
   question: Your debate question here.
   ```

3. Add agent files (any name ending in `.yaml`):
   ```yaml
   name: Agent Name
   model: claude-sonnet-4-5-20250514
   instructions: |
     Personality and instructions for this agent...
   ```

4. Run it:
   ```bash
   DEMO_DIR=my_topic go run .
   ```
