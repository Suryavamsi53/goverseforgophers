# GoVerse for Gophers 🐹🚀

GoVerse is an enterprise-grade, interactive learning platform and native code execution sandbox built exclusively for mastering the Go programming language.

## 🌟 Features

- **Native Go Sandbox**: A high-performance, browser-based IDE powered by Monaco Editor. Write, compile, and execute Go code instantly using a true native OS-level backend execution engine. No AI mocks—just real `go run` execution.
- **Interactive Learning Engine**: A custom Markdown parsing engine that renders beautiful, dynamic Go lessons directly to the browser.
- **Developer-First CLI**: Includes a standalone `goverse-cli` tool to compile and evaluate your code locally via terminal.
- **User Progress Dashboards**: Fully integrated PostgreSQL database to track learning streaks, profile metrics, and course completions.
- **Glassmorphism UI**: A gorgeous, modern frontend built with Alpine.js, HTMX, and Tailwind CSS.

## 🛠️ Tech Stack

- **Backend:** Go (Golang), Chi Router, standard `os/exec` isolated compilation
- **Database:** PostgreSQL 15, `pgx/v5`
- **Frontend:** HTMX, Alpine.js, Tailwind CSS, Monaco Editor
- **Infrastructure:** Docker Compose

## 🚀 Getting Started

### 1. Start the Database
Make sure you have Docker installed, then spin up the PostgreSQL instance:
```bash
cd deployments
docker compose up -d
```

### 2. Run the Platform
Start the Go backend server:
```bash
go run ./cmd/server/main.go
```
The application will be live at `http://localhost:8080`.

### 3. Quick Links
- **Dashboard**: `http://localhost:8080/dashboard`
- **Practice Sandbox**: `http://localhost:8080/practice`
- **Learn Modules**: `http://localhost:8080/learn`

## 💻 CLI Usage
To execute Go files quickly from your terminal using the GoVerse engine:
```bash
go run ./cmd/cli/main.go -file your_file.go
```

## 🗺️ GoVerse Platform Roadmap

> **Interactive Navigation:**
> - 🔭 **Macro View (Zoom Out):** View the high-level end-to-end lifecycle roadmap below.
> - 🔍 **Micro View (Zoom In):** Expand the collapsible sections under each phase to drill down into granular components, execution pipelines, and curriculum tracks.
> - 💡 *On GitLab & GitHub, hover over or click Mermaid diagrams to use built-in pan and zoom controls.*

```mermaid
%%{init: {
  'theme': 'base',
  'themeVariables': {
    'darkMode': true,
    'background': '#0b0f19',
    'primaryColor': '#111827',
    'primaryTextColor': '#f3f4f6',
    'primaryBorderColor': '#00ADD8',
    'lineColor': '#38bdf8',
    'secondaryColor': '#1f2937',
    'tertiaryColor': '#1e293b',
    'mainBkg': '#111827',
    'nodeBorder': '#00ADD8',
    'clusterBkg': '#0b1120',
    'clusterBorder': '#1e293b',
    'fontFamily': 'Inter, system-ui, -apple-system, sans-serif'
  }
}}%%
flowchart TB
    %% Master End-to-End Platform Roadmap Flowchart

    subgraph P1["⚡ PHASE 1: Core Foundation & Execution Engine"]
        direction TB
        P1_A["⚙️ Native OS Execution Engine<br/><code>os/exec</code> + PTY Streamer"]
        P1_B["📝 Custom Markdown Engine<br/>Goldmark + Prism Highlighter"]
        P1_C["💻 GoVerse Developer CLI<br/>Local Compiler & Benchmark Harness"]
        P1_A --> P1_B --> P1_C
    end

    subgraph P2["🌐 PHASE 2: Interactive Web Platform & Real-Time IDE"]
        direction TB
        P2_A["🎨 Modern Glassmorphism UI<br/>HTMX + Alpine.js + Tailwind CSS"]
        P2_B["🖥️ Monaco Browser IDE<br/>IntelliSense & Live Compilation"]
        P2_C["🗄️ PostgreSQL Persistence<br/><code>pgx/v5</code> + Migration Pipelines"]
        P2_D["🔐 Auth & Session Engine<br/>Argon2id + JWT + HTTP-only Cookies"]
        P2_A --> P2_B --> P2_C --> P2_D
    end

    subgraph P3["📚 PHASE 3: 27-Track Master Go Curriculum"]
        direction TB
        P3_A["🌱 Fundamentals & Memory<br/>Pointers, GC, Allocations, Escape Analysis"]
        P3_B["⚡ Concurrency & Goroutines<br/>Channels, Mutexes, Select, CSP Patterns"]
        P3_C["🏛️ Clean Architecture & APIs<br/>Domain-Driven Design, REST, gRPC, Protobuf"]
        P3_D["🏗️ Distributed Systems & K8s<br/>Kafka, Redis, Docker, Custom Operators"]
        P3_A --> P3_B --> P3_C --> P3_D
    end

    subgraph P4["☸️ PHASE 4: Cloud-Native & Production Isolation"]
        direction TB
        P4_A["📦 Multi-Stage Containerization<br/>Alpine Minimal Distroless Builds"]
        P4_B["🛡️ Sandboxed Isolation<br/>gVisor / Firecracker MicroVM Sandboxing"]
        P4_C["☸️ Kubernetes Orchestration<br/>Helm Charts, StatefulSets, Auto-Scaling"]
        P4_D["📊 Real-Time Observability<br/>OpenTelemetry, Prometheus & Grafana"]
        P4_A --> P4_B --> P4_C --> P4_D
    end

    subgraph P5["🚀 PHASE 5: AI Mentorship & Global Community"]
        direction TB
        P5_A["🧠 AI Code Reviewer<br/>AST-based hints, Idiomatic Go suggestions"]
        P5_B["🏆 Competitive Arena<br/>Live Coding Battles & Streak Leaderboards"]
        P5_C["🤝 Community Ecosystem<br/>User-created Lessons & Gopher Plugins"]
        P5_A --> P5_B --> P5_C
    end

    %% Phase-to-Phase Progression Connections
    P1 ==>|Engine Ready| P2
    P2 ==>|Platform Live| P3
    P3 ==>|Curriculum Scale| P4
    P4 ==>|Enterprise Scale| P5

    %% Node Styling Classes
    classDef completed fill:#064e3b,stroke:#10b981,stroke-width:2px,color:#ecfdf5;
    classDef inprogress fill:#0c4a6e,stroke:#00ADD8,stroke-width:2px,color:#f0f9ff;
    classDef planned fill:#3b0764,stroke:#a855f7,stroke-width:2px,color:#faf5ff;
    classDef future fill:#451a03,stroke:#f59e0b,stroke-width:2px,color:#fffbeb;

    class P1_A,P1_B,P1_C,P2_A,P2_B,P2_C completed;
    class P2_D,P3_A,P3_B,P3_C,P4_A inprogress;
    class P3_D,P4_B,P4_C,P4_D planned;
    class P5_A,P5_B,P5_C future;
```

