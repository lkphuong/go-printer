# API Job Print - Flow Diagram

## Overview

API Job Print (`POST /api/v1/printers/jobs`) là endpoint để gửi yêu cầu in một file với các thông số cụ thể.

---

## 🎯 Complete End-to-End Flow (Chi tiết nhất)

---

## 📊 Queue & Worker Processing Flow

```mermaid
graph TD
    subgraph "📥 Request Phase"
        A["🖥️ Client POST /api/v1/printers/jobs<br/>file + type + copies"]
        B["📥 Handler JobPrint<br/>Parse multipart form"]
        C["Validate inputs"]
        D["Call Service.JobPrint"]
    end

    subgraph "⚙️ Service Phase"
        E["📖 Load config.json"]
        F["🔍 Filter printers by type"]
        G{Printers<br/>found?}
    end

    subgraph "💾 File Phase"
        H["📝 Generate temp path<br/>uploads/YYYYMMDDhhmmss_file"]
        I["💾 Save uploaded file"]
        J["📤 Send to Telegram"]
        K["📂 Open file handle"]
    end

    subgraph "📊 Queue Phase"
        L["🔄 Loop: For each printer"]
        M["📨 Create PrintRequest{<br/>Printer, FilePath,<br/>Copies, StartTime,<br/>RetryCount}"]
        N["📤 Send to Queue Channel<br/>printQueue <- req"]
        O{Queue<br/>full?}
    end

    subgraph "🔌 Ping Phase"
        P["🔄 Loop: For each printer"]
        Q["🔌 TCP Connect<br/>ConnectPrinter"]
        R["❌ Close connection<br/>Ping signal sent"]
    end

    subgraph "✅ Response Phase"
        S["📝 Log: print done"]
        T["📤 Response 200 OK"]
        U["🔄 Handler returns"]
    end

    subgraph "👷 Background: Print Worker"
        V["👷 Goroutine listening<br/>for req := range printQueue"]
        W["🎯 Receive PrintRequest"]
        X["🔄 Inner loop: Retry logic"]
        Y["🔌 ConnectPrinter with<br/>TCP DialTimeout 5s"]
        Z["📊 Check printer status<br/>printStatus - bit flags"]
        AA["📂 Open & decode image"]
        AB["🖼️ Image processing:<br/>- Resize if needed<br/>- Convert to 1-bit B&W"]
        AC["🖨️ Build print command<br/>ESC/POS format"]
        AD["🔄 For each copy:<br/>- Write print cmd<br/>- Write blank lines<br/>- Write cut cmd<br/>- Sleep 1s"]
        AE["❌ Cleanup temp file"]
        AF["✅ Return nil<br/>r.Result <- err"]
    end

    subgraph "⏱️ Retry Mechanism"
        AG["❌ Error occurred<br/>in printFile?"]
        AH["RetryCount++"]
        AI["Check time > 1hr?"]
        AJ["🗑️ Delete temp file<br/>if timeout"]
        AK["😴 Sleep 30s"]
        AL["🔄 Retry printFile"]
    end

    A --> B --> C --> D --> E --> F --> G

    G -->|❌ Not Found| T
    G -->|✅ Found| H --> I --> J --> K

    K --> L --> M --> N --> O
    O -->|❌ Full| T
    O -->|✅ Sent| P

    P --> Q --> R --> S --> T --> U

    D -.->|Async| V -.-> W --> X

    X --> Y --> Z
    Z -->|❌ Error| AG
    Z -->|✅ OK| AA --> AB --> AC --> AD --> AE --> AF

    AG --> AH --> AI
    AI -->|❌ Not expired| AK --> AL
    AL --> X
    AI -->|✅ Timeout| AJ --> AF

    style A fill:#e3f2fd
    style T fill:#c8e6c9
    style V fill:#fff9c4
    style W fill:#fff9c4
    style X fill:#fff9c4
    style AG fill:#ffcccc
    style AJ fill:#ffcccc
```

---

## 🔄 Queue System Deep Dive

