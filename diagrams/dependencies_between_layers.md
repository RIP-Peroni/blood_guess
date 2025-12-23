```mermaid
flowchart LR
    E[External Layer] --> IA[Interface Adapters]
    IA --> A[Application Layer]
    A --> D[Domain Layer]
    
    style D fill:#e8f5e8
    style A fill:#fff3e0
    style IA fill:#e3f2fd
    style E fill:#f3e5f5
```