---

### 🎨 Roadmap Status Legend

| Status Icon | Status | Phase Coverage | Description |
|:---:|:---|:---|:---|
| 🟢 | **Completed** | Phase 1 & Core Phase 2 | Production-tested and operating in current releases |
| 🔵 | **In Active Progress** | Phase 2 & Phase 3 | Currently being developed and expanded |
| 🟣 | **Planned (Next Up)** | Phase 3 & Phase 4 | Scheduled for upcoming engineering sprints |
| 🟡 | **Future Vision** | Phase 5 | Long-term strategic evolution and AI roadmap |

---

### 🔍 Interactive Drill-Down (Zoom In / Zoom Out)

<details>
<summary><b>🔍 Zoom In: Phase 1 — Native Execution & Sandbox Architecture</b></summary>

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'darkMode': true, 'primaryColor': '#111827', 'primaryTextColor': '#fff', 'lineColor': '#10b981'}}}%%
flowchart LR
    UserCode["👨‍💻 Client Code Submission"] --> API["🌐 /api/v1/run Handler"]
    API --> Validator["🛡️ AST Safety Validator"]
    Validator --> TempFs["📁 Ephemeral RAM Disk (/tmp)"]
    TempFs --> Engine["⚡ Runner Engine (os/exec)"]
    Engine --> PTY["🖥️ Pseudo-Terminal (PTY Stream)"]
    Engine --> Timeout["⏱️ Context Timeout (10s Hard Cap)"]
    Engine --> MemoryCap["💾 Cgroup Memory Limiter (128MB)"]
    PTY --> OutputSanitizer["🧹 Output Stream Sanitizer"]
    OutputSanitizer --> ClientTerminal["📺 Monaco Interactive Console"]

    classDef pass fill:#064e3b,stroke:#10b981,stroke-width:2px,color:#ecfdf5;
    class UserCode,API,Validator,TempFs,Engine,PTY,Timeout,MemoryCap,OutputSanitizer,ClientTerminal pass;
```
</details>

<details>
<summary><b>🔍 Zoom In: Phase 2 — Web Platform, Monaco IDE & HTMX Reactive Flow</b></summary>

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'darkMode': true, 'primaryColor': '#111827', 'primaryTextColor': '#fff', 'lineColor': '#00ADD8'}}}%%
flowchart LR
    Browser["🌐 Browser (HTMX + Alpine.js)"] -->|HTTP / SSE| Chi["⚡ Chi Router Middleware"]
    Chi --> AuthMiddleware["🔐 Session / JWT Guard"]
    AuthMiddleware --> Delivery["📦 Web Delivery Controllers"]
    
    subgraph Services["Core Application Layer"]
        LearnSvc["📖 Learn & Progress Service"]
        PracticeSvc["🧪 Interactive Sandbox Service"]
        AuthSvc["👤 User Account & Profile Service"]
    end
    
    Delivery --> LearnSvc & PracticeSvc & AuthSvc
    LearnSvc & PracticeSvc & AuthSvc --> Repo["🗄️ PostgreSQL Repository (pgxpool)"]
    PracticeSvc --> Runner["🚀 Execution Engine (pkg/runner)"]

    classDef web fill:#0c4a6e,stroke:#00ADD8,stroke-width:2px,color:#f0f9ff;
    class Browser,Chi,AuthMiddleware,Delivery,LearnSvc,PracticeSvc,AuthSvc,Repo,Runner web;
```
</details>