```mermaid
graph LR
    subgraph CLIENT["🖥️ Client Request"]
        R1["Request 1"]
        R2["Request 2"]
        R3["Request N"]
    end

    subgraph HANDLER["📥 Handler Layer"]
        H1["Parse & Validate"]
        H2["Call Service"]
    end

    subgraph SERVICE["⚙️ Service Layer"]
        S1["Load Config"]
        S2["Filter Printers"]
        S3["Save File"]
        S4["Loop: PrintFileQueued<br/>for each printer"]
    end

    subgraph QUEUE["📊 QUEUE CHANNEL<br/>(Buffer Size: 1000)"]
        Q1["PrintRequest 1"]
        Q2["PrintRequest 2"]
        Q3["PrintRequest 3"]
        STATS["Current: 3<br/>Max: 1000"]
    end

    subgraph WORKER["👷 Print Worker<br/>(Single Goroutine)"]
        W1["Waiting..."]
        W2["Processing<br/>Request 1"]
        W3["Processing<br/>Request 2"]
    end

    subgraph PRINTER["🖨️ Printer Operations"]
        P1["TCP Connect"]
        P2["Check Status"]
        P3["Process Image"]
        P4["Send Print Cmd"]
        P5["Wait 1s between copies"]
        P6["Send Cut Command"]
        P7["Close Connection"]
    end

    subgraph RESULT["✅ Results"]
        RES1["r.Result <- nil"]
        RES2["r.Result <- error"]
    end

    CLIENT --> HANDLER
    HANDLER --> SERVICE
    SERVICE --> S1 --> S2 --> S3 --> S4

    S4 -->|Add to Queue| QUEUE
    QUEUE -.->|Non-blocking<br/>send| QUEUE
    QUEUE -->|Channel receive| WORKER

    WORKER --> W2
    W2 --> P1 --> P2
    P2 -->|Status OK| P3 --> P4
    P4 --> P5 --> P6 --> P7
    P7 --> RES1

    P2 -->|Status Error| RES2

    W1 -.->|Waiting for<br/>next request| W3

    style QUEUE fill:#ffeb3b
    style WORKER fill:#fff9c4
    style P1 fill:#e8f5e9
    style P7 fill:#e8f5e9
    style RES1 fill:#c8e6c9
    style RES2 fill:#ffcccc
```

---

## Print Worker Lifecycle

```mermaid
sequenceDiagram
    participant init as init()
    participant Worker as Print Worker<br/>(Goroutine)
    participant Queue as printQueue<br/>Channel
    participant Handler as PrintRequest<br/>Handler

    init->>Worker: startPrintWorker()
    Note over Worker: go func() {<br/>for req := range printQueue { ... }<br/>}()

    Note right of Worker: 🎯 State: IDLE<br/>Waiting for requests...

    activate Worker

    Handler->>Queue: PrintFileQueued(P1, path1, copies)
    Note over Queue: Add PrintRequest 1<br/>Buffer: 1/1000

    Queue-->>Worker: Send on channel
    Worker->>Worker: req := <-printQueue<br/>Receive P1

    Note right of Worker: 🎯 State: PROCESSING<br/>Processing P1...

    Worker->>Worker: go func(r PrintRequest) {<br/>printFile(P1)<br/>}(req)
    Note over Worker: ⏱️ Spawns sub-goroutine<br/>for actual print work

    Worker->>Worker: Continue loop<br/>back to wait
    Note right of Worker: 🎯 State: IDLE<br/>Ready for next request

    Handler->>Queue: PrintFileQueued(P2, path2, copies)
    Note over Queue: Add PrintRequest 2<br/>Buffer: 1/1000

    Queue-->>Worker: Send on channel
    Worker->>Worker: req := <-printQueue<br/>Receive P2

    par Sub-goroutine P1 (background)
        Note over Handler: P1 printing...<br/>with retry logic<br/>Max retry: 1 hour
    and Main Worker (ready)
        Worker->>Worker: Process P2
    end

    deactivate Worker
```

---

## Print File Processing with Retry

