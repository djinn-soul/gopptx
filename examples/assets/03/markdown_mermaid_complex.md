# Markdown To PPTX To Native PDF
- Complex markdown deck with many mermaid diagrams
- Generated through task example: 03-markdown-to-pptx
- Output artifacts:
  - 03_markdown_mermaid_complex.md
  - 03_markdown_mermaid_complex.pptx
  - 03_markdown_mermaid_complex.pdf

---

# Mermaid: Flowchart
```mermaid
flowchart LR
A[Start] --> B{Decision}
B -- Yes --> C[Ship]
B -- No --> D[Revise]
```

---

# Mermaid: Sequence
```mermaid
sequenceDiagram
Alice->>Bob: Hello Bob, how are you?
Bob-->>Alice: Jolly good!
```

---

# Mermaid: Pie
```mermaid
pie title Pets adopted by volunteers
"Dogs" : 386
"Cats" : 85
"Rats" : 15
```

---

# Mermaid: Gantt
```mermaid
gantt
title A Gantt Diagram
section Section
A task :a1, 2014-01-01, 30d
Another task :after a1, 20d
```

---

# Mermaid: Timeline
```mermaid
timeline
title History of Social Media Platform
2002 : LinkedIn
2004 : Facebook
: Google
```

---

# Mermaid: Quadrant
```mermaid
quadrantChart
title Reach and engagement of campaigns
x-axis Low Reach --> High Reach
y-axis Low Engagement --> High Engagement
quadrant-1 We should expand
Campaign A: [0.3, 0.6]
```

---

# Mermaid: Class
```mermaid
classDiagram
class Animal {
    +String name
    +isMammal()
}
class Dog {
    +bark()
}
Animal <|-- Dog
```

---

# Mermaid: State
```mermaid
stateDiagram-v2
[*] --> First
First --> Second
Second --> [*]
```

---

# Mermaid: ER
```mermaid
erDiagram
CUSTOMER ||--o{ ORDER : places
CUSTOMER {
    string name
    string email
}
ORDER {
    int orderNumber
}
```

---

# Mermaid: Mindmap
```mermaid
mindmap
root((mindmap))
    Origins
        Long history
    Research
        On effectiveness
```

---

# Mermaid: Journey
```mermaid
journey
title My working day
section Go to work
    Make tea: 5: Me
    Do work: 1: Me, Cat
```

---

# Mermaid: GitGraph
```mermaid
gitGraph
commit
commit
branch develop
checkout develop
commit
checkout main
merge develop
```

---

# Verification Checklist
| Step | Status |
|---|---|
| Markdown parsed | PASS |
| PPTX generated | PASS |
| Native PDF exported | PASS |

```go
slides, err := pptx.SlidesFromMarkdown(markdown)
if err != nil { return err }
err = export.PDFWithOptions("Deck", slides, "out.pdf",
  export.PDFOptions{Driver: export.PDFDriverNative})
```
