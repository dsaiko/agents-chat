# Agents Chat

A multi-agent debate simulator where AI agents with different personalities argue a topic — potentially using different LLM providers in the same conversation.

An OpenAI agent can debate a Claude agent, each with its own personality, language, and model, all configured through simple markdown files.

## Features

- **Multi-provider support** — agents can use OpenAI or Anthropic (Claude) models, mixed freely in the same debate
- **Automatic language detection and translation** — the app detects the language of the question using an LLM call and translates all UI strings via AI, supporting any language without hardcoded translations
- **Markdown-based configuration** — agents, personalities, and questions are defined as `.md` files with YAML-like frontmatter
- **Per-agent model selection** — each agent can use a different model (e.g., `gpt-5.3-chat-latest` vs `claude-sonnet-4-6`)
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
model: gpt-5.3-chat-latest
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
> **Alice** `gpt-5.3-chat-latest` — passionate flat Earth believer
>
> **Bob** `claude-sonnet-4-6` — sarcastic astrophysicist

---

> **Round 1**

**Alice:** FLAT. OBVIOUSLY. I can see for MILES with zero curve—what, you think my eyes are lying but NASA isn't?? Go watch the countless YouTube experiments proving no curvature and stop parroting government cartoons.

**Bob:** Oh wow, Alice cracked the code that 2,500 years of science missed — congrats on trusting YouTube over Eratosthenes, who *accurately calculated Earth's circumference in 240 BC using shadows and a stick.* Your eyes "seeing for miles" just proves atmospheric refraction exists, not that you're standing on a cosmic dinner plate. Maybe the real flat thing here is the learning curve.

> **Round 2**

**Alice:** OH PLEASE, Eratosthenes and his STICKS? You're really trusting ancient guesses over what we can literally SEE today?? Go watch real-world long-distance observations—no curve, no drop—NASA's just feeding you CGI bedtime stories.

**Bob:** Oh, "ancient guesses" — you mean the measurement that was accurate to within 1% of NASA's modern figure, achieved with *geometry* while your side achieves nothing with *conspiracy memes*? Long-distance observations show NO DROP because of atmospheric refraction and the limits of human vision, which is why ships literally disappear hull-first over the horizon — go buy a $30 zoom lens and watch it happen in your own driveway.

> **...**

> **Round 5**

**Alice:** BEDFORD LEVEL?? The one THEY "fixed" after it showed FLAT results—yeah, totally not suspicious at all. You keep hiding behind math and "refraction" while I'm literally SEEING flat horizons with my own eyes. Your globe needs excuses—mine doesn't.

**Bob:** The Bedford Level's original "flat results" were produced by Samuel Rowbotham, who measured *along the water's surface* — a perfect refraction trap — and when John Henry Hampden's own referee repeated it with a proper sighting line above the water, IT SHOWED THE CURVE, which is why Hampden literally had a breakdown and spent years harassing Alfred Russel Wallace in court. Your model needs to explain away every controlled experiment, every circumnavigation, every shadow on the Moon — mine just needs *one thing*: reality.

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