```mermaid
flowchart TD
    Start["🎯 printFile(printer, path, copies)"] --> Connect["🔌 ConnectPrinter<br/>TCP DialTimeout 5s"]

    Connect --> ConnOK{Connected?}
    ConnOK -->|❌ Error| Retry["⏱️ Check retry logic"]

    ConnOK -->|✅ OK| Status["📊 Check Printer Status<br/>Send: 0x1D 0x72 0x01"]

    Status --> StatusOK{Status OK?}
    StatusOK -->|❌ Error| Retry
    StatusOK -->|✅ OK| OpenFile["📂 Open image file"]

    OpenFile --> OpenOK{File OK?}
    OpenOK -->|❌ Error| Retry
    OpenOK -->|✅ OK| Decode["🖼️ Decode image<br/>JPEG, PNG, etc"]

    Decode --> DecodeOK{Decode OK?}
    DecodeOK -->|❌ Error| Retry
    DecodeOK -->|✅ OK| Resize["📏 Resize if needed<br/>MaxWidth: 576 dots"]

    Resize --> Convert["⚫ Convert to 1-bit<br/>Black & White binary"]
    Convert --> BuildCmd["🔨 Build ESC/POS<br/>print command"]
    BuildCmd --> PrintLoop["🔄 For i=0 to copies-1"]

    PrintLoop --> Write["📨 Write print command"]
    Write --> WriteOK{Written?}
    WriteOK -->|❌ Error| Retry
    WriteOK -->|✅ OK| Blank["📄 Write blank lines<br/>0x0A x4"]

    Blank --> Cut["✂️ Write cut command<br/>0x1D 0x56 0x00"]
    Cut --> Sleep["😴 Sleep 1 second"]
    Sleep --> MoreCopies{More copies?}
    MoreCopies -->|Yes| Write
    MoreCopies -->|No| Cleanup["🗑️ Delete temp file"]

    Cleanup --> CloseConn["❌ conn.Close()"]
    CloseConn --> Success["✅ Return nil"]

    Retry --> CheckTimeout["⏱️ Check timeout"]
    CheckTimeout --> TimeoutOK{time < 1hr?}

    TimeoutOK -->|❌ Expired| DeleteFile["🗑️ Delete temp file<br/>Job expired"]
    DeleteFile --> Failed["❌ Break retry loop"]

    TimeoutOK -->|✅ Within time| IncRetry["📈 RetryCount++"]
    IncRetry --> Wait["😴 Sleep 30 seconds"]
    Wait --> PrintAgain["🔄 Retry printFile()"]
    PrintAgain --> Connect

    Success --> End["✅ End"]
    Failed --> End

    style Start fill:#e3f2fd
    style Success fill:#c8e6c9
    style Failed fill:#ffcccc
    style Retry fill:#fff9c4
    style Connect fill:#f3e5f5
    style Status fill:#f3e5f5
    style PrintLoop fill:#fff9c4
    style DeleteFile fill:#ffcccc
```

---

## Complete Request Lifecycle Timeline

```mermaid
timeline
    title Job Print Request - Complete Timeline

    section Request
        T+0ms: Client sends POST request
        T+5ms: Handler parses form data
        T+10ms: Handler validates input
        T+15ms: Service.JobPrint() called

    section Config
        T+20ms: Load config.json from disk
        T+25ms: Parse JSON config
        T+30ms: Filter printers by type

    section File Operations
        T+40ms: Generate temp file path
        T+45ms: Save uploaded file
        T+50ms: Send file to Telegram
        T+55ms: Open file handle

    section Queue
        T+60ms: Loop - For each printer
        T+65ms: Create PrintRequest object
        T+70ms: Send to queue channel
        T+75ms: Handler returns immediately
        T+80ms: Client receives 200 OK response ✅

    section Ping
        T+85ms: Loop - For each printer
        T+90ms: TCP connect to printer
        T+95ms: Close connection (ping)
        T+100ms: Service returns

    section Background (Worker)
        T+70ms: Worker receives request from queue
        T+75ms: Spawn sub-goroutine for print
        T+100ms: Sub-goroutine connects to printer
        T+105ms: Check printer status
        T+110ms: Decode image file
        T+115ms: Build print command
        T+120ms: Send print command (copy 1)
        T+125ms: Sleep 1 second
        T+126ms: Send print command (copy 2)
        T+131ms: Send cut command
        T+135ms: Delete temp file
        T+140ms: Sub-goroutine completes ✅
```