<details>
<summary><b>🔍 Zoom In: Phase 3 — 27-Track Learning Curriculum Matrix</b></summary>

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'darkMode': true, 'primaryColor': '#111827', 'primaryTextColor': '#fff', 'lineColor': '#a855f7'}}}%%
mindmap
  root((🐹 GoVerse Curriculum))
    Core Foundations
      ::icon(fa fa-code)
      Go Fundamentals
      Memory Management & GC
      Clean Architecture
      Design Patterns
      Testing & Benchmarking
    High Concurrency & Systems
      ::icon(fa fa-bolt)
      Goroutines & Channels
      Sync & Atomic Primitives
      Networking & Sockets
      HTTP & REST Engineering
      gRPC & Protobuf
    Distributed Infrastructure
      ::icon(fa fa-server)
      PostgreSQL & Relational Design
      Redis Caching & PubSub
      Kafka Event Streaming
      Distributed Systems & Consensus
      System Design Case Studies
    Cloud-Native & Production
      ::icon(fa fa-cloud)
      Docker Containerization
      Kubernetes & Helm
      Custom Go K8s Operators
      Observability & OpenTelemetry
      Performance Profiling & pprof
      Security & Hardening
      Google Interview Prep
```
</details>

<details>
<summary><b>🔍 Zoom In: Phase 4 — Cloud-Native, Kubernetes & Production Scaling</b></summary>

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'darkMode': true, 'primaryColor': '#111827', 'primaryTextColor': '#fff', 'lineColor': '#38bdf8'}}}%%
flowchart TD
    Ingress["🌐 Cloud Ingress / TLS Termination"] --> Gateway["🚪 NGINX / Envoy Reverse Proxy"]
    
    subgraph K8s["☸️ Kubernetes Cluster (goverse-prod)"]
        direction TB
        AppPods["🐹 GoVerse Web Pods (Replicas: 3+)"]
        HPA["📈 Horizontal Pod Autoscaler (CPU/Mem metrics)"]
        AppPods --- HPA
        
        RunnerPool["📦 Isolated Sandbox Runner Pool<br/>(gVisor runsc runtime)"]
        AppPods -->|gRPC Internal| RunnerPool
    end
    
    Gateway --> AppPods
    AppPods --> Postgres["🐘 PostgreSQL StatefulSet (PersistentVolume)"]
    AppPods --> Redis["⚡ Redis Cluster (Session & Cache Store)"]
    
    subgraph Observability["📊 Telemetry & Health Monitoring"]
        Prometheus["🔥 Prometheus Metrics"]
        Grafana["📈 Grafana Dashboards"]
        Loki["📜 Loki Log Aggregation"]
    end
    
    AppPods -.-> Prometheus & Loki
    Prometheus --> Grafana

    classDef k8s fill:#1e1b4b,stroke:#818cf8,stroke-width:2px,color:#e0e7ff;
    class Ingress,Gateway,AppPods,HPA,RunnerPool,Postgres,Redis,Prometheus,Grafana,Loki k8s;
```
</details>

<details>
<summary><b>🔍 Zoom In: Phase 5 — AI Mentorship, Gamification & Global Community</b></summary>

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'darkMode': true, 'primaryColor': '#111827', 'primaryTextColor': '#fff', 'lineColor': '#f59e0b'}}}%%
flowchart LR
    Code["💻 User Solution"] --> AstEngine["🌳 AST Analyzer"]
    AstEngine --> AiMentor["🤖 GoVerse AI Mentor"]
    AiMentor --> Feedback["💡 Idiomatic Go Explanations<br/>& Efficiency Recommendations"]
    
    UserAction["🎯 Lesson Completion"] --> Gamification["🏆 Streak & XP Engine"]
    Gamification --> Leaderboard["🌍 Global Gopher Leaderboards"]
    Gamification --> Badges["🎖️ Verifiable Skill Badges"]
    
    Community["👥 Gopher Community"] --> Contribution["📝 Community Track Builder"]
    Contribution --> Review["🔍 Peer Review & QA"]
    Review --> PublishedTrack["🚀 Official Community Tracks"]

    classDef ai fill:#451a03,stroke:#f59e0b,stroke-width:2px,color:#fffbeb;
    class Code,AstEngine,AiMentor,Feedback,UserAction,Gamification,Leaderboard,Badges,Community,Contribution,Review,PublishedTrack ai;
```
</details>