---

## Flow Chart - Detailed Process

```mermaid
flowchart TD
    Start([Client gửi POST request]) --> Check1{Form data<br/>hợp lệ?}

    Check1 -->|Không| Error1["❌ Error:<br/>Failed to get form data"]
    Error1 --> Response1["Response 400"]

    Check1 -->|Có| Check2{Files<br/>có tồn tại?}
    Check2 -->|Không| Error2["❌ Error:<br/>no_files_uploaded"]
    Error2 --> Response2["Response 400"]

    Check2 -->|Có| SetDefault["✓ Set default copies = 1<br/>if empty"]

    SetDefault --> LoadConfig["📖 Load config.json"]
    LoadConfig --> ParseJSON["🔄 Parse JSON to array"]

    ParseJSON --> Filter["🔍 Filter printers by type<br/>(Match printType)"]

    Filter --> Check3{Printers<br/>found?}
    Check3 -->|Không| Error3["❌ Error:<br/>PRINT_NOT_FOUND"]
    Error3 --> Response3["Response 400"]

    Check3 -->|Có| GenPath["📝 Generate temp file path<br/>uploads/YYYYMMDDhhmmss_file"]

    GenPath --> SaveFile["💾 Save uploaded file<br/>to temp path"]
    SaveFile --> Check4{File saved<br/>successfully?}

    Check4 -->|Lỗi| Error4["❌ Error:<br/>Failed to save file"]
    Error4 --> Response4["Response 400"]

    Check4 -->|OK| SendTele["📤 Send file to<br/>Telegram Bot"]
    SendTele --> OpenFile["📂 Open file for reading"]

    OpenFile --> LoopPrint["🔄 Loop: For each printer"]
    LoopPrint --> Queue["📨 PrintFileQueued<br/>(printer, path, copies)"]
    Queue --> Check5{Queue<br/>success?}

    Check5 -->|Lỗi| Error5["❌ Error:<br/>Failed to queue print"]
    Error5 --> Response5["Response 400"]

    Check5 -->|OK| NextPrinter{More<br/>printers?}
    NextPrinter -->|Yes| Queue
    NextPrinter -->|No| LoopPing["🔄 Loop: Ping all printers"]

    LoopPing --> Connect["🔌 ConnectPrinter<br/>(TCP connection)"]
    Connect --> CloseConn["❌ Close connection<br/>(Ping signal)"]
    CloseConn --> NextPing{More<br/>printers?}

    NextPing -->|Yes| Connect
    NextPing -->|No| Log["📝 Log: print done"]

    Log --> Success["✅ Return nil<br/>(No error)"]
    Success --> Response200["Response 200<br/>OK"]

    Response1 --> End([End])
    Response2 --> End
    Response3 --> End
    Response4 --> End
    Response5 --> End
    Response200 --> End

    style Start fill:#e1f5e1
    style End fill:#ffe1e1
    style Error1 fill:#ffcccc
    style Error2 fill:#ffcccc
    style Error3 fill:#ffcccc
    style Error4 fill:#ffcccc
    style Error5 fill:#ffcccc
    style Response200 fill:#ccffcc
```

---

## File Locations

| Loại          | Vị trí                            |
| ------------- | --------------------------------- |
| Temp files    | `uploads/YYYYMMDDhhmmss_filename` |
| Config        | `config/config.json`              |
| Device config | `config/device.json`              |

---

## Notes

1. **Copies**: Được gửi tới `PrintFileQueued()` để kiểm soát số bản in
2. **Temp Files**: Lưu với timestamp để tránh trùng tên
3. **Telegram**: Gửi file để logging hoặc monitoring
4. **Ping**: Kết nối TCP cuối cùng để đánh thức printer
5. **Default**: Nếu không có copies → mặc định in 1 bản
